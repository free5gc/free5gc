package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type OriginatorInfo struct { /* Sequence Type */
	OriginatorIMSI            *IMSI            `ber:"tagNum:0,optional"`
	OriginatorMSISDN          *MSISDN          `ber:"tagNum:1,optional"`
	OriginatorOtherAddress    *SMAddressInfo   `ber:"tagNum:2,optional"`
	OriginatorSCCPAddress     *AddressString   `ber:"tagNum:3,optional"`
	OriginatorReceivedAddress *SMAddressInfo   `ber:"tagNum:4,optional"`
	SMOriginatorInterface     *SMInterface     `ber:"tagNum:5,optional"`
	SMOriginatorProtocolID    *asn.OctetString `ber:"tagNum:6,optional"`
	/* Sequence of = 35, FULL Name = struct OriginatorInfo__originatorOtherAddresses */
	/* SMAddressInfo */
	OriginatorOtherAddresses []SMAddressInfo `ber:"tagNum:7,optional"`
}
