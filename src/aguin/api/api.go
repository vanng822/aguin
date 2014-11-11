package api

import (
	"aguin/crypto"
	"aguin/model"
	"aguin/utils"
	"aguin/validator"
	"aguin/config"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type RequestData struct {
	crypted bool
	app     model.Application
	message map[string]interface{}
}

type AguinSetting struct {
	dbSession *mgo.Session
	log       *log.Logger // General log
}

func VerifyRequest() interface{} {
	return func(c martini.Context, res http.ResponseWriter, req *http.Request, render render.Render) {
		log := utils.GetLogger("aguin")
		// Recover from panic and serve internal server error in json format
		defer func() {
			if err := recover(); err != nil {
				stack := utils.Stack(3)
				log.Printf("PANIC: %s\n%s", err, stack)
				serveInternalServerError(render)
			}
		}()
		appConfig := config.AppConf()
		apiKey := req.Header.Get("X-AGUIN-API-KEY")

		if apiKey == "" || !validator.ValidAPIKey(apiKey) {
			serveForbidden(render)
			return
		}
		// If empty then data must be encrypted
		apiSecret := req.Header.Get("X-AGUIN-API-SECRET")

		app := model.Application{}

		setting := AguinSetting{}
		requestData := RequestData{}
		setting.log = log
		setting.dbSession = model.Session()
		defer setting.dbSession.Close()
		err := model.AppCollection(setting.dbSession).FindId(bson.ObjectIdHex(apiKey)).One(&app)
		if err != nil {
			serveForbidden(render)
			return
		}

		requestData.app = app

		// parse form for data
		req.ParseForm()
		if appConfig.EncryptionEnabled {
			// decrypting data here
			authKey := []byte(app.Secret)
			decryptedMessage, err := crypto.Decrypt(req.Form.Get("message"), authKey, authKey)
			if err != nil {
				setting.log.Print(err, req.Form.Get("message"))
				serveInternalServerError(render)
				return
			}
			requestData.message = decryptedMessage.(map[string]interface{})
			requestData.crypted = true
		} else {
			// have secret, compare with the one in database
			if apiSecret != app.Secret {
				serveForbidden(render)
				return
			}
			data3, err := utils.Bytes2json([]byte(req.Form.Get("message")))
			if err != nil {
				setting.log.Print(err)
				serveBadRequestJson(render)
				return
			}
			requestData.message = data3.(map[string]interface{})
			requestData.crypted = false
		}
		c.Map(requestData)
		c.Map(setting)
		c.Next()
	}
}

func IndexGet(res http.ResponseWriter, req *http.Request, render render.Render, requestData RequestData, setting AguinSetting) {
	criteria := validator.ValidateSearch(requestData.message)
	if criteria.Validated == false {
		setting.log.Print(requestData.message)
		serveBadRequestData(render)
		return
	}
	var results []model.Entity
	err := model.EntityCollection(setting.dbSession).Find(bson.M{"name": criteria.Entity, "appid": requestData.app.Id}).All(&results)
	if err != nil {
		setting.log.Println(err)
		serveInternalServerError(render)
		return
	}
	ServeResponse(http.StatusOK, render, results, requestData, setting)
}

func IndexPost(res http.ResponseWriter, req *http.Request, render render.Render, requestData RequestData, setting AguinSetting) {
	result := validator.ValidateEntity(requestData.message)
	if result.Validated && result.Entity != "" {
		doc := model.Entity{}
		doc.Name = result.Entity
		doc.AppId = requestData.app.Id
		doc.CreatedAt = time.Now()
		doc.Data = result.Data
		err := doc.Save(setting.dbSession)
		if err == nil {
			serveOK(render)
			return
		} else {
			setting.log.Print(err)
			serveInternalServerError(render)
			return
		}
	}
	setting.log.Print(requestData.message)
	serveBadRequestData(render)
}

func IndexStatus(render render.Render) {
	serveOK(render)
}

func NotFound(render render.Render) {
	serveNotFound(render)
}

/* Serving responses */
func ServeResponse(status int, render render.Render, result interface{}, requestData RequestData, setting AguinSetting) {
	if requestData.crypted {
		ServeSignedResponse(status, render, result, requestData, setting)
	} else {
		ServeUnsignedResponse(status, render, result, requestData, setting)
	}
}

func ServeUnsignedResponse(status int, render render.Render, result interface{}, requestData RequestData, setting AguinSetting) {
	render.JSON(status, map[string]interface{}{"result": result, "encrypted": false})
}

func ServeSignedResponse(status int, render render.Render, result interface{}, requestData RequestData, setting AguinSetting) {
	authKey := []byte(requestData.app.Secret)
	cryptedResult, err := crypto.Encrypt(result, authKey, authKey)
	if err != nil {
		serveInternalServerError(render)
		return
	}
	render.JSON(status, map[string]interface{}{"result": cryptedResult, "encrypted": true})
}

func serveOK(render render.Render) {
	render.JSON(http.StatusOK, map[string]interface{}{"result": map[string]interface{}{"status": "OK"}, "encrypted": false})
}

func serveForbidden(render render.Render) {
	render.JSON(http.StatusForbidden, map[string]interface{}{"result": map[string]interface{}{"status": "FORBIDDEN"}, "encrypted": false})
}

func serveBadRequestData(render render.Render) {
	render.JSON(http.StatusBadRequest, map[string]interface{}{"result": map[string]interface{}{"status": "ERROR", "msg": "Invalid data"}, "encrypted": false})
}

func serveBadRequestJson(render render.Render) {
	render.JSON(http.StatusBadRequest, map[string]interface{}{"result": map[string]interface{}{"status": "ERROR", "msg": "Invalid json"}, "encrypted": false})
}

func serveInternalServerError(render render.Render) {
	render.JSON(http.StatusInternalServerError, map[string]interface{}{"result": map[string]interface{}{"status": "ERROR", "msg": "Internal server error"}, "encrypted": false})
}

func serveNotFound(render render.Render) {
	render.JSON(http.StatusNotFound, map[string]interface{}{"result": map[string]interface{}{"status": "ERROR", "msg": "Not found"}, "encrypted": false})
}
