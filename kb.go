package typefaster

const (
	// consonants
	/* 24 of them, 12 from the set without using the thumb
	   12 more from the set that uses 1 
	*/
	N  = 1 << 4
	T  = 1 << 3
	R  = 1 << 1
	S  = 1 << 2
	D  = 1 << 5
	L  = 1<<1 | 1<<4
	DH = 1<<2 | 1<<3
	Z  = 1<<3 | 1<<4
	M  = 1<<1 | 1<<2
	K  = 1<<2 | 1<<3 | 1<<4
	V  = 1<<1 | 1<<3
	W  = 1<<1 | 1<<2 | 1<<3 | 1<<4
	P  = 1<<1 | 1<<2 | 1<<3
	F  = 1<<1 | 1<<5
	B  = 1<<4 | 1<<5
	HH = 1<<2 | 1<<4
	NG = 1<<2 | 1<<3 | 1<<4 | 1<<5
	SH = 1<<1 | 1<<3 | 1<<4
	G  = 1<<3 | 1<<4 | 1<<5
	Y  = 1<<1 | 1<<2 | 1<<3 | 1<<4 | 1<<5
	CH = 1<<2 | 1<<5
	JH = 1<<1 | 1<<4 | 1<<5
	TH = 1<<1 | 1<<2 | 1<<4
	ZH = 1<<1 | 1<<3 | 1<<4 | 1<<5

	// vowels
	/* 17 of them, 16 from the set that begin with 0
	   and one from the unused set that begins with 01
	*/
	AA  = 1
	IH2 = 1 | 1<<4 // maybe get rid of this one? dennis?
	AO  = 1 | 1<<2
	IH  = 1 | 1<<1
	AE  = 1 | 1<<3
	EH  = 1 | 1<<2 | 1<<3 | 1<<4
	IY  = 1 | 1<<2 | 1<<3
	EY  = 1 | 1<<5
	AH  = 1 | 1<<3 | 1<<4
	UW  = 1 | 1<<2 | 1<<3 | 1<<4 | 1<<5
	AY  = 1 | 1<<4 | 1<<5
	OW  = 1 | 1<<2 | 1<<4
	UH  = 1 | 1<<3 | 1<<4 | 1<<5
	ER  = 1 | 1<<2 | 1<<5
	AW  = 1 | 1<<2 | 1<<3 | 1<<5
	YU  = 1 | 1<<3 | 1<<5 // really Y UW
	OY  = 1 | 1<<2 | 1<<4 | 1<<5
)

type Phone struct {
	Cmu    string // from cmupd
	Klat   string // klattese
	Ipa    string // IPA
	Mbrola string // for mbrola
	Espeak string // for espeak
	Ispeak string // for apple speak
	Des    string // deseret alphabet
}

var Phones = map[uint8]Phone{
	AA: Phone{Cmu: "AA"},
	AE: Phone{Cmu: "AE"},
	AH: Phone{Cmu: "AH"},
	AO: Phone{Cmu: "AO"},
	AW: Phone{Cmu: "AW"},
	AY: Phone{Cmu: "AY"},
	B:  Phone{Cmu: "B"},
	CH: Phone{Cmu: "CH"},
	D:  Phone{Cmu: "D"},
	DH: Phone{Cmu: "DH"},
	EH: Phone{Cmu: "EH"},
	ER: Phone{Cmu: "ER"},
	EY: Phone{Cmu: "EY"},
	F:  Phone{Cmu: "F"},
	G:  Phone{Cmu: "G"},
	HH: Phone{Cmu: "HH"},
	IH: Phone{Cmu: "IH"},
	IY: Phone{Cmu: "IY"},
	JH: Phone{Cmu: "JH"},
	K:  Phone{Cmu: "K"},
	L:  Phone{Cmu: "L"},
	M:  Phone{Cmu: "M"},
	N:  Phone{Cmu: "N"},
	NG: Phone{Cmu: "NG"},
	OW: Phone{Cmu: "OW"},
	OY: Phone{Cmu: "OY"},
	P:  Phone{Cmu: "P"},
	R:  Phone{Cmu: "R"},
	S:  Phone{Cmu: "S"},
	SH: Phone{Cmu: "SH"},
	T:  Phone{Cmu: "T"},
	TH: Phone{Cmu: "TH"},
	UH: Phone{Cmu: "UH"},
	UW: Phone{Cmu: "UW"},
	V:  Phone{Cmu: "V"},
	W:  Phone{Cmu: "W"},
	Y:  Phone{Cmu: "Y"},
	Z:  Phone{Cmu: "Z"},
	ZH: Phone{Cmu: "ZH"},
}

