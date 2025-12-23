package routes

import (
	"fmt"
	"school-information-system/config"
	"school-information-system/internal/libs/replylib"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Route struct {
	db *gorm.DB
}

func NewRoute(db *gorm.DB) *Route {
	return &Route{db}
}

func (rt *Route) RegisterBase(r *gin.Engine) {
	r.GET("/api/health", func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))
		response := map[string]any{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		dbcheck := c.Query("dbcheck")
		if dbcheck == "true" {
			db, err := rt.db.DB()
			if err != nil {
				details := config.SplitByEnv("", err.Error())
				rp.Error(replylib.CodeServerError, "Database connection error", reply.OptErrorPayload{Details: details}).FailJSON()
				return
			}
			err = db.Ping()
			if err != nil {
				details := config.SplitByEnv("", err.Error())
				rp.Error(replylib.CodeServerError, "Database ping failed", reply.OptErrorPayload{Details: details}).FailJSON()
				return
			}
			response["database"] = "healthy"
		}

		rp.Success(response).OkJSON()
	})

	r.NoRoute(func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))
		msg := fmt.Sprintf("can not %s %s", c.Request.Method, c.Request.URL.Path)
		rp.Error(replylib.CodeNotFound, msg).FailJSON()
	})
}
