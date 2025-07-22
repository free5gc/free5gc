package context

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" // #nosec G505
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	gtpv1 "github.com/wmnsk/go-gtp/gtpv1"

	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/pkg/factory"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/sctp"
	"github.com/free5gc/util/idgenerator"
	"github.com/free5gc/util/ippool"
)

type n3iwf interface {
	Config() *factory.Config
	CancelContext() context.Context
}

type N3IWFContext struct {
	n3iwf

	// ID generator
	RANUENGAPIDGenerator *idgenerator.IDGenerator
	TEIDGenerator        *idgenerator.IDGenerator

	// Pools
	AMFPool                sync.Map // map[string]*N3IWFAMF, SCTPAddr as key
	AMFReInitAvailableList sync.Map // map[string]bool, SCTPAddr as key
	IKESA                  sync.Map // map[uint64]*IKESecurityAssociation, SPI as key
	ChildSA                sync.Map // map[uint32]*ChildSecurityAssociation, inboundSPI as key
	GTPConnectionWithUPF   sync.Map // map[string]*gtpv1.UPlaneConn, UPF address as key
	AllocatedUEIPAddress   sync.Map // map[string]*N3IWFIkeUe, IPAddr as key
	AllocatedUETEID        sync.Map // map[uint32]*RanUe, TEID as key
	IKEUePool              sync.Map // map[uint64]*N3IWFIkeUe, SPI as key
	RANUePool              sync.Map // map[int64]*RanUe, RanUeNgapID as key
	IKESPIToNGAPId         sync.Map // map[uint64]RanUeNgapID, SPI as key
	NGAPIdToIKESPI         sync.Map // map[uint64]SPI, RanUeNgapID as key

	// Security data
	CertificateAuthority []byte
	N3IWFCertificate     []byte
	N3IWFPrivateKey      *rsa.PrivateKey

	IPSecInnerIPPool *ippool.IPPool
	// TODO: [TWIF] TwifUe may has its own IP address pool

	// XFRM interface
	XfrmIfaces          sync.Map // map[uint32]*netlink.Link, XfrmIfaceId as key
	XfrmParentIfaceName string
	// Every UE's first UP IPsec will use default XFRM interface, additoinal UP IPsec will offset its XFRM id
	XfrmIfaceIdOffsetForUP uint32
}

func NewContext(n3iwf n3iwf) (*N3IWFContext, error) {
	n := &N3IWFContext{
		n3iwf:                n3iwf,
		RANUENGAPIDGenerator: idgenerator.NewGenerator(0, math.MaxInt64),
		TEIDGenerator:        idgenerator.NewGenerator(1, math.MaxUint32),
	}
	cfg := n3iwf.Config()

	// Private key
	block, _, err := decodePEM(cfg.GetIKECertKeyPath())
	if err != nil {
		return nil, errors.Wrapf(err, "IKE PrivKey")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.CtxLog.Warnf("Parse PKCS8 private key failed: %v", err)
		logger.CtxLog.Info("Parse using PKCS1...")

		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Errorf("Parse PKCS1 pricate key failed: %v", err)
		}
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.Errorf("Private key is not an rsa private key")
	}
	n.N3IWFPrivateKey = rsaKey

	// Certificate authority
	block, _, err = decodePEM(cfg.GetIKECAPemPath())
	if err != nil {
		return nil, errors.Wrapf(err, "IKE CA")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Errorf("Parse certificate authority failed: %v", err)
	}
	// Get sha1 hash of subject public key info
	sha1Hash := sha1.New() // #nosec G401
	_, err = sha1Hash.Write(cert.RawSubjectPublicKeyInfo)
	if err != nil {
		return nil, errors.Errorf("Hash function writing failed: %+v", err)
	}
	n.CertificateAuthority = sha1Hash.Sum(nil)

	// Certificate
	block, _, err = decodePEM(cfg.GetIKECertPemPath())
	if err != nil {
		return nil, errors.Wrapf(err, "IKE Cert")
	}
	n.N3IWFCertificate = block.Bytes

	// UE IP address range
	ueIPPool, err := ippool.NewIPPool(cfg.GetUEIPAddrRange())
	if err != nil {
		return nil, errors.Errorf("NewContext(): %+v", err)
	}
	n.IPSecInnerIPPool = ueIPPool

	// XFRM related
	ikeBindIfaceName, err := getInterfaceName(cfg.GetIKEBindAddr())
	if err != nil {
		return nil, err
	}
	n.XfrmParentIfaceName = ikeBindIfaceName

	return n, nil
}

