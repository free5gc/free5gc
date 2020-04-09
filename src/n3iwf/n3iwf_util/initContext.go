package n3iwf_util

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net"
	"strings"

	"github.com/sirupsen/logrus"

	"free5gc/lib/path_util"
	"free5gc/src/n3iwf/factory"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
)

var contextLog *logrus.Entry

func init() {
	contextLog = logger.ContextLog
}

func InitN3IWFContext() bool {
	var ok bool

	if factory.N3iwfConfig.Configuration == nil {
		contextLog.Error("No N3IWF configuration found")
		return false
	}

	n3iwfContext := n3iwf_context.N3IWFSelf()
	n3iwfContext.NFInfo = factory.N3iwfConfig.Configuration.N3IWFInfo

	if ok = formatSupportedTAList(&n3iwfContext.NFInfo); !ok {
		return false
	}

	// IKE bind address
	if factory.N3iwfConfig.Configuration.IKEBindAddr == "" {
		contextLog.Error("IKE bind address is empty")
		return false
	} else {
		n3iwfContext.IKEBindAddress = factory.N3iwfConfig.Configuration.IKEBindAddr
	}

	// IPSec gateway address
	if factory.N3iwfConfig.Configuration.IPSecGatewayAddr == "" {
		contextLog.Error("IPSec interface address is empty")
		return false
	} else {
		n3iwfContext.IPSecGatewayAddress = factory.N3iwfConfig.Configuration.IPSecGatewayAddr
	}

	// GTP bind address
	if factory.N3iwfConfig.Configuration.GTPBindAddr == "" {
		contextLog.Error("GTP bind address is empty")
		return false
	} else {
		n3iwfContext.GTPBindAddress = factory.N3iwfConfig.Configuration.GTPBindAddr
	}

	// TCP port
	if factory.N3iwfConfig.Configuration.TCPPort == 0 {
		contextLog.Error("TCP port is not defined")
		return false
	} else {
		n3iwfContext.TCPPort = factory.N3iwfConfig.Configuration.TCPPort
	}

	// FQDN
	if factory.N3iwfConfig.Configuration.FQDN == "" {
		contextLog.Error("FQDN is empty")
		return false
	} else {
		n3iwfContext.FQDN = factory.N3iwfConfig.Configuration.FQDN
	}

	// Private key
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.PrivateKey == "" {
			contextLog.Warn("No private key file path specified, load default key file...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.key")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.PrivateKey
		}

		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			contextLog.Errorf("Cannot read private key data from file: %+v", err)
			return false
		}
		block, _ := pem.Decode(content)
		if block == nil {
			contextLog.Error("Parse pem failed")
			return false
		}
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			contextLog.Warnf("Parse PKCS8 private key failed: %+v", err)
			contextLog.Info("Parse using PKCS1...")

			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				contextLog.Errorf("Parse PKCS1 pricate key failed: %+v", err)
				return false
			}
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			contextLog.Error("Private key is not an rsa private key")
			return false
		}

		n3iwfContext.N3IWFPrivateKey = rsaKey
	}

	// Certificate authority
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.CertificateAuthority == "" {
			contextLog.Warn("No certificate authority file path specified, load default CA certificate...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.pem")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.CertificateAuthority
		}

		// Read .pem
		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			contextLog.Errorf("Cannot read certificate authority data from file: %+v", err)
			return false
		}
		// Decode pem
		block, _ := pem.Decode(content)
		if block == nil {
			contextLog.Error("Parse pem failed")
			return false
		}
		// Parse DER-encoded x509 certificate
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			contextLog.Errorf("Parse certificate authority failed: %+v", err)
			return false
		}
		// Get sha1 hash of subject public key info
		sha1Hash := sha1.New()
		if _, err := sha1Hash.Write(cert.RawSubjectPublicKeyInfo); err != nil {
			contextLog.Errorf("Hash function writing failed: %+v", err)
			return false
		}

		n3iwfContext.CertificateAuthority = sha1Hash.Sum(nil)
	}

	// Certificate
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.Certificate == "" {
			contextLog.Warn("No certificate file path specified, load default certificate...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.pem")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.Certificate
		}

		// Read .pem
		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			contextLog.Errorf("Cannot read certificate data from file: %+v", err)
			return false
		}
		// Decode pem
		block, _ := pem.Decode(content)
		if block == nil {
			contextLog.Error("Parse pem failed")
			return false
		}

		n3iwfContext.N3IWFCertificate = block.Bytes
	}

	// UE IP address range
	if factory.N3iwfConfig.Configuration.UEIPAddressRange == "" {
		contextLog.Error("UE IP address range is empty")
		return false
	} else {
		_, ueIPRange, err := net.ParseCIDR(factory.N3iwfConfig.Configuration.UEIPAddressRange)
		if err != nil {
			contextLog.Errorf("Parse CIDR failed: %+v", err)
			return false
		}
		n3iwfContext.Subnet = ueIPRange
	}

	if factory.N3iwfConfig.Configuration.InterfaceMark == 0 {
		contextLog.Warn("IPSec interface mark is not defined, set to default value 7")
		n3iwfContext.Mark = 7
	} else {
		n3iwfContext.Mark = factory.N3iwfConfig.Configuration.InterfaceMark
	}

	return true
}

