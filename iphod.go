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

// puts the iphod into map[string]iphodrecord
func Readiphod(iphod string) error {
	IPHOD = make(map[string]iphodrecord)
	cb := func(v, phonemes string, nphones int) {
		IPHOD[v] = iphodrecord{nphones, phonemes}
	}
	return readiphod(iphod, cb)
}

// make a compact prefix tree out of iphod for the keyboard to use
// spelling words out of phonemes
func Maketree(iphod string) (*Tree, error) {
	tree := NewTree("root")
	cb := func(word, phonemes string, nphones int) {
		tree.Insert(tree.Root, word, phonemes)
	}
	err := readiphod(iphod, cb)
	if err != nil {
		return tree, err
	} 
	return tree, nil
}

// calls cb with each word phoneme nphones fields in the iphod
func readiphod(iphod string, cb func(word, phonemes string, nphones int)) error {
	fh, err := os.Open(iphod)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		l := strings.Fields(scanner.Text())
		v := strings.ToLower(l[1])
		nphones, err := strconv.Atoi(l[5])
		if err != nil {
			fmt.Println(err)
			nphones = 0
		}
		phonemes := l[2]
		cb(v, phonemes, nphones)
	}
	return nil
}



