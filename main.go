package main

import (
	"github.com/edgera/pig/hid"
)

func main() {
	m := hid.NewMouse("/dev/hidg0")
	defer m.Close()
	
	m.Zero()
	m.Click(hid.B1)
	m.Move(0, 540)
	m.Move(960, 0)
	m.Move(-960, -540)
}
