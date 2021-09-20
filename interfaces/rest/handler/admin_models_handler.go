package handler

import (
	"fmt"
	"github.com/AeroAgency/go-admin-api/application/service"
	appErrors "github.com/AeroAgency/go-admin-api/infrastructure/errors"
	"github.com/AeroAgency/go-admin-api/interfaces/rest/dto/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	pkgErrors "github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type AdminModelsHandler struct {
	errorService service.Error
	modelService *service.ModelsService
}

func NewPingHandler(errorService service.Error, modelService *service.ModelsService) *AdminModelsHandler {
	return &AdminModelsHandler{errorService: errorService, modelService: modelService}
}

func (ph AdminModelsHandler) PingModel(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	err := ph.modelService.PingModel(modelCode)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (ph AdminModelsHandler) PingModelField(c *gin.Context) {

	fieldType := c.Query("type")
	referenceLinkCode := c.Query("referenceCode")
	modelLinkCode := c.Query("modelCode")
	multiple, err := strconv.ParseBool(c.Query("multiple"))
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	modelFiledAddParams := models.ModelFieldAddParamsApiDto{
		Type:          fieldType,
		ReferenceCode: referenceLinkCode,
		ModelCode:     modelLinkCode,
		Multiple:      multiple,
	}

	modelCode := c.Params.ByName("model-code")
	modelFieldCode := c.Params.ByName("model-field-code")
	err = ph.modelService.PingModelField(modelCode, modelFieldCode, modelFiledAddParams)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (ph AdminModelsHandler) GetModelElementsPermissions(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	err := ph.modelService.PingModel(modelCode)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, []string{"read", "create", "edit", "delete"})
}

func (ph AdminModelsHandler) GetModelFilterValues(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		ph.errorService.HandleError(appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("value of param limit is invalid"))}, c)
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		ph.errorService.HandleError(appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("value of param offset is invalid"))}, c)
		return
	}
	dto := models.ModelFilterValuesParamsDto{
		ModelCode:               c.Params.ByName("model-code"),
		ModelFieldCode:          c.Params.ByName("model-field-code"),
		ModelFieldType:          c.Query("model_field_type"),
		ModelFieldModelCode:     c.Query("model_field_model_code"),
		ModelFieldReferenceCode: c.Query("model_field_reference_code"),
		Limit:                   limit,
		Offset:                  offset,
		Query:                   c.Query("query"),
		FilterId:                c.Query("filter_id"),
		Multiple:                c.Query("multiple") == "true",
	}
	filterValues, err := ph.modelService.GetModelFilterValues(dto)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, filterValues)
}

func (ph AdminModelsHandler) GetModelElementsList(c *gin.Context) {
	var modelElementsListParamsApiDto models.ModelElementsListParamsApiDto
	err := c.ShouldBindWith(&modelElementsListParamsApiDto, binding.JSON)
	if err != nil {
		err = appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("some of input params is invalid. err: %s", err))}
		ph.errorService.HandleError(err, c)
		return
	}
	modelElements, err := ph.modelService.GetModelElementsList(modelElementsListParamsApiDto)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, modelElements)
}

func (ph AdminModelsHandler) GetModelElement(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	modelElementId := c.Params.ByName("model-element-id")
	selectFields := strings.Split(c.Query("select"), ",")
	modelElement, err := ph.modelService.GetModelElement(modelCode, modelElementId, selectFields)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, modelElement)
}

func (ph AdminModelsHandler) CreateModelElement(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	var modelElementCreateApiDto models.ModelElementCreateApiDto
	err := c.ShouldBindWith(&modelElementCreateApiDto, binding.JSON)
	if err != nil {
		err = appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("some of input params is invalid. err: %s", err))}
		ph.errorService.HandleError(err, c)
		return
	}
	modelElementId, err := ph.modelService.CreateModelElement(modelCode, modelElementCreateApiDto)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, modelElementId)
}

func (ph AdminModelsHandler) EditModelElement(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	modelElementId := c.Params.ByName("model-element-id")
	var modelElementCreateApiDto models.ModelElementCreateApiDto
	err := c.ShouldBindWith(&modelElementCreateApiDto, binding.JSON)
	if err != nil {
		err = appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("some of input params is invalid. err: %s", err))}
		ph.errorService.HandleError(err, c)
		return
	}
	fmt.Println("service EditModelElement st")
	modelElementIdDto, err := ph.modelService.EditModelElement(modelCode, modelElementId, modelElementCreateApiDto)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	fmt.Println("service EditModelElement end")
	c.JSON(http.StatusOK, modelElementIdDto)
}

func (ph AdminModelsHandler) DeleteModelElement(c *gin.Context) {
	modelCode := c.Params.ByName("model-code")
	modelElementId := c.Params.ByName("model-element-id")
	err := ph.modelService.DeleteModelElement(modelCode, modelElementId)
	if err != nil {
		ph.errorService.HandleError(err, c)
		return
	}
	c.JSON(http.StatusOK, nil)
}
