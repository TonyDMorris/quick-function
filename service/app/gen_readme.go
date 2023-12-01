package app

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func (a *App) HandleStrapiWebhook(c *gin.Context) {
	bytes, _ := io.ReadAll(c.Request.Body)
	fmt.Println(string(bytes))

	c.JSON(200, gin.H{
		"success": "true", "message": "ok",
	})
}
