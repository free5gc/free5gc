package PDUSession

import (
	"github.com/gin-gonic/gin"
	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/smf/smf_pfcp/pfcp_udp"
	"log"
)

func DummyServer() {
	router := gin.Default()

	AddService(router)

	go pfcp_udp.Run()

	smfKeyLogPath := path_util.Gofree5gcPath("free5gc/smfsslkey.log")
	smfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/smf.pem")
	smfkeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/smf.key")

	server, _ := http2_util.NewServer(":29502", smfKeyLogPath, router)

	err := server.ListenAndServeTLS(smfPemPath, smfkeyPath)

	if err != nil {
		log.Fatal(err)
	}
}
