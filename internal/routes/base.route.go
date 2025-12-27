package routes

import (
	"fmt"
	"net/http"
	"school-information-system/config"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/repos"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Route struct {
	db *gorm.DB
	rp *repos.Repos
}

func NewRoute(db *gorm.DB, repositories *repos.Repos) *Route {
	return &Route{db, repositories}
}

type health struct {
	Status    string `json:"status" example:"healthy"`
	Database  string `json:"database,omitempty" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2006-01-02T15:04:05Z07:00"`
}

// @Summary      Check server's health
// @Tags         health
// @Accept       json
// @Produce      json
// @Param				 dbcheck  query 		bool 	false	"set 'true' if want to check db status"
// @Success      200  		{object}  swaglib.Envelope{data=health}
// @Response     default  {object}  swaglib.Envelope{data=reply.ErrorPayload}
// @Router       /health 	[get]
func (rt *Route) RegisterBase(r *gin.Engine) {
	r.GET("/api/health", func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))
		response := health{Status: "healthy", Timestamp: time.Now().Format(time.RFC3339)}

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
			response.Database = "healthy"
		}

		rp.Success(response).OkJSON()
	})

	r.NoRoute(func(c *gin.Context) {
		rp := replylib.Client.New(adapter.AdaptGin(c))
		msg := fmt.Sprintf("can not %s %s", c.Request.Method, c.Request.URL.Path)
		rp.Error(replylib.CodeNotFound, msg).FailJSON()
	})

	r.GET("/swagger", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/swaggerview/index.html")
	})
	r.GET("/swaggerview/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
