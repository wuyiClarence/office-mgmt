package routers

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"platform-backend/config"
	"platform-backend/middleware"
	"platform-backend/utils/customvalidator"
	"platform-backend/utils/format"
)

var Router *router

func init() {
	Router = NewRouter()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("checkTime", customvalidator.CheckTime)
		v.RegisterValidation("checkDate", customvalidator.CheckDate)
	}
}

func NewRouter() *router {
	return &router{}
}

type router struct{}

// Init
// @title Office-Platform API
// @version 1.0
// @description This is an example API documentation
// @host localhost:8080
// @BasePath /
func (router *router) Init() *gin.Engine {
	r := gin.New()

	gin.SetMode(config.MyConfig.AppConfig.Mod)
	if gin.IsDebugging() {
		pprof.Register(r, "/debug/pprof")
	}

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.Recover())
	r.Use(middleware.CrossMiddleware())

	r.NoRoute(routeNotFound)
	r.NoMethod(methodNotFound)

	r.StaticFS("/public", http.Dir("public"))

	routerV1(r.Group("/api"))
	return r
}

func methodNotFound(context *gin.Context) {
	format.NewResponseJson(context).ErrorWithHttpCode(http.StatusNotFound, http.StatusNotFound)
}

func routeNotFound(context *gin.Context) {
	format.NewResponseJson(context).ErrorWithHttpCode(http.StatusNotFound, http.StatusNotFound, "Not Found")
}
