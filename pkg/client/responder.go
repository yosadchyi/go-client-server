package client

import "fmt"

// Responder is responsible for providing feedback on command execution
type Responder interface {
	// Error reports an error
	Error(err error)
	// Ok reports success
	Ok()
	// Bye reports shutdown
	Bye()
	// Help shows help
	Help()
}

type interactiveResponder struct {
}

// NewInteractiveResponder returns new interactive responder, used when input taken from stdin
func NewInteractiveResponder() Responder {
	return &interactiveResponder{}
}

func (r *interactiveResponder) Error(err error) {
	println(err.Error())
}

func (r *interactiveResponder) Ok() {
	println("OK")
}

func (r *interactiveResponder) Bye() {
	println("Bye")
}

func (r *interactiveResponder) Help() {
	fmt.Println(`Commands:
	+ITEM
		add item with data 'ITEM'
	-INDEX
		remove item with index INDEX, where index is an integer number
	<INDEX
		get item with index INDEX, where index is an integer number
	*
		list all items
	^C
		quit`)
}

type batchResponder struct {
}

// NewBatchResponder returns new batch responder, used when input taken from file
func NewBatchResponder() Responder {
	return &batchResponder{}
}

func (r *batchResponder) Error(err error) {
	println(err.Error())
}

func (r *batchResponder) Ok() {
}

func (r *batchResponder) Bye() {
}

func (r *batchResponder) Help() {
}
