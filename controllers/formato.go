package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_formato_mid/helpers"
	"github.com/udistrital/planeacion_formato_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// FormatoController operations for Formato
type FormatoController struct {
	beego.Controller
}

func (c *FormatoController) URLMapping() {
	c.Mapping("ConsultarFormato", c.ConsultarFormato)
}

// ConsultarFormato ...
// @Title ConsultarFormato
// @Description Consulta el Formato por id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formato
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FormatoController) ConsultarFormato() {

	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "FormatoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	var hijos []models.Nodo
	var plan map[string]interface{}
	var hijosID []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		request.LimpiezaRespuestaRefactor(res, &hijos)
		request.LimpiezaRespuestaRefactor(res, &hijosID)
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res)
		if err != nil {
			return
		}
		request.LimpiezaRespuestaRefactor(res, &plan)
		helpers.Limpia(plan)
		arbol := helpers.ConstruirArbol(hijos, hijosID)
		c.Data["json"] = arbol
	} else {
		panic(err)
	}
	c.ServeJSON()
}
