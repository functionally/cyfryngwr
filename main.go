package main

import (
	"github.com/functionally/cyfryngwr/cwtch"
)

func main() {
	cwtchbot := cwtch.Connect(".cyfryngwr/", "cyfryngwr", "Cyfryngwr, a cwtch agent")
	cwtch.Loop(cwtchbot)
}
