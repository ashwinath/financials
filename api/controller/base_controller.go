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

func (c *controller) ServiceUnavailable(w http.ResponseWriter, body interface{}) {
	c.WriteJSON(w, http.StatusServiceUnavailable, body)
}

func (c *controller) Ok(w http.ResponseWriter, body interface{}) {
	c.WriteJSON(w, http.StatusOK, body)
}

func (c *controller) WriteJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
