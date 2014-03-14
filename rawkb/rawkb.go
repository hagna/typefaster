package rawkb

/*
#include <stdlib.h>
#include "test.h"
#cgo CFLAGS:-Wno-error
*/
import "C"
import "unsafe"

func SetupKeyboard() int {
	i := C.setupKeyboard()
	return int(i)
}

func RestoreKeyboard() {
	C.restoreKeyboard()
}

func ReadOnce() uint8 {
	var ch C.char = 23
	i := C.read(0, (*C.char)ch, 1)
	return uint8(i)	
}
