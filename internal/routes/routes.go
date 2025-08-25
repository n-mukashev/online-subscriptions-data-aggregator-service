package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mukashev-n/online-subscriptions-data-aggregator-service/docs"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/validators"
)

func RegisterRoutes(server *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterValidators(v)
	}
	following := server.Group("/subscription")
	{
		following.GET("/:id", getById)
		following.GET("/all", getAll)
		following.POST("", create)
		following.PUT("", update)
		following.DELETE("/:id", delete)
		following.POST("/invoice", getSubscriptionsInvoice)
	}
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
