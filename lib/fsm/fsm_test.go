//go:binary-only-package

package fsm_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/fsm"
	"testing"
)

const (
	ACITVE    fsm.State = "ACTIVE"
	INACITVE  fsm.State = "INACITVE"
	EXCEPTION fsm.State = "EXCEPTION"
)

const (
	MESSAGE fsm.Event = "MESSAGE"
)

func Activefunc(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {}
func Inactivefunc(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {}
func Exceptionfunc(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {}

func TestInitFSM(t *testing.T) {}
