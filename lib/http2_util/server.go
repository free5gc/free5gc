//go:binary-only-package

package http2_util

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func NewServer(bindAddr string, tlskeylog string, handler http.Handler) (server *http.Server, err error) {}
