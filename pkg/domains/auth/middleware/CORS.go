package middleware

import "github.com/gin-gonic/gin"

func HandleCORS(c *gin.Context) {
	c.Header("Referrer-Policy", "no-referrer-when-downgrade")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}
