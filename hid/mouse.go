package hid

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type Mouse struct {
	deviceName string
	device     *os.File

	buttonState byte
}

func NewMouse(deviceName string) *Mouse {
	f, err := os.OpenFile(deviceName, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &Mouse{
		deviceName:  deviceName,
		device:      f,
		buttonState: 0,
	}
}

type Button byte

const (
	B1 Button = 1 << iota
	B2
	B3
)

func (m *Mouse) Press(b Button) {
	m.buttonState |= byte(b)
	m.Report(0, 0)
}

func (m *Mouse) Click(b Button) {
	m.Press(b)
	m.Release(b)
}

func (m *Mouse) Release(b Button) {
	m.buttonState &= ^byte(b)
	m.Report(0, 0)
}

// Zeros mouse to the upper left corner
func (m *Mouse) Zero() {
	for i := 0; i < 20; i++ {
		m.Report(-100, -100)
	}
}

func (m *Mouse) Close() {
	m.device.Close()
}

const (
	_xMax     int = (1 << 7) - 1
	_xRateMax int = 30 // pixel per milli
	_yMax     int = (1 << 7) - 1
	_yRateMax int = 30 // pixel per milli

)

func divmod(numerator, denominator int) (quotient, remainder int) {
	return numerator / denominator, numerator % denominator
}
func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

// Relative movement
func (m *Mouse) Move(x, y int) {
	xSteps, xR := divmod(abs(x), _xRateMax)
	//fmt.Printf("x %v %v \n", xSteps, xR)
	if xR > 0 {
		xSteps += 1
	}
	ySteps, yR := divmod(abs(y), _yRateMax)
	if yR > 0 {
		ySteps += 1
	}
	//fmt.Printf("y %v %v \n", ySteps, yR)

	nSteps := xSteps
	if ySteps > nSteps {
		nSteps = ySteps
	}

	path := make([]struct{ x, y int }, 1+nSteps)
	path[0] = struct{ x, y int }{x: 0, y: 0,}
	// Compute path
	for i := 1; i < nSteps; i++ {
		path[i] = struct{ x, y int }{
			x: (i * x) / nSteps,
			y: (i * y) / nSteps,
		}
	}
	// Report movements
	for i := 1; i < nSteps; i++ {
		m.Report(int8(path[i].x-path[i-1].x), int8(path[i].y-path[i-1].y))
	}

}

func (m *Mouse) Report(relX, relY int8) {
	report := make([]byte, 4)
	n := m.fill_report(report, relX, relY)

	n, err := m.device.Write(report[:n])
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Millisecond)
}

func (m *Mouse) fill_report(report []byte, x, y int8) int {
	report[0] |= m.buttonState

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, interface{}(x))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	report[1] = buf.Bytes()[0] //  X

	buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, interface{}(y))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	report[2] = buf.Bytes()[0] //  Y

	report[3] = 0 //  makes windows10 happy

	return 4 // FIXME the usbgadget report defines 8 bytesof space...
}
