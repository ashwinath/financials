package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ashwinath/financials/api/context"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type controller struct {
	context   *context.Context
	decoder   *schema.Decoder
	validator *validator.Validate
}

// getBody unmarshals the body into a Go struct
// Uses the json struct tag
func (c *controller) getBody(r *http.Request, dst interface{}) error {
	d := json.NewDecoder(r.Body)

	if err := d.Decode(dst); err != nil {
		return err
	}

	return c.validate(dst)
}

// getparams unmarshals the params into a Go struct
// Uses the schema struct tag
func (c *controller) getParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := c.decoder.Decode(dst, r.Form); err != nil {
		return err
	}

	return c.validate(dst)
}

func (c *controller) validate(dst interface{}) error {
	// Here is is possible to give custom messages.
	// See https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
	return c.validator.Struct(dst)
}

func serviceUnavailable(w http.ResponseWriter, body interface{}) {
	writeJSON(w, http.StatusServiceUnavailable, body)
}

func ok(w http.ResponseWriter, body interface{}) {
	writeJSON(w, http.StatusOK, body)
}

func created(w http.ResponseWriter, body interface{}) {
	writeJSON(w, http.StatusCreated, body)
}

type errorResponse struct {
	Description string `json:"description"`
	Message     string `json:"message"`
}

func badRequest(w http.ResponseWriter, description string, message string) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Description: description,
		Message:     message,
	})
}

func internalServiceError(w http.ResponseWriter, description string, message string) {
	writeJSON(w, http.StatusInternalServerError, errorResponse{
		Description: description,
		Message:     message,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
