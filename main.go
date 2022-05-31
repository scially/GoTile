package main

import (
	"GoTile/config"
	"GoTile/gotile"
	"github.com/gin-contrib/gzip"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	if config.Configiure.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	gotile.RegisterRouter(router)

	var port int = config.Configiure.Server.Port
	if port == 0 {
		port = 9090
	}
	log.Printf("server listen on 0.0.0.0:%d", port)
	router.Run("0.0.0.0:" + strconv.Itoa(port))
}
