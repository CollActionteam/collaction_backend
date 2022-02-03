package http

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())
	return r
}
