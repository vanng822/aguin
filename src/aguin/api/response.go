package api

import (
	"aguin/crypto"
	"github.com/martini-contrib/render"
	"net/http"
)


/* Serving responses */
func ServeResponse(status int, render render.Render, result interface{}, requestData RequestData, setting AguinSetting) {
	if setting.conf.EncryptionEnabled {
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
