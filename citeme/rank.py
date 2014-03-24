import sys, operator

fh = open("iphod.txt" or sys.argv[1])
res = {}
for i in fh:
    phones = i.split()[2]
    p = phones.split('.')
    for j in p:
        res.setdefault(j, 0)
        res[j] += 1

sorted_res = sorted(res.iteritems(), key=operator.itemgetter(1))
for j in sorted_res:
    print j[1], j[0]
