package main

import (
	"school-information-system/config"
	"school-information-system/database"
	"school-information-system/database/seeds"
	"school-information-system/internal/cron"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/repos"
	"school-information-system/internal/routes"

	"github.com/gin-gonic/gin"
)

// INFO: i just put some test routes here for quick testing of auth service
func main() {
	// check env
	if err := config.EnvCheck(); err != nil {
		panic(err)
	}

	// setup database
	db := database.Connect()
	seeds.Migrate(db)

	// setup validator
	validatorlib.RegisterTagName()

	// setup router
	r := gin.Default()
	api := r.Group("/api")
	router := routes.NewRoute(db, repos.New(db))

	// register routes
	{
		router.RegisterBase(r)
		router.RegisterAuth(api.Group("/auth"))
		router.RegisterAdmin(api.Group("/admins"))
	}

	// start cron jobs
	cronjob := cron.New(db)
	cronjob.Start()
	defer cronjob.Stop()

	// port is automatically set to PORT env by gin
	r.Run()
}
