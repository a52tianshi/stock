#!/usr/bin/env python3
import os
print("开始")
dir=os.getenv("STOCKDATA")+"/raw1"
print(dir)
try :os.remove("8月振幅")
except IOError: pass
outf = open("8月振幅","w")
num = 0
for  a in os.listdir(dir):
    max = float(0.0)
    min =float(1000.0)
    num+=1
    print(a,num)
    sumamp = float(0.0)
    for fname in os.listdir(dir+"/"+a):

        f=open(dir+"/"+a+"/"+fname,encoding="GBK")
        i=0
        ##print("1",f)
        daymax = float(0.0)
        daymin = float(1000.0)
        for line in f:
            ##print("2")
            if i!=0 and len(line)>10:
                price=line.split("\t",2)
                ##print(price[1])
                price2=float(price[1])
                if price2>max:
                    max=price2
                if price2>daymax:
                    daymax=price2
                if price2<daymin:
                    daymin=price2
                if price2<min:
                    min=price2
            ##print("3")
            i+=1
        sumamp += (daymax - daymin)
        f.close()
    outf.write(a+","+str(max)+","+str(min)+","+str(sumamp)+"\n")
outf.close()
