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

var encodemap map[string]uint8
var decodemap map[uint8]string

func init() {
	s := "AA.AE.AH.AO.AW.AY.B.CH.D.DH.EH.ER.EY.F.G.HH.IH.IY.JH.K.L.M.N.NG.OW.OY.P.R.S.SH.T.TH.UH.UW.V.W.Y.Z.ZH"
	encodemap = make(map[string]uint8)
	decodemap = make(map[uint8]string)
	for i, v := range strings.Split(s, ".") {
		j := 'A' + uint8(i)
		encodemap[v] = j
		decodemap[j] = v
	}
}

/* maybe this is silly but the cpt uses strings and each char is an item, so
we use this to encode cmu style phonemes (AH AA etc.) to single chars for use
in the cpt, and then decode when we want to display them.  Alternatively we
could change cpt to use []string, but it's already working with string so why
bother.
*/
func encode(p string) string {
	res := ""
	for _, v := range strings.Split(p, ".") {
		res += string(byte(encodemap[v]))
	}
	return res

}

func decode(p string) string {
	res := []string{}
	for _, v := range p {
		res = append(res, decodemap[uint8(v)])
	}
	return strings.Join(res, ".")
}

// make a compact prefix tree out of iphod for the keyboard to use
// spelling words out of phonemes
func Maketree(iphod string) (*MemTree, error) {
	tree := NewMemTree("root")
	cb := func(word, phonemes string, nphones int) {
		phonemes = encode(phonemes)
		tree.Insert(phonemes, word)
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
	scanner.Scan() // chomp first line
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