func decodePEM(path string) (*pem.Block, []byte, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Cannot read file(%s)", path)
	}
	p, rest := pem.Decode(content)
	if p == nil {
		return nil, nil, errors.Errorf("Decode pem failed")
	}
	return p, rest, nil
}

func getInterfaceName(IPAddress string) (interfaceName string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "nil", err
	}

	res, err := net.ResolveIPAddr("ip4", IPAddress)
	if err != nil {
		return "", fmt.Errorf("Error resolving address [%s]: %v", IPAddress, err)
	}
	IPAddress = res.String()

	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			return "nil", err
		}
		for _, addr := range addrs {
			if IPAddress == addr.String()[0:strings.Index(addr.String(), "/")] {
				return inter.Name, nil
			}
		}
	}
	return "", fmt.Errorf("cannot find interface name for IP[%s]", IPAddress)
}

func (c *N3IWFContext) NewN3iwfIkeUe(spi uint64) *N3IWFIkeUe {
	n3iwfIkeUe := &N3IWFIkeUe{
		N3iwfCtx: c,
	}
	n3iwfIkeUe.init()
	c.IKEUePool.Store(spi, n3iwfIkeUe)
	return n3iwfIkeUe
}

func (c *N3IWFContext) NewN3iwfRanUe() *N3IWFRanUe {
	ranUeNgapId, err := c.RANUENGAPIDGenerator.Allocate()
	if err != nil {
		logger.CtxLog.Errorf("New N3IWF UE failed: %+v", err)
		return nil
	}
	n3iwfRanUe := &N3IWFRanUe{
		RanUeSharedCtx: RanUeSharedCtx{
			N3iwfCtx: c,
		},
	}
	n3iwfRanUe.init(ranUeNgapId)
	c.RANUePool.Store(ranUeNgapId, n3iwfRanUe)

	return n3iwfRanUe
}

func (c *N3IWFContext) DeleteRanUe(ranUeNgapId int64) {
	c.RANUePool.Delete(ranUeNgapId)
	c.DeleteIkeSPIFromNgapId(ranUeNgapId)
}

func (c *N3IWFContext) DeleteIKEUe(spi uint64) {
	c.IKEUePool.Delete(spi)
	c.DeleteNgapIdFromIkeSPI(spi)
}

func (c *N3IWFContext) IkeUePoolLoad(spi uint64) (*N3IWFIkeUe, bool) {
	ikeUe, ok := c.IKEUePool.Load(spi)
	if ok {
		return ikeUe.(*N3IWFIkeUe), ok
	} else {
		return nil, ok
	}
}

func (c *N3IWFContext) RanUePoolLoad(id interface{}) (RanUe, bool) {
	var ranUeNgapId int64

	cfgLog := logger.CfgLog
	switch id := id.(type) {
	case int64:
		ranUeNgapId = id
	default:
		cfgLog.Warnf("RanUePoolLoad unhandle type: %t", id)
		return nil, false
	}

	ranUe, ok := c.RANUePool.Load(ranUeNgapId)
	if ok {
		return ranUe.(RanUe), ok
	} else {
		return nil, ok
	}
}

func (c *N3IWFContext) IkeSpiNgapIdMapping(spi uint64, ranUeNgapId int64) {
	c.IKESPIToNGAPId.Store(spi, ranUeNgapId)
	c.NGAPIdToIKESPI.Store(ranUeNgapId, spi)
}

func (c *N3IWFContext) IkeSpiLoad(ranUeNgapId int64) (uint64, bool) {
	spi, ok := c.NGAPIdToIKESPI.Load(ranUeNgapId)
	if ok {
		return spi.(uint64), ok
	}
	return 0, false
}

func (c *N3IWFContext) NgapIdLoad(spi uint64) (int64, bool) {
	ranNgapId, ok := c.IKESPIToNGAPId.Load(spi)
	if ok {
		return ranNgapId.(int64), ok
	}
	return 0, false
}

