package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func CheckError(controller beego.Controller, err error, status string) {
	if err != nil {
		logs.Error(err)
		resp := map[string]interface{}{
			"data ":   make([]interface{}, 0),
			"message": err.Error(),
			"status":  status,
		}
		Send(controller, resp)
	}
}

func Send(controller beego.Controller, resp interface{}) {
	controller.Data["json"] = resp
	serveJsonError := controller.ServeJSON()
	if serveJsonError != nil {
		logs.Error(serveJsonError)
		controller.Abort("500")
	}
	controller.Finish()
}

func CheckExists(controller beego.Controller, exists bool, status string) {
	if !exists {
		resp := map[string]interface{}{
			"data ":   make([]interface{}, 0),
			"message": "Not Exists",
			"status":  status,
		}
		Send(controller, resp)
	}
}

func CheckCustomError(controller beego.Controller, err, customError error, status, customMessage string) {
	var resp interface{}
	if err != nil {
		resp = map[string]interface{}{
			"data ":   make([]interface{}, 0),
			"message": "Not Exists",
			"status":  status,
		}
		if err == customError {
			resp = map[string]interface{}{
				"data ":   make([]interface{}, 0),
				"message": customMessage,
				"status":  status,
			}
		}
		Send(controller, resp)
	}
}
