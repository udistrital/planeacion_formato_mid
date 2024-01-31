package controllers

import (
	"net/http"
	"testing"
)

func TestConsultarFormato(t *testing.T) {
	if response, err := http.Get("http://localhost:9011/v1/formato/61305c7edf020f065956ba26"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ConsultarFormato Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ConsultarFormato Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ConsultarFormato:", err.Error())
		t.Fail()
	}
}
