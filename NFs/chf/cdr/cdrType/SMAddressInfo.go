package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SMAddressInfo struct { /* Sequence Type */
	SMAddressType   *SMAddressType     `ber:"tagNum:0,optional"`
	SMAddressData   *asn.GraphicString `ber:"tagNum:1,optional"`
	SMAddressDomain *SMAddressDomain   `ber:"tagNum:2,optional"`
}
