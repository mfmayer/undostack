# Go Undo Stack

`undostack` allows operations to be done (executed) and keeps track of them in order to be able to undo and redo them.

An operation can consist of one or more actions, whose `Do()` methods are called when the opration is done/redone and whose `Undo()` methods are called in reverse order when the operation is undone.

The `Action` interface wraps the `Do()` and `Undo()` method that must be implemented for arbitrary actions:

```go 
type Action interface {
	Do() error
	Undo() error
}
```

type Action interface {
	Do() error
	Undo() error
}
