package context

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fiorix/go-diameter/diam/sm"

	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/free5gc/util/idgenerator"
)

var chfContext CHFContext

func Init() {
	InitChfContext(&chfContext)
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &CHFContext{}

type CHFContext struct {
	NfId                      string
	Name                      string
	Url                       string
	UriScheme                 models.UriScheme
	BindingIPv4               string
	RegisterIPv4              string
	SBIPort                   int
	NfService                 map[models.ServiceName]models.NrfNfManagementNfService
	RecordSequenceNumber      map[string]int64
	LocalRecordSequenceNumber uint64
	NrfUri                    string
	NrfCertPem                string
	UePool                    sync.Map
	OAuth2Required            bool

	RatingCfg *sm.Settings
	AbmfCfg   *sm.Settings

	RatingSessionIdGenerator  *idgenerator.IDGenerator
	AccountSessionIdGenerator *idgenerator.IDGenerator
	sync.Mutex
}

func (c *CHFContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.UtilLog.Debugf("CHFContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}

	logger.UtilLog.Debugf("CHFContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.NrfCertPem)
}

func (context *CHFContext) AddChfUeToUePool(ue *ChfUe, supi string) {
	if len(supi) == 0 {
		logger.CtxLog.Errorf("Supi is nil")
	}
	ue.Supi = supi
	context.UePool.Store(ue.Supi, ue)
}

// Allocate CHF Ue with supi and add to chf Context and returns allocated ue
func (context *CHFContext) NewCHFUe(supi string) (*ChfUe, error) {
	if ue, ok := context.ChfUeFindBySupi(supi); ok {
		return ue, nil
	}
	if strings.HasPrefix(supi, "imsi-") {
		ue := ChfUe{}
		ue.init()

		if supi != "" {
			context.AddChfUeToUePool(&ue, supi)
		}

		return &ue, nil
	} else {
		return nil, fmt.Errorf(" add Ue context fail ")
	}
}

func (context *CHFContext) ChfUeFindBySupi(supi string) (*ChfUe, bool) {
	if value, ok := context.UePool.Load(supi); ok {
		return value.(*ChfUe), ok
	}
	return nil, false
}

func GenerateRatingSessionId() uint32 {
	if id, err := chfContext.RatingSessionIdGenerator.Allocate(); err == nil {
		return uint32(id)
	}
	return 0
}

func GenerateAccountSessionId() uint32 {
	if id, err := chfContext.AccountSessionIdGenerator.Allocate(); err == nil {
		return uint32(id)
	}
	return 0
}

func GetSelf() *CHFContext {
	return &chfContext
}

func (c *CHFContext) GetSelfID() string {
	return c.NfId
}

func (c *CHFContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_CHF, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}
