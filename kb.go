package typefaster

const (
	// consonants
	/* 24 of them, 12 from the set without using the thumb
	   12 more from the set that uses 1 on the thumb
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

type phone struct {
	Cmu    string // from cmupd
	Klat   string // klattese
	Ipa    string // IPA
	Mbrola string // for mbrola
	Espeak string // for espeak
	Ispeak string // for apple speak
	Des    string // deseret alphabet
}

var Phones = map[uint8]phone{
	AA: phone{Cmu: "AA"},
	AE: phone{Cmu: "AE"},
	AH: phone{Cmu: "AH"},
	AO: phone{Cmu: "AO"},
	AW: phone{Cmu: "AW"},
	AY: phone{Cmu: "AY"},
	B:  phone{Cmu: "B"},
	CH: phone{Cmu: "CH"},
	D:  phone{Cmu: "D"},
	DH: phone{Cmu: "DH"},
	EH: phone{Cmu: "EH"},
	ER: phone{Cmu: "ER"},
	EY: phone{Cmu: "EY"},
	F:  phone{Cmu: "F"},
	G:  phone{Cmu: "G"},
	HH: phone{Cmu: "HH"},
	IH: phone{Cmu: "IH"},
	IY: phone{Cmu: "IY"},
	JH: phone{Cmu: "JH"},
	K:  phone{Cmu: "K"},
	L:  phone{Cmu: "L"},
	M:  phone{Cmu: "M"},
	N:  phone{Cmu: "N"},
	NG: phone{Cmu: "NG"},
	OW: phone{Cmu: "OW"},
	OY: phone{Cmu: "OY"},
	P:  phone{Cmu: "P"},
	R:  phone{Cmu: "R"},
	S:  phone{Cmu: "S"},
	SH: phone{Cmu: "SH"},
	T:  phone{Cmu: "T"},
	TH: phone{Cmu: "TH"},
	UH: phone{Cmu: "UH"},
	UW: phone{Cmu: "UW"},
	V:  phone{Cmu: "V"},
	W:  phone{Cmu: "W"},
	Y:  phone{Cmu: "Y"},
	Z:  phone{Cmu: "Z"},
	ZH: phone{Cmu: "ZH"},
}

