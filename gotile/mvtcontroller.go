package gotile

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func errHandler(err error, context *gin.Context) {
	if err != nil {
		context.Status(http.StatusNotFound)
		context.Error(err)
	}
}
func GetTileMvt(context *gin.Context) {
	x, err := strconv.ParseInt(context.Param("x"), 10, 64)
	if err != nil {
		errHandler(err, context)
		return
	}
	y, err := strconv.ParseInt(context.Param("y"), 10, 64)
	if err != nil {
		errHandler(err, context)
		return
	}
	z, err := strconv.ParseInt(context.Param("z"), 10, 64)
	if err != nil {
		errHandler(err, context)
		return
	}
	format := context.Param("format")
	tablename := context.Param("tablename")
	pbf, err := GetMvtTileBinary(x, y, z, tablename, format)
	if err != nil {
		errHandler(err, context)
		return
	}
	context.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%d-%d-%d.%s", z, x, y, format))
	context.Data(202, "application/x-protobuf", pbf)
}

func RegisterRouter(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.GET("/service/:tablename/:x/:y/:z/:format", GetTileMvt)
}
