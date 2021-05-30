# Go Undo Stack

`undostack` allows operations to be done (executed) and keeps track of them in order to be able to undo and redo them.

[![GoDoc](https://godoc.org/github.com/mfmayer/undostack?status.svg)](http://godoc.org/github.com/mfmayer/undostack)

An `Operation` can consist of one or more `Action`s, whose `Do()` methods are called when the opration is done/redone and whose `Undo()` methods are called in reverse order when the operation is undone.

`Action` is the interface that wraps the `Do()` and `Undo()` method that must be implemented for arbitrary actions:

```go
type Action interface {
  Do() error
  Undo() error
}
```

## Example

[Not quite serious example with an operation that holds two actions](examples/main.go):

```go
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
```
