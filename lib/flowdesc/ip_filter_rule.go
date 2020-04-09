//go:binary-only-package

package flowdesc

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type IPFilterRule interface {
	// initial structure
	//Init() error
	// Set action of IPFilterRule
	SetAction(bool) error
	// Set Direction of IPFilterRule
	SetDirection(bool) error
	// Set Protocol of IPFilterRule
	// 0xfc stand for ip (any)
	SetProtocal(int) error
	// Set Source IP of IPFilterRule
	// format: IP or IP/mask or "any"
	SetSourceIp(string) error
	// Set Source port of IPFilterRule
	// format: {port/port-port}[,ports[,...]]
	SetSourcePorts(string) error
	// Set Destination IP of IPFilterRule
	// format: IP or IP/mask or "assigned"
	SetDestinationIp(string) error
	// Set Destination port of IPFilterRule
	// format: {port/port-port}[,ports[,...]]
	SetDestinationPorts(string) error
	// Encode the IPFilterRule
	Encode() (string, error)
	// Decode the IPFilterRule
	Decode() error
}

type ipFilterRule struct {
	action   bool   // true: permit, false: deny
	dir      bool   // false: in, true: out
	proto    int    // protocal number
	srcIp    string // <address/mask>
	srcPorts string // [ports]
	dstIp    string // <address/mask>
	dstPorts string // [ports]
}

func NewIPFilterRule() *ipFilterRule {}

func (r *ipFilterRule) SetAction(action bool) error {}

func (r *ipFilterRule) SetDirection(dir bool) error {}

func (r *ipFilterRule) SetProtocal(proto int) error {}

func (r *ipFilterRule) SetSourceIp(networkStr string) error {}

func (r *ipFilterRule) SetSourcePorts(ports string) error {}

func (r *ipFilterRule) SetDestinationIp(networkStr string) error {}

func (r *ipFilterRule) SetDestinationPorts(ports string) error {}

func (r *ipFilterRule) Encode() (string, error) {}

func (r *ipFilterRule) Decode() error {}
