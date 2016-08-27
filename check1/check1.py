#!/usr/bin/env python3
import os
print("开始")
dir1=os.getenv("STOCKDATA")+"/raw1"
dir2=os.getenv("STOCKDATA")+"/raw2"
print(dir)
for  a in os.listdir(dir1):

    i =0
    for f2 in os.listdir(dir2+"/"+a):
        f = open(dir2 + "/" + a + "/" + f2)
        #print(len(f.readlines()))
        i+=len(f.readlines())
        f.close()
    if len(os.listdir(dir1+"/"+a))!=i:
        print(a+"  "+str(i)+ " 实际"+ str(len(os.listdir(dir1+"/"+a))))
        for f3 in os.listdir(dir1+"/"+a):
          #  print(f3)
            ok = False
            for f in os.listdir(dir2+"/"+a+"/"):
                f4 = open(dir2 + "/" + a+"/"+f)
                for line in f4.readlines():
                   # print(line)
                    if line[0:10] == f3:
                        ##print(line[0:10] + " "+ f3)
                        ok=True
                        break
                f4.close()
            if ok==False:
                print(f3)
#    if os.listdir(dir1+"/"+a)==
#    for fname in os.listdir(dir+"/"+a):

