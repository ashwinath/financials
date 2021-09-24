package controller

import "net/http"

type healthController struct {
	controller
}

func (c *healthController) Alive(w http.ResponseWriter, _ *http.Request) {
	c.Ok(w, struct{}{})
}

func (c *healthController) Ready(w http.ResponseWriter, _ *http.Request) {
	db, err := c.context.DB.DB()
	if err != nil {
		c.ServiceUnavailable(w, struct{}{})
		return
	}

	if err := db.Ping(); err != nil {
		c.ServiceUnavailable(w, struct{}{})
		return
	}

	c.Ok(w, struct{}{})
}
