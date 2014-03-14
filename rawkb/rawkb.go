package rawkb

/*
#include <stdlib.h>
#include "test.h"
#cgo CFLAGS:-Wno-error
*/
import "C"

func SetupKeyboard() int {
    i := C.setupKeyboard()
    return int(i)
}

func RestoreKeyboard() {
    C.restoreKeyboard()
}
