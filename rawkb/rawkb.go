package rawkb

/*
#include <stdlib.h>
#include "test.h"
#cgo CFLAGS:-Wno-error
*/
import "C"
import "fmt"

func Random() int {
    fmt.Println("setupKeyboard gave us")
    fmt.Println(C.setupKeyboard())
    fmt.Println("and now restore keyboard")
    C.restoreKeyboard()
    return int(C.random())
}

func Seed(i int) {
    C.srandom(C.uint(i))
}
