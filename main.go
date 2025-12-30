package main

import (
	_ "school-information-system/config"
	"school-information-system/database"
	"school-information-system/database/seeds"
	_ "school-information-system/docs"
	"school-information-system/internal/cron"
	_ "school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/repos"
	"school-information-system/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title			School Information System API
// @description	This is an API used for manages school environment.
// @version 1.0
// @host		localhost:8080
// @BasePath	/api
func main() {
	// setup database
	db := database.Connect()
	seeds.Migrate(db)

	// setup router
	r := gin.Default()
	api := r.Group("/api")
	router := routes.NewRoute(db, repos.New(db))

	// register routes
	{
		router.RegisterBase(r)
		router.RegisterAuth(api.Group("/auth"))
		router.RegisterAdmin(api.Group("/admins"))
		router.RegisterSubject(api.Group("/subjects"))
		router.RegisterClass(api.Group("/classes"))
	}

	// start cron jobs
	cronjob := cron.New(db)
	cronjob.Start()
	defer cronjob.Stop()

	// port is automatically set to PORT env by gin
	r.Run()
}
