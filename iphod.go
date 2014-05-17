package typefaster

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)
type iphodrecord struct {
	Nphones  int
	Phonemes string
}

var IPHOD map[string]iphodrecord

func Readiphod(iphod string) error {
	IPHOD = make(map[string]iphodrecord)
	fh, err := os.Open(iphod)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		l := strings.Fields(scanner.Text())
		v := l[1]
		nphones, err := strconv.Atoi(l[5])
		if err != nil {
			fmt.Println(err)
			nphones = 0
		}
		phonemes := strings.ToLower(l[2])
		IPHOD[v] = iphodrecord{nphones, phonemes}
	}
	return nil
}



