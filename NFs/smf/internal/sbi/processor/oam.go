package processor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

type PDUSessionInfo struct {
	Supi         string
	PDUSessionID string
	Dnn          string
	Sst          string
	Sd           string
	AnType       models.AccessType
	PDUAddress   string
	SessionRule  models.SessionRule
	UpCnxState   models.UpCnxState
	Tunnel       context.UPTunnel
}

func (p *Processor) HandleOAMGetUEPDUSessionInfo(c *gin.Context, smContextRef string) {
	smContext := context.GetSMContextByRef(smContextRef)
	if smContext == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	pduSessionInfo := &PDUSessionInfo{
		Supi:         smContext.Supi,
		PDUSessionID: strconv.Itoa(int(smContext.PDUSessionID)),
		Dnn:          smContext.Dnn,
		Sst:          strconv.Itoa(int(smContext.SNssai.Sst)),
		Sd:           smContext.SNssai.Sd,
		AnType:       smContext.AnType,
		PDUAddress:   smContext.PDUAddress.String(),
		UpCnxState:   smContext.UpCnxState,
		// Tunnel: context.UPTunnel{
		// 	//UpfRoot:  smContext.Tunnel.UpfRoot,
		// 	ULCLRoot: smContext.Tunnel.UpfRoot,
		// },
	}
	c.JSON(http.StatusOK, pduSessionInfo)
}

func (p *Processor) HandleGetSMFUserPlaneInfo(c *gin.Context) {
	c.JSON(http.StatusOK, factory.SmfConfig.Configuration.UserPlaneInformation)
}
