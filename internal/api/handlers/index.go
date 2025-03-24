package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexHandler struct{}

func (h *IndexHandler) RootHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, "<h1>Hello World</h1>")
}