package handler

import (
	"github.com/AeroAgency/go-admin-api/application/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PingHandler struct {
	errorService service.Error
	modelService *service.ModelsService
}

func NewPingHandler(errorService service.Error, modelService *service.ModelsService) *PingHandler {
	return &PingHandler{errorService: errorService, modelService: modelService}
}

func (ph PingHandler) PingModel(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	err := ph.modelService.PingModel(modelCode)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (ph PingHandler) PingModelField(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	modelFieldCode := c.Params.ByName("model-field-code")
	err := ph.modelService.PingModelField(modelCode, modelFieldCode)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, nil)
}
