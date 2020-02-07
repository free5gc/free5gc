package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/nrf/logger"
	. "github.com/free5gc/openapi/models"
	"github.com/free5gc/path_util"
)

var (
	NrfLogPath = path_util.Free5gcPath("github.com/free5gc/nrf/management/sslkeylog.log")
	NrfPemPath = path_util.Free5gcPath("free5gc/support/TLS/nrf.pem")
	NrfKeyPath = path_util.Free5gcPath("free5gc/support/TLS/nrf.key")
)

func main() {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	router.POST("", func(c *gin.Context) {
		/*buf, err := c.GetRawData()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Remove NL line feed, new line character
		//requestBody = string(buf[:len(b uf)-1])*/
		var ND NotificationData

		if err := c.ShouldBindJSON(&ND); err != nil {
			log.Panic(err.Error())
		}
		fmt.Println(ND)
		c.JSON(http.StatusNoContent, gin.H{})
	})

	srv, err := http2_util.NewServer(":30678", NrfLogPath, router)
	if err != nil {
		log.Panic(err.Error())
	}

	err2 := srv.ListenAndServeTLS(NrfPemPath, NrfKeyPath)
	if err2 != nil && err2 != http.ErrServerClosed {
		log.Panic(err2.Error())
	}
}
