package swagger

/*
#include "mica_gpio.h"
void call(int, int, void*);
typedef void (*closure)(int, int, void*);
#cgo LDFLAGS: -lhidapi-libusb
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	// DirectionInput INPUT
	DirectionInput = iota
	// DirectionOutput OUTPUT
	DirectionOutput
)

const (
	// StateLow LOW
	StateLow = iota
	// StateHigh HIGH
	StateHigh
)

// Callback Use for callback registration
type Callback func(int, int)

var callback Callback

//export call
func call(id C.int, state C.int, data unsafe.Pointer) {
	callback(int(id), int(state))
}

/*
SetCallback Sets the callback
*/
func SetCallback(_callback Callback) {
	callback = _callback
	var p unsafe.Pointer
	C.mica_gpio_set_callback(C.closure(C.call), p)
}

/*
GetDirection Gets the direction of id
*/
func GetDirection(id int) Direction {
	var direction = C.mica_gpio_get_direction(C.uchar(id))
	if direction == 0 {
		return INPUT
	}
	return OUTPUT
}

/*
SetDirection Sets the direction of id
*/
func setDirection(id int, direction Direction) {
	var dir uint32
	if direction == INPUT {
		dir = 0
	} else {
		dir = 1
	}
	C.mica_gpio_set_direction(C.uchar(id), dir)
}

/*
GetState Gets the state of id
*/
func GetState(id int) State {
	var s = C.mica_gpio_get_state(C.uchar(id))

	switch s {
	case 0:
		return LOW
	case 1:
		return HIGH
	default:
		return UNKNOWN
	}
}

/*
SetState Sets the state of id
*/
func setState(id int, state string) error {
	var s uint32
	switch state {
	case "LOW":
		s = 0
	case "HIGH":
		s = 1
	default:
		return errors.New("Unknown State")
	}
	C.mica_gpio_set_state(C.uchar(id), s)
	return nil
}

/*
GetEnable Gets the enable of id
*/
func getEnable(id int) bool {
	switch C.mica_gpio_get_enable(C.uchar(id)) {
	case 1:
		return true
	default:
		return false
	}
}

/*
SetEnable Sets the enable of id
*/
func SetEnable(id int, enable bool) {
	var e uint16
	if enable {
		e = 1
	} else {
		e = 0
	}
	C.mica_gpio_set_enable(C.uchar(id), C.uchar(e))
}

/*
Count get no of pins
*/
func Count() int {
	return C.MICA_GPIO_SIZE
}
