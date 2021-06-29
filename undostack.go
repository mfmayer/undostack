package undostack

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Action interface that wraps Do and Undo methods that are called when an operation is done/redone or undone.
type Action interface {
	Do() error
	Undo() error
}

// Operation can consist of one or more actions, whose Do methods are called in order when the opration is done/reonde and
// whose Undo methods are called in reverse order when the operation is undone.
type Operation struct {
	Name    string
	Actions []Action
}

func (o *Operation) do() error {
	var err *multierror.Error
	for _, a := range o.Actions {
		e := a.Do()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}
	return err.ErrorOrNil()
}

func (o *Operation) undo() error {
	var err *multierror.Error
	for i := len(o.Actions) - 1; i >= 0; i-- {
		a := o.Actions[i]
		e := a.Undo()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}
	return err.ErrorOrNil()
}

// UndoStack enables operations to be done (executed) and keeps track of them in order to be able to undo and redo them.
type UndoStack struct {
	operations  []*Operation
	nextDoIndex int
	mtx         sync.Mutex
}

func (us *UndoStack) checkNextDoIndex() {
	if us.nextDoIndex > len(us.operations) {
		us.nextDoIndex = len(us.operations)
	}
	if us.nextDoIndex < 0 {
		us.nextDoIndex = 0
	}
}

// Do (execute) a new operation and put it on the stack.
//
// Previously undone operations are dropped (at least if they haven been redone before).
func (us *UndoStack) Do(op *Operation) (err error) {
	// if op != nil : clear Opartions stack > next index, add op to operations stack and execute next
	// else         : execute next
	//                   +---- nextIndex
	//   len(Operations) = 4    |
	// +---+---+---+---+ V
	// | 0 | 1 | 2 | 3 | 4
	// +---+---+---+---+
	us.mtx.Lock()
	defer us.mtx.Unlock()
	us.checkNextDoIndex()
	if op != nil {
		if us.nextDoIndex < len(us.operations) {
			// drop operations with index >= nextDoIndex
			us.operations = us.operations[:us.nextDoIndex]
		}
		us.operations = append(us.operations, op)
	}
	if us.nextDoIndex >= len(us.operations) {
		// nothing to do
		return
	}
	err = us.operations[us.nextDoIndex].do()
	us.nextDoIndex++
	return
}

// Redo a previously undone operation.
//
// Calls the Do methods of the operation's actions.
func (us *UndoStack) Redo() (err error) {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	us.checkNextDoIndex()
	if us.nextDoIndex >= len(us.operations) {
		// nothing to redo
		return
	}
	err = us.operations[us.nextDoIndex].do()
	us.nextDoIndex++
	return
}

// Undo a previously done operation.
//
// Calls the Do methods of the operation's actions in reverse order.
func (us *UndoStack) Undo() (err error) {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	us.checkNextDoIndex()
	if us.nextDoIndex <= 0 {
		// nothing to undo
		return
	}
	us.nextDoIndex--
	err = us.operations[us.nextDoIndex].undo()
	return
}

// Clear the undo stack and drop all operations
func (us *UndoStack) Clear() {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	us.operations = []*Operation{}
	us.checkNextDoIndex()
}

// State returns the number of operations that can be undone and redone
func (us *UndoStack) State() (canBeUndone int, canBeRedone int) {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	canBeUndone = us.nextDoIndex
	canBeRedone = len(us.operations) - us.nextDoIndex
	return
}
