package main

import (
	. ".."
)

func main() {
	InitConfig()
	a := AptMethod{}
	a.SendCapabilities()
	a.Run()
}