func formatSupportedTAList(info *n3iwf_context.N3IWFNFInfo) bool {
	for taListIndex := range info.SupportedTAList {

		supportedTAItem := &info.SupportedTAList[taListIndex]

		// Checking TAC
		if supportedTAItem.TAC == "" {
			contextLog.Error("TAC is mandatory.")
			return false
		}
		if len(supportedTAItem.TAC) < 6 {
			contextLog.Trace("Detect configuration TAC length < 6")
			supportedTAItem.TAC = strings.Repeat("0", 6-len(supportedTAItem.TAC)) + supportedTAItem.TAC
			contextLog.Tracef("Changed to %s", supportedTAItem.TAC)
		} else if len(supportedTAItem.TAC) > 6 {
			contextLog.Error("Detect configuration TAC length > 6")
			return false
		}

		// Checking SST and SD
		for plmnListIndex := range supportedTAItem.BroadcastPLMNList {

			broadcastPLMNItem := &supportedTAItem.BroadcastPLMNList[plmnListIndex]

			for sliceListIndex := range broadcastPLMNItem.TAISliceSupportList {

				sliceSupportItem := &broadcastPLMNItem.TAISliceSupportList[sliceListIndex]

				// SST
				if sliceSupportItem.SNSSAI.SST == "" {
					contextLog.Error("SST is mandatory.")
				}
				if len(sliceSupportItem.SNSSAI.SST) < 2 {
					contextLog.Trace("Detect configuration SST length < 2")
					sliceSupportItem.SNSSAI.SST = "0" + sliceSupportItem.SNSSAI.SST
					contextLog.Tracef("Change to %s", sliceSupportItem.SNSSAI.SST)
				} else if len(sliceSupportItem.SNSSAI.SST) > 2 {
					contextLog.Error("Detect configuration SST length > 2")
					return false
				}

				// SD
				if sliceSupportItem.SNSSAI.SD != "" {
					if len(sliceSupportItem.SNSSAI.SD) < 6 {
						contextLog.Trace("Detect configuration SD length < 6")
						sliceSupportItem.SNSSAI.SD = strings.Repeat("0", 6-len(sliceSupportItem.SNSSAI.SD)) + sliceSupportItem.SNSSAI.SD
						contextLog.Tracef("Change to %s", sliceSupportItem.SNSSAI.SD)
					} else if len(sliceSupportItem.SNSSAI.SD) > 6 {
						contextLog.Error("Detect configuration SD length > 6")
						return false
					}
				}

			}
		}

	}

	return true
}
