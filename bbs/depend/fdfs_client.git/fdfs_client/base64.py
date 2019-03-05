# -*- coding: utf-8 -*-
import logging
from struct import *

class Base64(object):
    BASE64_IGNORE = -1
    
    BASE64_PAD       = -2
    #char array
    valueToChar       = []
    #int array
    charToValue       = []
    #char value
    charToPad          = '.'
    
    lineLength          = 72
    
    def __init__(self):
        self.init('-', '_', '.')
    
    def init(self, chPlus, chSplash, chPad):
        #0..25 -> 'A'..'Z'
        for i in range(26):
            self.valueToChar.append(chr(ord('A') + i))
        #26..51 -> 'a'..'z'
        for i in range(26):
            self.valueToChar.append(chr(ord('a') + i))
        #52..61 -> '0'..'9'
        for i in range(10):
            self.valueToChar.append(chr(ord('0') + i))
            
        self.valueToChar.append(chPlus)
        self.valueToChar.append(chSplash)
        
        #build translate defaultCharToValue table only once.
        for i in range(256):
            self.charToValue.append(self.BASE64_IGNORE)
            
        for i in range(64):
            self.charToValue[ord(self.valueToChar[i])] = i
            
        self.charToValue[ord(chPad)] = self.BASE64_PAD
        """for i in range(4):
            self.charToPad.append(self.BASE64_PAD)"""
        self.charToPad = chPad
    
    def base64_decode_auto(self, file_name):
        nRemain = len(file_name) % 4
        if nRemain == 0:
            return self.base64_decode(file_name)
        else:
            char_list = [self.charToPad for _ in range(4 - nRemain)]
            file_name = pack("%ds%dc" % (len(file_name), len(char_list)), file_name, *char_list)
            return self.base64_decode(file_name)
        
    def buff2int(self, name_array):
        return ord(name_array[0]) << 24 | ord(name_array[1]) << 16 | ord(name_array[2]) << 8 | ord(name_array[3])
    
    def buff2long(self, name_array):
        #python 移位计算的bug，用于小文件，只取后4位
        """return (ord(name_array[0]) << 56) | (ord(name_array[1]) << 48) \
               | (ord(name_array[2]) << 40) | (ord(name_array[3]) << 32) \
               | (ord(name_array[4]) << 24) | (ord(name_array[5]) << 16) \
               | (ord(name_array[6]) << 8) | ord(name_array[7])"""
        return (ord(name_array[4]) << 24) | (ord(name_array[5]) << 16) \
               | (ord(name_array[6]) << 8) | ord(name_array[7])
        
    def base64_decode(self, file_name):
        file_name_len = len(file_name)        
        cycle = 0
        combined = 0
        dummies = 0
        pDest = []
        for i in range(file_name_len):
            char_i = file_name[i]
            value_i = self.BASE64_IGNORE
            if ord(char_i) <= 255:
                value_i = self.charToValue[ord(char_i)]
            if value_i == self.BASE64_IGNORE:
                continue            
            else:
                if value_i == self.BASE64_PAD:
                    value_i = 0
                    dummies += 1
                if cycle == 0:
                    combined = value_i
                    cycle = 1                    
                elif cycle == 1:
                    combined <<= 6
                    combined |= value_i
                    cycle = 2                    
                elif cycle == 2:
                    combined <<= 6
                    combined |= value_i
                    cycle = 3                    
                elif cycle == 3:
                    combined <<= 6
                    combined |= value_i
                    
                    #we have just completed a cycle of 4 chars.
                    #the four 6-bit values are in combined in big-endian order
                    #peel them off 8 bits at a time working lsb to msb
                    #to get our original 3 8-bit bytes back
                    
                    pDest.append(chr(combined >> 16))
                    pDest.append(chr((combined & 0x0000FF00) >> 8))
                    pDest.append(chr(combined & 0x000000FF))                    
                    cycle = 0
        
        if cycle <> 0:
            logging.error("Input to decode not an even multiple of 4 characters; pad with %c" % self.charToPad)
            return []
        
        dest_len = len(pDest) - dummies
        pDest = pDest[0:dest_len]
        return pDest
            
            