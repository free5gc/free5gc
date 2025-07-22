package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type N3gaLocation struct { /* Sequence Type */
	N3gppTai       *TAI       `ber:"tagNum:0,optional"`
	N3IwfId        *N3IwFId   `ber:"tagNum:1,optional"`
	UeIpv4Addr     *IPAddress `ber:"tagNum:2,optional"`
	UeIpv6Addr     *IPAddress `ber:"tagNum:3,optional"`
	PortNumber     *int64     `ber:"tagNum:4,optional"`
	TnapId         *TNAPId    `ber:"tagNum:5,optional"`
	TwapId         *TWAPId    `ber:"tagNum:6,optional"`
	HfcNodeId      *HFCNodeId `ber:"tagNum:7,optional"`
	W5gbanLineType *LineType  `ber:"tagNum:8,optional"`
	Gli            *GLI       `ber:"tagNum:9,optional"`
	Gci            *GCI       `ber:"tagNum:10,optional"`
}
