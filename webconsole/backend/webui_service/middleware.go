package webui_service

import (
	"github.com/gin-gonic/gin"
	"free5gc/lib/path_util"
)

var PublicPath string

func init() {
	PublicPath = path_util.Gofree5gcPath("free5gc/webconsole/public")
}

func ReturnPublic() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		if method == "GET" {
			destPath := PublicPath + context.Request.RequestURI
			if destPath[len(destPath)-1] == '/' {
				destPath = destPath[:len(destPath)-1]
			}
			context.File(destPath)
		} else {
			context.Next()
		}
	}
}
