package routes

import (
	"fmt"
	"os"

	"github.com/alexdor/github-user-interaction-fetcher/app/controllers"
	"github.com/gin-gonic/gin"
)

func Init() {
	err := controllers.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetupRoutes(router *gin.Engine) {
	// router.Use(static.Serve("/", static.LocalFile("/app/static", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("app/static/index.html")
	})
	router.POST("/api/v1/userInfo", controllers.GetUserInfo)

}
