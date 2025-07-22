package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

const (
	DIAMETER_LOGOUT               TerminationCause = 1
	DIAMETER_SERVICE_NOT_PROVIDED TerminationCause = 2
	DIAMETER_BAD_ANSWER           TerminationCause = 3
	DIAMETER_ADMINISTRATIVE       TerminationCause = 4
	DIAMETER_LINK_BROKEN          TerminationCause = 5
	DIAMETER_AUTH_EXPIRED         TerminationCause = 6
	DIAMETER_USER_MOVED           TerminationCause = 7
	DIAMETER_SESSION_TIMEOUT      TerminationCause = 8
)

type TerminationCause diam_datatype.Enumerated
