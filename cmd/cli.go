package cmd

import (
	"fmt"
	"log"

	"github.com/alexdor/github-user-interaction-fetcher/app/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	version bool

	Cmd = &cobra.Command{
		Use:   "",
		Short: "Fetch contributions links from github users",

		Run: func(ccmd *cobra.Command, args []string) {
			if version {
				fmt.Println("v1.0")
			} else {
				start()
			}
		},
	}
)

func start() {
	routes.Init()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.SetupRoutes(router)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
