package context

import (
	"crypto/rsa"
	"crypto/x509"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

var (
	NrfNfProfile     models.NfProfile
	Nrf_NfInstanceID string
	RootPrivKey      *rsa.PrivateKey
	RootCert         *x509.Certificate
	NrfPrivKey       *rsa.PrivateKey
	NrfPubKey        *rsa.PublicKey
	NrfCert          *x509.Certificate
)

func InitNrfContext() error {
	config := factory.NrfConfig
	logger.InitLog.Infof("nrfconfig Info: Version[%s] Description[%s]",
		config.Info.Version, config.Info.Description)
	configuration := config.Configuration

	Nrf_NfInstanceID = uuid.New().String()
	NrfNfProfile.NfInstanceId = Nrf_NfInstanceID
	NrfNfProfile.NfType = models.NfType_NRF
	NrfNfProfile.NfStatus = models.NfStatus_REGISTERED

	serviceNameList := configuration.ServiceNameList

	if config.GetOAuth() {
		var err error
		rootPrivKeyPath := config.GetRootPrivKeyPath()
		RootPrivKey, err = openapi.ParsePrivateKeyFromPEM(rootPrivKeyPath)
		if err != nil {
			logger.InitLog.Warnf("No root private key: %v; generate new one", err)
			err = makeDir(rootPrivKeyPath)
			if err != nil {
				return errors.Wrapf(err, "NRF init")
			}
			RootPrivKey, err = openapi.GenerateRSAKeyPair("", rootPrivKeyPath)
			if err != nil {
				return errors.Wrapf(err, "NRF init")
			}
		}

		rootCertPath := config.GetRootCertPemPath()
		RootCert, err = openapi.ParseCertFromPEM(rootCertPath)
		if err != nil {
			logger.InitLog.Warnf("No root cert: %v; generate new one", err)
			err = makeDir(rootCertPath)
			if err != nil {
				return errors.Wrapf(err, "NRF init")
			}
			RootCert, err = openapi.GenerateRootCertificate(rootCertPath, RootPrivKey)
			if err != nil {
				return errors.Wrapf(err, "NRF init")
			}
		}

		nrfPrivKeyPath := config.GetNrfPrivKeyPath()
		NrfPrivKey, err = openapi.ParsePrivateKeyFromPEM(nrfPrivKeyPath)
		if err != nil {
			logger.InitLog.Warnf("No NF priv key: %v; generate new one", err)
			NrfPrivKey, err = openapi.GenerateRSAKeyPair("", nrfPrivKeyPath)
			if err != nil {
				return errors.Wrapf(err, "NRF init")
			}
		}
		NrfPubKey = &NrfPrivKey.PublicKey

		nrfCertPath := config.GetNrfCertPemPath()
		logger.InitLog.Infof("generate new NRF cert")
		NrfCert, err = openapi.GenerateCertificate(
			string(NrfNfProfile.NfType), Nrf_NfInstanceID,
			nrfCertPath, NrfPubKey, RootCert, RootPrivKey)
		if err != nil {
			return errors.Wrapf(err, "NRF init")
		}
	}

	NFServices := InitNFService(serviceNameList, config.Info.Version)
	NrfNfProfile.NfServices = &NFServices
	return nil
}

func InitNFService(srvNameList []string, version string) []models.NfService {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	NFServices := make([]models.NfService, len(srvNameList))
	for index, nameString := range srvNameList {
		name := models.ServiceName(nameString)
		NFServices[index] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(factory.NrfConfig.GetSbiScheme()),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       factory.NrfConfig.GetSbiUri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: factory.NrfConfig.GetSbiRegisterIP(),
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(factory.NrfConfig.GetSbiPort()),
				},
			},
		}
	}
	return NFServices
}

func makeDir(filePath string) error {
	dir, _ := filepath.Split(filePath)
	if err := os.MkdirAll(dir, 0o775); err != nil {
		return errors.Wrapf(err, "makeDir(%s):", dir)
	}
	return nil
}

func SignNFCert(nfType, nfId string) error {
	nfCertPath := openapi.GetNFCertPath(factory.NrfConfig.GetCertBasePath(), nfType)

	// Get NF's Certificate from file
	nfCert, err := openapi.ParseCertFromPEM(nfCertPath)
	if err != nil {
		logger.ManagementLog.Warnf("No NF cert: %v; generate new one", err)

		// Get NF's Public key from file
		var nfPubKey *rsa.PublicKey
		nfPubKey, err = openapi.ParsePublicKeyFromPEM(nfCertPath)
		if err != nil {
			return errors.Wrapf(err, "sign NF cert")
		}

		// Generate new NF's Certificate to file
		_, err = openapi.GenerateCertificate(
			nfType, nfId, nfCertPath, nfPubKey, RootCert, RootPrivKey)
		if err != nil {
			return errors.Wrapf(err, "sign NF cert")
		}
	} else {
		nfPubkey, ok := nfCert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return errors.Errorf("No public key in NF cert")
		}

		// Re-generate new NF's Certificate to file
		_, err = openapi.GenerateCertificate(
			nfType, nfId, nfCertPath, nfPubkey, RootCert, RootPrivKey)
		if err != nil {
			return errors.Wrapf(err, "sign NF cert")
		}
	}

	return nil
}
