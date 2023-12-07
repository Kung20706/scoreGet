package controller

import (
	"net/http"
	"test/model"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @Summary Get a list of items
// @Description Get a list of items
// @Tags items
// @Accept json
// @Produce json
// @Param name query string true "Item Name"
// @Success 200 {object} string "OK"
// @Router /items [get]
func GetItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var tickets []model.TicketNumber
	db.Find(&tickets)

	c.JSON(http.StatusOK, tickets)

}

// 添加其他路由和处理程序...

// Swagger 文档路由
func SwaggerHandler(c *gin.Context) {
	ginSwagger.WrapHandler(swaggerfiles.Handler)(c)
}
