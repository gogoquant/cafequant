# -*- encoding:utf-8 -*-

'''
    @brief used for bitmap support
    @author mhw
    @data 2016-12-10

'''

class BitMap(object):

    def __init__(self,max):

        '确定所需数组个数'
        self.osize = 31
        self.size = int ((max + self.osize - 1) / self.osize)
        self.array = [0 for i in range(self.size)]

    def bitindex(self,num):
        '确定数组中元素的位索引'
        return num % self.osize

    def set(self,num):
        '将元素所在的位置1'
        elemindex = num / self.osize
        byteindex = self.bitindex(num)
        ele = self.array[elemindex]
        self.array[elemindex] = ele | (1 << byteindex)

    def clean(self,num):
        '将元素所在的位置1'
        elemindex = num / self.osize
        byteindex = self.bitindex(num)
        ele = self.array[elemindex]
        self.array[elemindex] = ele & (~(1 << byteindex))

    def check(self,i):
        '检测元素存在的位置'
        elemindex = i / self.osize
        byteindex = self.bitindex(i)
        if self.array[elemindex] & (1 << byteindex):
            return True
        return False

    def getsize(self):
        return self.size

if __name__ == '__main__':
    Max = ord('z')
    suffle_array = [x for x in 'qwelmfg']
    result = []
    bitmap = Bitmap(Max)
    for c in suffle_array:
        bitmap.set(ord(c))
    for i in range(Max+1):
        if bitmap.check(i):
            result.append(chr(i))
    print u'原始数组为:    %s' % suffle_array
    print u'排序后的数组为: %s' % result
