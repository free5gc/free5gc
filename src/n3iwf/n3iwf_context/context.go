package n3iwf_context

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math"
	"math/big"
	"net"

	"github.com/sirupsen/logrus"
	gtpv1 "github.com/wmnsk/go-gtp/v1"
	"golang.org/x/net/ipv4"

	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
)

var contextLog *logrus.Entry

var n3iwfContext = N3IWFContext{}
var ranUeNgapIdGenerator int64 = 0
var teidGenerator uint32 = 1

type N3IWFContext struct {
	NFInfo                 N3IWFNFInfo
	UePool                 map[int64]*N3IWFUe                   // RanUeNgapID as key
	AMFPool                map[string]*N3IWFAMF                 // SCTPAddr as key
	AMFReInitAvailableList map[string]bool                      // SCTPAddr as key
	IKESA                  map[uint64]*IKESecurityAssociation   // SPI as key
	ChildSA                map[uint32]*ChildSecurityAssociation // SPI as key
	GTPConnectionWithUPF   map[string]*gtpv1.UPlaneConn         // UPF address as key
	AllocatedUEIPAddress   map[string]*N3IWFUe                  // IPAddr as key
	AllocatedUETEID        map[uint32]*N3IWFUe                  // TEID as key

	// N3IWF FQDN
	FQDN string

	// security data
	CertificateAuthority []byte
	N3IWFCertificate     []byte
	N3IWFPrivateKey      *rsa.PrivateKey

	// UEIPAddressRange
	Subnet *net.IPNet

	// Network interface mark for xfrm
	Mark uint32

	// N3IWF local address
	IKEBindAddress      string
	IPSecGatewayAddress string
	GTPBindAddress      string
	TCPPort             uint16

	// N3IWF N1 interface raw socket
	N1RawSocket *ipv4.RawConn
}

func init() {
	// init log
	contextLog = logger.ContextLog

	// init context
	N3IWFSelf().UePool = make(map[int64]*N3IWFUe)
	N3IWFSelf().AMFPool = make(map[string]*N3IWFAMF)
	N3IWFSelf().AMFReInitAvailableList = make(map[string]bool)
	N3IWFSelf().IKESA = make(map[uint64]*IKESecurityAssociation)
	N3IWFSelf().ChildSA = make(map[uint32]*ChildSecurityAssociation)
	N3IWFSelf().GTPConnectionWithUPF = make(map[string]*gtpv1.UPlaneConn)
	N3IWFSelf().AllocatedUEIPAddress = make(map[string]*N3IWFUe)
	N3IWFSelf().AllocatedUETEID = make(map[uint32]*N3IWFUe)
}

// Create new N3IWF context
func N3IWFSelf() *N3IWFContext {
	return &n3iwfContext
}

func (context *N3IWFContext) NewN3iwfUe() *N3IWFUe {
	n3iwfUe := &N3IWFUe{}
	n3iwfUe.init()

	ranUeNgapIdGenerator %= MaxValueOfRanUeNgapID
	ranUeNgapIdGenerator++
	for {
		if _, double := context.UePool[ranUeNgapIdGenerator]; double {
			ranUeNgapIdGenerator++
		} else {
			break
		}
	}

	n3iwfUe.RanUeNgapId = ranUeNgapIdGenerator
	n3iwfUe.AmfUeNgapId = AmfUeNgapIdUnspecified
	context.UePool[n3iwfUe.RanUeNgapId] = n3iwfUe
	return n3iwfUe
}

func (context *N3IWFContext) NewN3iwfAmf(sctpAddr string) *N3IWFAMF {
	if amf, ok := context.AMFPool[sctpAddr]; ok {
		contextLog.Warn("[Context] NewN3iwfAmf(): AMF entry already exists.")
		return amf
	} else {
		amf = &N3IWFAMF{
			SCTPAddr:              sctpAddr,
			N3iwfUeList:           make(map[int64]*N3IWFUe),
			AMFTNLAssociationList: make(map[string]*AMFTNLAssociationItem),
		}
		context.AMFPool[sctpAddr] = amf
		return amf
	}
}

func (context *N3IWFContext) FindAMFBySCTPAddr(sctpAddr string) (*N3IWFAMF, error) {
	amf, ok := context.AMFPool[sctpAddr]
	if !ok {
		return nil, fmt.Errorf("[Context] FindAMF(): AMF not found. sctpAddr: %s", sctpAddr)
	}
	return amf, nil
}

func (context *N3IWFContext) FindUeByRanUeNgapID(ranUeNgapID int64) *N3IWFUe {
	if n3iwfUE, ok := context.UePool[ranUeNgapID]; ok {
		return n3iwfUE
	} else {
		return nil
	}
}

// returns true means reinitialization is available, and false is unavailable.
func (context *N3IWFContext) CheckAMFReInit(sctpAddr string) bool {

	if check, ok := context.AMFReInitAvailableList[sctpAddr]; ok {
		return check
	}
	return true
}

func (context *N3IWFContext) NewIKESecurityAssociation() *IKESecurityAssociation {
	ikeSecurityAssociation := &IKESecurityAssociation{}

	var maxSPI *big.Int = new(big.Int).SetUint64(math.MaxUint64)
	var localSPIuint64 uint64

	for {
		localSPI, err := rand.Int(rand.Reader, maxSPI)
		if err != nil {
			contextLog.Error("[Context] Error occurs when generate new IKE SPI")
			return nil
		}
		localSPIuint64 = localSPI.Uint64()
		if _, duplicate := context.IKESA[localSPIuint64]; !duplicate {
			break
		}
	}

	ikeSecurityAssociation.LocalSPI = localSPIuint64
	context.IKESA[localSPIuint64] = ikeSecurityAssociation

	return ikeSecurityAssociation
}

func (context *N3IWFContext) NewTEID(ue *N3IWFUe) uint32 {
	for {
		if teidGenerator == 0 {
			teidGenerator++
			continue
		}
		if _, double := context.AllocatedUETEID[teidGenerator]; double {
			teidGenerator++
		} else {
			break
		}
	}

	context.AllocatedUETEID[teidGenerator] = ue

	return teidGenerator
}

func (context *N3IWFContext) FindIKESecurityAssociationBySPI(spi uint64) *IKESecurityAssociation {
	if ikeSecurityAssociation, ok := context.IKESA[spi]; ok {
		return ikeSecurityAssociation
	} else {
		return nil
	}
}

func (context *N3IWFContext) AMFSelection(ueSpecifiedGUAMI *ngapType.GUAMI) *N3IWFAMF {
	for _, n3iwfAMF := range context.AMFPool {
		if n3iwfAMF.FindAvalibleAMFByCompareGUAMI(ueSpecifiedGUAMI) {
			return n3iwfAMF
		}
	}
	return nil
}
