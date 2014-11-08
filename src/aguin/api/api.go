package api

import (
	"aguin/model"
	"aguin/utils"
	"aguin/validator"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RequestData struct {
	app       model.Application
	data      url.Values
}

type ApplicationSetting struct {
	dbSession *mgo.Session
	log *log.Logger // General log
}

func serveOK(render render.Render) {
	render.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func serveForbidden(render render.Render) {
	render.JSON(http.StatusForbidden, map[string]interface{}{"status": "FORBIDDEN"})
}

func serveBadRequestData(render render.Render) {
	render.JSON(http.StatusBadRequest, map[string]interface{}{"status": "ERROR", "msg": "Invalid data"})
}

func serveBadRequestJson(render render.Render) {
	render.JSON(http.StatusBadRequest, map[string]interface{}{"status": "ERROR", "msg": "Invalid json"})
}

func serveInternalServerError(render render.Render) {
	render.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "ERROR"})
}

func VerifyRequest() interface{} {
	return func(c martini.Context, res http.ResponseWriter, req *http.Request, render render.Render) {
		log := utils.GetLogger("aguin")
		// Recover from panic and serve internal server error in json format
		defer func() {
			if r := recover(); r != nil {
				log.Print(r)
				serveInternalServerError(render)
			}
		}()
	
		apiKey := req.Header.Get("X-AGUIN-API-KEY")

		if apiKey == "" {
			serveForbidden(render)
			return
		}
		// If empty then data must be encrypted
		apiSecret := req.Header.Get("X-AGUIN-API-SECRET")

		app := model.Application{}
		
		setting := ApplicationSetting{}
		requestData := RequestData{}
		setting.log = log
		setting.dbSession = model.Session()
		
		defer setting.dbSession.Close()

		err := model.AppCollection(setting.dbSession).FindId(bson.ObjectIdHex(apiKey)).One(&app)

		log.Print(app, err)

		if err != nil {
			serveForbidden(render)
			return
		}

		requestData.app = app

		// parse form for data
		req.ParseForm()
		// if secret empty it must be crypted
		if apiSecret == "" {
			// decrypting data here

		} else {
			// have secret, compare with the one in database
			if apiSecret != app.Secret {
				serveForbidden(render)
				return
			}
		}
		requestData.data = req.Form
		c.Map(requestData)
		c.Map(setting)
		c.Next()
	}
}

func IndexGet(res http.ResponseWriter, req *http.Request, render render.Render, requestData RequestData, setting ApplicationSetting) {
	log := setting.log
	
	var results []model.Entity
	da := []byte(requestData.data.Get("message"))
	
	data3, err := utils.Bytes2json(&da)
	if err != nil {
		log.Print(err)
		serveBadRequestJson(render)
		return
	}
	criteria := validator.ValidateSearch(data3)
	if criteria.Validated == false {
		serveBadRequestData(render)
		return
	}
	err = model.EntityCollection(setting.dbSession).Find(bson.M{"name": criteria.Entity, "appid": requestData.app.Id}).All(&results)
	
	render.JSON(http.StatusOK, results)
}

func IndexPost(res http.ResponseWriter, req *http.Request, render render.Render, requestData RequestData, setting ApplicationSetting) {
	log := setting.log
	da := []byte(requestData.data.Get("message"))
	data3, err := utils.Bytes2json(&da)
	if err != nil {
		log.Print(err)
		serveBadRequestJson(render)
		return
	}
	entity, data, validated := validator.ValidateEntity(data3)
	if validated && entity != "" {
		doc := model.Entity{}
		doc.Name = entity
		doc.AppId = requestData.app.Id
		doc.CreatedAt = time.Now()
		doc.Data = data
		err = doc.Save(setting.dbSession)
		if err == nil {
			serveOK(render)
			return
		} else {
			log.Print(err)
			serveInternalServerError(render)
			return
		}
	}
	serveBadRequestData(render)
}

func IndexStatus(render render.Render) {
	serveOK(render)
}

func NotFound(render render.Render) {
	render.JSON(http.StatusNotFound, map[string]interface{}{"status": "ERROR", "msg": "Not found"})
}
