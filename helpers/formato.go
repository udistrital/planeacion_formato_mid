package helpers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_formato_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var DatavalidaT = []string{}
var estadoPlan string

func Limpia(plan map[string]interface{}) {
	DatavalidaT = []string{}
	jsonString, _ := json.Marshal(plan["estado_plan_id"])
	json.Unmarshal(jsonString, &estadoPlan)
}

func ConstruirArbol(hijos []models.Nodo, hijosID []map[string]interface{}) [][]map[string]interface{} {
	var arbol []map[string]interface{}
	var requeridos []map[string]interface{}
	var nodo []models.NodoDetalle
	var res map[string]interface{}
	var resultado [][]map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if hijos[i].Activo {
			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijosID[i]["_id"]
			forkData["nombre"] = hijos[i].Nombre
			jsonString, _ := json.Marshal(hijosID[i]["_id"])
			json.Unmarshal(jsonString, &id)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &res); err == nil {
				request.LimpiezaRespuestaRefactor(res, &nodo)
				if len(nodo) > 0 {
					var deta map[string]interface{}
					json.Unmarshal([]byte(nodo[0].Dato), &deta)

					if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] == nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
					} else if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] != nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
						forkData["options"] = deta["options"]
					} else {
						forkData["type"] = " "
						forkData["required"] = " "
					}
				}
			}
			if len(hijos[i].Hijos) > 0 {
				if len(ConsultarHijos(hijos[i].Hijos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = make([]map[string]interface{}, len(ConsultarHijos(hijos[i].Hijos)))
					forkData["sub"] = ConsultarHijos(hijos[i].Hijos)
				}
			}
			arbol = append(arbol, forkData)
			add(id)
		}
	}
	requeridos = Convertir(DatavalidaT)
	resultado = append(resultado, arbol)
	resultado = append(resultado, requeridos)
	return resultado
}

func ConsultarHijos(hijos []string) (Arbolhijos []map[string]interface{}) {
	var res map[string]interface{}
	var resp map[string]interface{}
	var nodo models.Nodo
	var nodoId map[string]interface{}
	var detalle []models.NodoDetalle

	for _, hijo := range hijos {
		forkData := make(map[string]interface{})
		var id string

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijo, &res); err != nil {
			return
		}
		request.LimpiezaRespuestaRefactor(res, &nodo)
		request.LimpiezaRespuestaRefactor(res, &nodoId)
		if nodo.Activo == true {
			forkData["id"] = nodoId["_id"]
			forkData["nombre"] = nodo.Nombre
			jsonString, _ := json.Marshal(nodoId["_id"])
			json.Unmarshal(jsonString, &id)

			if err_ := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
				request.LimpiezaRespuestaRefactor(resp, &detalle)
				if len(detalle) > 0 {
					var deta map[string]interface{}
					json.Unmarshal([]byte(detalle[0].Dato), &deta)

					if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] == nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
					} else if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] != nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
						forkData["options"] = deta["options"]
					} else {
						forkData["type"] = " "
						forkData["required"] = " "
					}
				}
			}
			if len(nodo.Hijos) > 0 {
				if len(ConsultarHijos(nodo.Hijos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = ConsultarHijos(nodo.Hijos)
				}
			}
			Arbolhijos = append(Arbolhijos, forkData)
		}
		add(id)
	}
	return
}

func add(id string) {
	if !request.Contains(DatavalidaT, id) {
		DatavalidaT = append(DatavalidaT, id)
	}
}

func Convertir(valido []string) []map[string]interface{} {
	var validadores []map[string]interface{}
	forkData := make(map[string]interface{})

	for _, v := range valido {
		if v == "" {

		} else {
			forkData[v] = ""
			forkData[v+"_o"] = ""
		}
	}

	validadores = append(validadores, forkData)
	return validadores
}
