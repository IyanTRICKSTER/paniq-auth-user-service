package middleware

import "github.com/gin-gonic/gin"

func HandleCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}
