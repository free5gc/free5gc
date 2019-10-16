//go:binary-only-package

package fsm

import (
	"errors"
	"fmt"
)

type State string
type Event string
type HandleFunc func(*FSM, Event, Args) error
type FuncTable map[State]HandleFunc

const (
	EVENT_ENTRY Event = "ENTRY EVENT"
)

type Args map[string]interface{}

type FSM struct {
	state     State
	funcTable FuncTable
}

func NewFuncTable() (table FuncTable) {}

func NewFSM(initState State, table FuncTable) (fsm *FSM, err error) {}

func (fsm *FSM) Current() State {}

func (fsm *FSM) Check(state State) bool {}

func (fsm *FSM) AddState(state State, callback HandleFunc) {}

func (fsm *FSM) SendEvent(event Event, args Args) (err error) {}

/* args is for ENTRY params*/
func (fsm *FSM) Transfer(trans State, args Args) error {}

func (fsm *FSM) AllStates() (states []State) {}

func (fsm *FSM) PrintStates() {}
