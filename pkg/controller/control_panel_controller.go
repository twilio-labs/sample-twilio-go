package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControlPanelController struct {
}

func NewControlPanelController() *ControlPanelController {
	return &ControlPanelController{}
}

func (ctr *ControlPanelController) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "campaign-controls.html", gin.H{})
}
