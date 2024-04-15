package services

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_formato_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var DatavalidaT = []string{}
var estadoPlan string

func ConsultarFormato(id string) (interface{}, error) {
	var respuesta map[string]interface{}
	var hijos []models.Nodo
	var plan map[string]interface{}
	var hijosID []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &hijos)
		request.LimpiezaRespuestaRefactor(respuesta, &hijosID)
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {
			request.LimpiezaRespuestaRefactor(respuesta, &plan)
			Limpia(plan)
			if arbol, err := ConstruirArbol(hijos, hijosID); err == nil {
				return arbol, nil
			} else {
				logs.Error("Error en Service.ConsultarFormato -->", arbol)
				return nil, errors.New(err.Error())
			}	
		} else {
			logs.Error("Error en Service.ConsultarFormato -->", respuesta)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error en Service.ConsultarFormato -->", respuesta)
		return nil, errors.New(err.Error())
	}
}

func Limpia(plan map[string]interface{}) {
	DatavalidaT = []string{}
	if jsonString, err := json.Marshal(plan["estado_plan_id"]); err == nil {
		if err := json.Unmarshal(jsonString, &estadoPlan); err != nil {
			logs.Error("Error en Service.Limpia -->", err)
		}
	} else {
		logs.Error("Error en Service.Limpia -->", err)
	}
}

func ConstruirArbol(hijos []models.Nodo, hijosID []map[string]interface{}) ([][]map[string]interface{}, error) {
	var arbol []map[string]interface{}
	var requeridos []map[string]interface{}
	var nodo []models.NodoDetalle
	var respuesta map[string]interface{}
	var resultado [][]map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if hijos[i].Activo {
			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijosID[i]["_id"]
			forkData["nombre"] = hijos[i].Nombre
			jsonString, _ := json.Marshal(hijosID[i]["_id"])
			if err := json.Unmarshal(jsonString, &id); err != nil {
				logs.Error("Error en Service.ConstruirArbol -->", err)
				return nil, errors.New(err.Error())
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &respuesta); err == nil {
				request.LimpiezaRespuestaRefactor(respuesta, &nodo)
				if len(nodo) > 0 {
					var deta map[string]interface{}
					if err := json.Unmarshal([]byte(nodo[0].Dato), &deta); err == nil {
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
					} else {
						logs.Error("Error en Service.ConstruirArbol -->", err)
						return nil, errors.New(err.Error())
					}
				}
			} else {
				logs.Error("Error en Service.ConstruirArbol -->", respuesta)
				return nil, errors.New(err.Error())
			}
			if len(hijos[i].Hijos) > 0 {
				if respuestaHijos, err:= ConsultarHijos(hijos[i].Hijos); err == nil{
					if len(respuestaHijos) == 0 {
						forkData["sub"] = ""
					} else {
						forkData["sub"] = make([]map[string]interface{}, len(respuestaHijos))
						forkData["sub"] = respuestaHijos
					}
				}else{
					logs.Error("Error en Service.ConstruirArbol -->", respuestaHijos)
					return nil, errors.New(err.Error())
				}
			}
			arbol = append(arbol, forkData)
			add(id)
		}
	}
	requeridos = Convertir(DatavalidaT)
	resultado = append(resultado, arbol)
	resultado = append(resultado, requeridos)
	return resultado, nil
}

func ConsultarHijos(hijos []string) ([]map[string]interface{}, error) {
	var res map[string]interface{}
	var resp map[string]interface{}
	var nodo models.Nodo
	var nodoId map[string]interface{}
	var detalle []models.NodoDetalle
	var Arbolhijos []map[string]interface{}

	for _, hijo := range hijos {
		forkData := make(map[string]interface{})
		var id string

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijo, &res); err == nil {
			request.LimpiezaRespuestaRefactor(res, &nodo)
			request.LimpiezaRespuestaRefactor(res, &nodoId)
			if nodo.Activo == true {
				forkData["id"] = nodoId["_id"]
				forkData["nombre"] = nodo.Nombre
				if jsonString, err := json.Marshal(nodoId["_id"]); err == nil {
					if err := json.Unmarshal(jsonString, &id); err == nil {
						if err_ := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
							request.LimpiezaRespuestaRefactor(resp, &detalle)
							if len(detalle) > 0 {
								var deta map[string]interface{}
								if err := json.Unmarshal([]byte(detalle[0].Dato), &deta); err == nil {
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
								} else {
									logs.Error("Error en Service.ConsultarHijos -->", err)
									return nil, errors.New(err.Error())
								}
							}
						} else {
							logs.Error("Error en Service.ConsultarHijos -->", resp)
							return nil, errors.New(err_.Error())
						}

						if len(nodo.Hijos) > 0 {
							if respuestaHijos, err:= ConsultarHijos(nodo.Hijos); err == nil{
								if len(respuestaHijos) == 0 {
									forkData["sub"] = ""
								} else {
									forkData["sub"] = respuestaHijos
								}
							} else {
								logs.Error("Error en Service.ConsultarHijos -->", respuestaHijos)
								return nil, errors.New(err.Error())
							}
							
						}
						Arbolhijos = append(Arbolhijos, forkData)
					}else {
						logs.Error("Error en Service.ConsultarHijos -->", err)
						return nil, errors.New(err.Error())
					}
				} else {
					logs.Error("Error en Service.ConsultarHijos -->", jsonString)
					return nil, errors.New(err.Error())
				}
			}
			add(id)
		} else {
			logs.Error("Error en Service.ConsultarHijos -->", res)
			return nil, errors.New(err.Error())
		}
	}
	return Arbolhijos, nil
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