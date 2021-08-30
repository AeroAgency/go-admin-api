package service

import "github.com/gin-gonic/gin"

type Error interface {
	HandleError(err error, c *gin.Context)
}
