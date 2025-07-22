package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type RecipientInfo struct { /* Sequence Type */
	RecipientIMSI            *IMSI            `ber:"tagNum:0,optional"`
	RecipientMSISDN          *MSISDN          `ber:"tagNum:1,optional"`
	RecipientOtherAddress    *SMAddressInfo   `ber:"tagNum:2,optional"`
	RecipientSCCPAddress     *AddressString   `ber:"tagNum:3,optional"`
	RecipientReceivedAddress *SMAddressInfo   `ber:"tagNum:4,optional"`
	SMDestinationInterface   *SMInterface     `ber:"tagNum:5,optional"`
	SMRecipientProtocolID    *asn.OctetString `ber:"tagNum:6,optional"`
	/* Sequence of = 35, FULL Name = struct RecipientInfo__recipientOtherAddresses */
	/* SMAddressInfo */
	RecipientOtherAddresses []SMAddressInfo `ber:"tagNum:7,optional"`
}
