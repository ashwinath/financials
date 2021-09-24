package controller

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
)

type fooController struct {
	controller
}

type params struct {
	Foo string `schema:"foo" validate:"required,lte=5"`
}

type body struct {
	Bar string `json:"bar" validate:"required,gte=5"`
}

func (c *fooController) testController(_ http.ResponseWriter, r *http.Request) (*params, *body, error) {
	p := params{}
	err := c.getParams(r, &p)
	if err != nil {
		return nil, nil, err
	}

	b := body{}
	err = c.getBody(r, &b)
	if err != nil {
		return nil, nil, err
	}
	return &p, &b, nil
}

func TestBaseController(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		hasError       bool
		expectedParams string
		expectedBody   string
	}{
		{
			name:           "nominal",
			method:         http.MethodPost,
			path:           "/foo?foo=hello",
			body:           `{"bar": "hello-world"}`,
			hasError:       false,
			expectedParams: "hello",
			expectedBody:   "hello-world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "Application/JSON")
			assert.Nil(t, err)

			c := fooController{
				controller: controller{
					decoder:   schema.NewDecoder(),
					validator: validator.New(),
				},
			}

			p, b, err := c.testController(nil, r)
			if tt.hasError {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.expectedParams, p.Foo)
			assert.Equal(t, tt.expectedBody, b.Bar)
		})
	}
}
