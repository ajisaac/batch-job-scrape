package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"scrapebatch-go/api/database"
	"scrapebatch-go/api/middleware"
	"scrapebatch-go/api/route"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var server = Server{}

func Run() {
	server.Router = gin.Default()
	server.Router.Use(middleware.CORSMiddleware())

	server.DB = database.InitDatabase()
	defer database.Close()

	route.InitializeRoutes(server.Router)
	_ = server.Router.Run()

}
