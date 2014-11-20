package api

import (
	"aguin/config"
	"aguin/crypto"
	"aguin/model"
	"aguin/utils"
	"aguin/validator"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type RequestData struct {
	app     model.Application
	message map[string]interface{}
}

type AguinSetting struct {
	dbSession *mgo.Session
	log       *utils.Logger // General log
	conf      *config.AppConfig
}

func VerifyRequest() interface{} {
	return func(c martini.Context, res http.ResponseWriter, req *http.Request, render render.Render) {
		log := utils.GetLogger("aguin")
		// Recover from panic and serve internal server error in json format
		defer func() {
			if err := recover(); err != nil {
				stack := utils.Stack(3)
				log.Error("PANIC: %s\n%s", err, stack)
				serveInternalServerError(render)
			}
		}()

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
			if err.Error() == "not found" {
				log.Info("Try to find application for key %s which is not found, error: %v", apiKey, err)
				serveForbidden(render)
			} else {
				log.Error("Error when finding application for key %s, error: %v", apiKey, err)
				serveInternalServerError(render)
			}
			return
		}

		if req.RequestURI == "/status" {
			IndexStatus(render)
			return
		}

		requestData.app = app
		setting.conf = config.AppConf()
		// parse form for data
		req.ParseForm()
		message := req.Form.Get("message")
		if message == "" {
			setting.log.Info("Request made without message")
			serveBadRequestJson(render)
			return
		}
		if setting.conf.EncryptionEnabled {
			// decrypting data here
			authKey := []byte(app.Secret)
			decryptedMessage, err := crypto.Decrypt(message, authKey, authKey)
			if err != nil || decryptedMessage == nil {
				setting.log.Error("error: %v, message: %v, decryptedMessage: %v", err, message, decryptedMessage)
				serveInternalServerError(render)
				return
			}
			requestData.message = decryptedMessage.(map[string]interface{})
		} else {
			// have secret, compare with the one in database
			if apiSecret != app.Secret {
				serveForbidden(render)
				return
			}
			data3, err := utils.Bytes2json([]byte(message))
			if err != nil {
				setting.log.Error("%v", err)
				serveBadRequestJson(render)
				return
			}
			requestData.message = data3.(map[string]interface{})
		}
		c.Map(requestData)
		c.Map(setting)
		c.Next()
	}
}

func IndexGet(res http.ResponseWriter, req *http.Request, render render.Render, requestData RequestData, setting AguinSetting) {
	criteria := validator.ValidateSearch(requestData.message)
	if criteria.Validated == false {
		setting.log.Info("%v", requestData.message)
		serveBadRequestData(render)
		return
	}
	setting.log.Debug("%v", requestData.message)
	setting.log.Debug("%v", criteria)
	var results []model.Entity
	err := model.EntityCollection(setting.dbSession).Find(
		bson.M{"name": criteria.Entity,
			"appid":      requestData.app.Id,
			"createdat": bson.M{"$gte": criteria.StartDate, "$lte": criteria.EndDate}}).All(&results)

	if err != nil {
		setting.log.Error("%v", err)
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
			setting.log.Error("%v", err)
			serveInternalServerError(render)
			return
		}
	}
	setting.log.Info("%v", requestData.message)
	serveBadRequestData(render)
}

// Should not do anything but return OK
func IndexStatus(render render.Render) {
	serveOK(render)
}

func NotFound(render render.Render) {
	serveNotFound(render)
}
