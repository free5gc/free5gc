package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SessionAMBR struct { /* Sequence Type */
	AmbrUL Bitrate `ber:"tagNum:1"`
	AmbrDL Bitrate `ber:"tagNum:2"`
}
