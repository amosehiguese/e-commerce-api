package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