func (c *N3IWFContext) DeleteNgapIdFromIkeSPI(spi uint64) {
	c.IKESPIToNGAPId.Delete(spi)
}

func (c *N3IWFContext) DeleteIkeSPIFromNgapId(ranUeNgapId int64) {
	c.NGAPIdToIKESPI.Delete(ranUeNgapId)
}

func (c *N3IWFContext) RanUeLoadFromIkeSPI(spi uint64) (RanUe, error) {
	ranNgapId, ok := c.IKESPIToNGAPId.Load(spi)
	if ok {
		ranUe, err := c.RanUePoolLoad(ranNgapId.(int64))
		if !err {
			return nil, fmt.Errorf("cannot find RanUE from RanNgapId : %+v", ranNgapId)
		}
		return ranUe, nil
	}
	return nil, fmt.Errorf("cannot find RanNgapId from IkeUe SPI : %+v", spi)
}

func (c *N3IWFContext) IkeUeLoadFromNgapId(ranUeNgapId int64) (*N3IWFIkeUe, error) {
	spi, ok := c.NGAPIdToIKESPI.Load(ranUeNgapId)
	if ok {
		ikeUe, err := c.IkeUePoolLoad(spi.(uint64))
		if !err {
			return nil, fmt.Errorf("cannot find IkeUe from spi : %+v", spi)
		}
		return ikeUe, nil
	}
	return nil, fmt.Errorf("cannot find SPI from NgapId : %+v", ranUeNgapId)
}

func (c *N3IWFContext) NewN3iwfAmf(sctpAddr string, conn *sctp.SCTPConn) *N3IWFAMF {
	amf := new(N3IWFAMF)
	amf.init(sctpAddr, conn)
	item, loaded := c.AMFPool.LoadOrStore(sctpAddr, amf)
	if loaded {
		logger.CtxLog.Warn("[Context] NewN3iwfAmf(): AMF entry already exists.")
		return item.(*N3IWFAMF)
	}
	return amf
}

func (c *N3IWFContext) DeleteN3iwfAmf(sctpAddr string) {
	c.AMFPool.Delete(sctpAddr)
}

func (c *N3IWFContext) AMFPoolLoad(sctpAddr string) (*N3IWFAMF, bool) {
	amf, ok := c.AMFPool.Load(sctpAddr)
	if ok {
		return amf.(*N3IWFAMF), ok
	}
	return nil, false
}

func (c *N3IWFContext) DeleteAMFReInitAvailableFlag(sctpAddr string) {
	c.AMFReInitAvailableList.Delete(sctpAddr)
}

func (c *N3IWFContext) AMFReInitAvailableListLoad(sctpAddr string) (bool, bool) {
	flag, ok := c.AMFReInitAvailableList.Load(sctpAddr)
	if ok {
		return flag.(bool), ok
	}
	return true, false
}

func (c *N3IWFContext) AMFReInitAvailableListStore(sctpAddr string, flag bool) {
	c.AMFReInitAvailableList.Store(sctpAddr, flag)
}

func (c *N3IWFContext) NewIKESecurityAssociation() *IKESecurityAssociation {
	ikeSecurityAssociation := new(IKESecurityAssociation)

	maxSPI := new(big.Int).SetUint64(math.MaxUint64)
	var localSPIuint64 uint64

	for {
		localSPI, err := rand.Int(rand.Reader, maxSPI)
		if err != nil {
			logger.CtxLog.Error("[Context] Error occurs when generate new IKE SPI")
			return nil
		}
		localSPIuint64 = localSPI.Uint64()
		_, duplicate := c.IKESA.LoadOrStore(localSPIuint64, ikeSecurityAssociation)
		if !duplicate {
			break
		}
	}

	ikeSecurityAssociation.LocalSPI = localSPIuint64

	return ikeSecurityAssociation
}

func (c *N3IWFContext) DeleteIKESecurityAssociation(spi uint64) {
	c.IKESA.Delete(spi)
}

func (c *N3IWFContext) IKESALoad(spi uint64) (*IKESecurityAssociation, bool) {
	securityAssociation, ok := c.IKESA.Load(spi)
	if ok {
		return securityAssociation.(*IKESecurityAssociation), ok
	}
	return nil, false
}

