package adapter

import (
	"errors"
	"net/url"
	"strings"
)

// errMissingConfiguration is returned by Apply when an NF block has no
// `configuration` section.
var errMissingConfiguration = errors.New("configuration block is missing")

// validLogLevels are the logger levels every free5GC NF accepts.
var validLogLevels = map[string]bool{
	"trace": true, "debug": true, "info": true, "warn": true,
	"error": true, "fatal": true, "panic": true,
}

func isHTTPScheme(s string) bool { return s == "http" || s == "https" }

func isAbsoluteURL(s string) bool {
	u, err := url.Parse(s)
	return s != "" && err == nil && u.Scheme != "" && u.Host != ""
}

func isMongoURL(s string) bool {
	return strings.HasPrefix(s, "mongodb://") || strings.HasPrefix(s, "mongodb+srv://")
}
