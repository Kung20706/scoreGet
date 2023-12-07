package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InjectDB middleware injects the DB instance into the Gin context
func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
