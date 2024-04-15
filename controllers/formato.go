package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_formato_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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
	defer errorhandler.HandlePanic(&c.Controller)

	id := c.Ctx.Input.Param(":id")

	if resultado, err := services.ConsultarFormato(id); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}