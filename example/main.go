package main

import (
	"fmt"

	"github.com/mfmayer/undostack"
)

// Implement welcome action
type welcomeAction struct{}

func (wa *welcomeAction) Do() error {
	fmt.Println("Welcome my friend.")
	return nil
}

func (wa *welcomeAction) Undo() error {
	fmt.Println("Go away, I don't like you.")
	return nil
}

// Implement action to offer a seat
type offerSeatAction struct{}

func (sfa *offerSeatAction) Do() error {
	fmt.Println("Please have a seat.")
	return nil
}

func (sfa *offerSeatAction) Undo() error {
	fmt.Println("Please stand up.")
	return nil
}

func main() {
	// initalize the undo stack
	undoStack := undostack.UndoStack{}

	// create the receive operation and add the actions welcome and offer a seat
	receiveOperation := undostack.Operation{
		Name: "Receive Guest",
		Actions: []undostack.Action{
			&welcomeAction{},
			&offerSeatAction{},
		},
	}

	// Do, Undo and Redo the opearion
	fmt.Println("# Do receive:")
	undoStack.Do(&receiveOperation)

	fmt.Println("\n# Undo receive:")
	undoStack.Undo()

	fmt.Println("\n# Redo receive:")
	undoStack.Redo()
}
