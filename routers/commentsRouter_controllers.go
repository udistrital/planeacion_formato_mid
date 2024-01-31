package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_formato_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_formato_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "ConsultarFormato",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
