# pkb sample.txt
37 phonemes
12 unknown words of total length 40

# pkb -decode phones.txt
the
butcher
likes
the
meat|meet

# pkb -v sample.txt 
AH.B.AE.N.D.AH.N.D
JH.AA.G.ER
Unknown:Rhonda
JH.OY.N

# pkb -genfromiphod iphod.txt -treename iphod
putting tree in "iphod"
100 nodes
50 interior nodes
50 leaves

# pkb -genfromcmupd cmudist -treename qux
putting tree in "qux"
100 nodes
50 interior nodes
50 leaves


# pkb -v -t iphod -decode sample.txt
meet
meat

# pkb -t ./iphod -print
aardvark
aaron
ab
aback
...



