package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// CertificateMode 定義證書管理模式
type CertificateMode string

const (
	// ModeExternal 僅使用外部證書，不存在則報錯
	ModeExternal CertificateMode = "external"
	// ModeAuto 優先使用外部證書，不存在則自動生成
	ModeAuto CertificateMode = "auto"
	// ModeGenerate 總是生成新證書，忽略外部證書
	ModeGenerate CertificateMode = "generate"
)

// Manager 管理證書的載入和生成
type Manager struct {
	certPath string
	keyPath  string
	mode     CertificateMode
}

// NewManager 創建新的證書管理器
func NewManager(certPath, keyPath string) *Manager {
	return &Manager{
		certPath: certPath,
		keyPath:  keyPath,
		mode:     ModeAuto, // 默認模式
	}
}

// CertificateExists 檢查外部證書是否存在且可讀
func (m *Manager) CertificateExists() (bool, error) {
	fileInfo, err := os.Stat(m.certPath)
	if os.IsNotExist(err) {
		// 文件不存在
		return false, nil
	}
	if err != nil {
		// 其他錯誤（例如權限問題）
		return false, fmt.Errorf("failed to check certificate: %w", err)
	}

	// 文件存在，檢查是否為普通文件
	if fileInfo.IsDir() {
		return false, fmt.Errorf("certificate path is a directory: %s", m.certPath)
	}

	// 嘗試打開文件以檢查讀取權限
	file, err := os.Open(m.certPath)
	if err != nil {
		// 無法打開文件（可能是權限問題）
		return false, fmt.Errorf("cannot open certificate file: %w", err)
	}
	file.Close()

	return true, nil
}

// ExtractNfInstanceID 從證書中提取 NF Instance ID
// 根據 3GPP TS 33.310，URI SAN 格式為：urn:uuid:<uuid>
func (m *Manager) ExtractNfInstanceID(cert *x509.Certificate) (string, error) {
	if len(cert.URIs) == 0 {
		return "", fmt.Errorf("no URI SAN in certificate")
	}

	// 使用第一個 URI
	uri := cert.URIs[0]

	// 檢查 URI scheme
	if uri.Scheme != "urn" {
		return "", fmt.Errorf("invalid URI scheme: expected 'urn', got '%s'", uri.Scheme)
	}

	// 提取 UUID from urn:uuid:<uuid>
	// uri.Opaque 包含 "uuid:<uuid>" 部分
	if len(uri.Opaque) < 5 || uri.Opaque[:5] != "uuid:" {
		return "", fmt.Errorf("invalid UUID format in URI: expected 'uuid:' prefix, got '%s'", uri.Opaque)
	}

	uuid := uri.Opaque[5:] // 去掉 "uuid:" 前綴
	if uuid == "" {
		return "", fmt.Errorf("empty UUID in URI SAN")
	}

	return uuid, nil
}

// LoadOrGenerateCertificate 根據模式載入或生成證書
// generateFn 是生成新證書的函數
// 返回：證書、提取的 UUID、錯誤
func (m *Manager) LoadOrGenerateCertificate(
	generateFn func() (*x509.Certificate, error),
) (*x509.Certificate, string, error) {
	exists, err := m.CertificateExists()
	if err != nil {
		return nil, "", err
	}

	// 根據模式處理
	switch m.mode {
	case ModeExternal:
		// 僅使用外部證書模式
		if !exists {
			return nil, "", fmt.Errorf("external certificate not found at %s (mode: external)", m.certPath)
		}
		return m.loadExternalCertificate()

	case ModeAuto:
		// 自動模式：優先外部證書
		if exists {
			cert, uuid, err := m.loadExternalCertificate()
			if err != nil {
				// 外部證書載入失敗，記錄警告並回退到生成
				return m.generateCertificate(generateFn)
			}
			return cert, uuid, nil
		}
		// 外部證書不存在，生成新證書
		return m.generateCertificate(generateFn)

	case ModeGenerate:
		// 生成模式：總是生成新證書
		return m.generateCertificate(generateFn)

	default:
		return nil, "", fmt.Errorf("invalid certificate mode: %s", m.mode)
	}
}

// loadExternalCertificate 載入外部證書並提取 UUID
func (m *Manager) loadExternalCertificate() (*x509.Certificate, string, error) {
	// 讀取證書文件
	certPEM, err := os.ReadFile(m.certPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read certificate file: %w", err)
	}

	// 解析證書（此處簡化，實際需要解析 PEM 格式）
	// 暫時先用簡單實現
	cert, err := parseCertFromPEM(certPEM)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	// 提取 UUID
	uuid, err := m.ExtractNfInstanceID(cert)
	if err != nil {
		return nil, "", fmt.Errorf("failed to extract NF Instance ID from certificate: %w", err)
	}

	return cert, uuid, nil
}

// generateCertificate 使用提供的函數生成新證書
func (m *Manager) generateCertificate(
	generateFn func() (*x509.Certificate, error),
) (*x509.Certificate, string, error) {
	if generateFn == nil {
		return nil, "", fmt.Errorf("certificate generation function is nil")
	}

	cert, err := generateFn()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate certificate: %w", err)
	}

	// 從生成的證書中提取 UUID
	uuid, err := m.ExtractNfInstanceID(cert)
	if err != nil {
		return nil, "", fmt.Errorf("failed to extract NF Instance ID from generated certificate: %w", err)
	}

	return cert, uuid, nil
}

// SetMode 設置證書管理模式
func (m *Manager) SetMode(mode CertificateMode) {
	m.mode = mode
}

// parseCertFromPEM 解析 PEM 格式的證書
func parseCertFromPEM(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("PEM block type is not CERTIFICATE, got: %s", block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}
