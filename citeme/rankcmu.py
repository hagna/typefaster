import sys, operator
import re

num = re.compile('\d+')

fh = open("cmudict.0.7a" or sys.argv[1])
res = {}
for i in fh:
    if i and i[0] == ';':
        continue
    phones = i.split()[1:]
    p = phones
    for j in p:
        k = re.sub(num, '', j) 
        res.setdefault(k, 0)
        res[k] += 1

sorted_res = sorted(res.iteritems(), key=operator.itemgetter(1))
for j in sorted_res:
    print j[1], j[0]