func (c *N3IWFContext) DeleteGTPConnection(upfAddr string) {
	c.GTPConnectionWithUPF.Delete(upfAddr)
}

func (c *N3IWFContext) GTPConnectionWithUPFLoad(upfAddr string) (*gtpv1.UPlaneConn, bool) {
	conn, ok := c.GTPConnectionWithUPF.Load(upfAddr)
	if ok {
		return conn.(*gtpv1.UPlaneConn), ok
	}
	return nil, false
}

func (c *N3IWFContext) GTPConnectionWithUPFStore(upfAddr string, conn *gtpv1.UPlaneConn) {
	c.GTPConnectionWithUPF.Store(upfAddr, conn)
}

func (c *N3IWFContext) NewIPsecInnerUEIP(ikeUe *N3IWFIkeUe) (net.IP, error) {
	var ueIPAddr net.IP
	var err error
	cfg := c.Config()
	ipsecGwAddr := cfg.GetIPSecGatewayAddr()

	for {
		ueIPAddr, err = c.IPSecInnerIPPool.Allocate(nil)
		if err != nil {
			return nil, errors.Wrapf(err, "NewIPsecInnerUEIP()")
		}
		if ueIPAddr.String() == ipsecGwAddr {
			continue
		}
		_, ok := c.AllocatedUEIPAddress.LoadOrStore(ueIPAddr.String(), ikeUe)
		if ok {
			logger.CtxLog.Warnf("NewIPsecInnerUEIP(): IP(%v) is used by other IkeUE",
				ueIPAddr.String())
		} else {
			break
		}
	}

	return ueIPAddr, nil
}

func (c *N3IWFContext) DeleteInternalUEIPAddr(ipAddr string) {
	c.AllocatedUEIPAddress.Delete(ipAddr)
}

func (c *N3IWFContext) AllocatedUEIPAddressLoad(ipAddr string) (*N3IWFIkeUe, bool) {
	ikeUe, ok := c.AllocatedUEIPAddress.Load(ipAddr)
	if ok {
		return ikeUe.(*N3IWFIkeUe), ok
	}
	return nil, false
}

func (c *N3IWFContext) NewTEID(ranUe RanUe) uint32 {
	teid64, err := c.TEIDGenerator.Allocate()
	if err != nil {
		logger.CtxLog.Errorf("New TEID failed: %+v", err)
		return 0
	}
	if teid64 < 0 || teid64 > math.MaxUint32 {
		logger.CtxLog.Warnf("NewTEID teid64 out of uint32 range: %d, use maxUint32", teid64)
		return 0
	}
	teid32 := uint32(teid64)

	c.AllocatedUETEID.Store(teid32, ranUe)

	return teid32
}

func (c *N3IWFContext) DeleteTEID(teid uint32) {
	c.TEIDGenerator.FreeID(int64(teid))
	c.AllocatedUETEID.Delete(teid)
}

func (c *N3IWFContext) AllocatedUETEIDLoad(teid uint32) (RanUe, bool) {
	ranUe, ok := c.AllocatedUETEID.Load(teid)
	if ok {
		return ranUe.(RanUe), ok
	}
	return nil, false
}

func (c *N3IWFContext) AMFSelection(
	ueSpecifiedGUAMI *ngapType.GUAMI,
	ueSpecifiedPLMNId *ngapType.PLMNIdentity,
) *N3IWFAMF {
	var availableAMF, defaultAMF *N3IWFAMF
	c.AMFPool.Range(func(key, value interface{}) bool {
		amf := value.(*N3IWFAMF)
		if defaultAMF == nil {
			defaultAMF = amf
		}
		if amf.FindAvalibleAMFByCompareGUAMI(ueSpecifiedGUAMI) {
			availableAMF = amf
			return false
		} else {
			// Fail to find through GUAMI served by UE.
			// Try again using SelectedPLMNId
			if amf.FindAvalibleAMFByCompareSelectedPLMNId(ueSpecifiedPLMNId) {
				availableAMF = amf
				return false
			}
			return true
		}
	})
	if availableAMF == nil &&
		defaultAMF != nil {
		availableAMF = defaultAMF
	}
	return availableAMF
}
