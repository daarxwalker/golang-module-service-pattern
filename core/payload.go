package core

import (
	"fmt"
	"net/http"

	"example/core/helper/jsonHelper"
)

type Payload interface {
	unauthorized()
	forbidden()
	unknownAction()
	ok(data interface{})
	error(e string)
	formErrors(data interface{})
}

type payload struct {
	provideService ProvideService
}

type responsePayload struct {
	Payload    interface{} `json:"payload,omitempty"`
	Module     string      `json:"module,omitempty"`
	Status     int         `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	FormErrors interface{} `json:"formErrors,omitempty"`
}

func newPayload(provideService ProvideService) Payload {
	return payload{
		provideService,
	}
}

func makePayload(module string, status int, payload interface{}) responsePayload {
	return responsePayload{
		Module:  module,
		Status:  status,
		Payload: payload,
	}
}

func makeFormErrorPayload(module string, status int, formErrors interface{}) responsePayload {
	return responsePayload{
		Module:     module,
		Status:     status,
		FormErrors: formErrors,
	}
}

func makeErrorPayload(module string, status int, message string) responsePayload {
	return responsePayload{
		Module:  module,
		Status:  status,
		Message: message,
	}
}

func getContentType(ct string) string {
	if len(ct) > 0 {
		return ct
	}
	return "application/json"
}

func (p payload) setPayload(status int, payload responsePayload) {
	w := p.provideService.GetResponse()
	r := jsonHelper.ParseJSON(payload)

	if w.Header().Get("Content-Transfer-Encoding") != "Binary" {
		w.Header().Set("Content-Type", getContentType(w.Header().Get("Content-Type")))
		w.WriteHeader(status)
		_, err := w.Write(r)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (p payload) unauthorized() {
	p.setPayload(
		http.StatusUnauthorized,
		makeErrorPayload(
			p.provideService.GetModuleName(),
			http.StatusUnauthorized,
			unauthorizedMessage,
		),
	)
}

func (p payload) forbidden() {
	p.setPayload(
		http.StatusForbidden,
		makeErrorPayload(
			p.provideService.GetModuleName(),
			http.StatusForbidden,
			forbiddenMessage,
		),
	)
}

func (p payload) unknownAction() {
	p.setPayload(
		http.StatusNotFound,
		makeErrorPayload(
			"",
			http.StatusNotFound,
			unknownActionMessage,
		),
	)
}

func (p payload) ok(data interface{}) {
	p.setPayload(
		http.StatusOK,
		makePayload(
			p.provideService.GetModuleName(),
			http.StatusOK,
			data,
		),
	)
}

func (p payload) error(e string) {
	p.setPayload(
		http.StatusBadRequest,
		makeErrorPayload(
			p.provideService.GetModuleName(),
			http.StatusBadRequest,
			e,
		),
	)
}

func (p payload) formErrors(data interface{}) {
	p.setPayload(
		http.StatusBadRequest,
		makeFormErrorPayload(
			p.provideService.GetModuleName(),
			http.StatusBadRequest,
			data,
		),
	)
}
