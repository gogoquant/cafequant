package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"gopkg.in/logger.v1"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func initCryptoCipher() (cipher.Block, error) {
	key := "abcdefghijklmnopqrstuvwxyzABCDEF"
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		return nil, err
	}
	return c, nil
}

//EnCrypto 加密
func EnCrypto(plainText string) ([]byte, error) {
	block, err := initCryptoCipher()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	plaintext := []byte(plainText)

	streamCFB := cipher.NewCFBEncrypter(block, commonIV)
	cipherText := make([]byte, len(plaintext))
	streamCFB.XORKeyStream(cipherText, plaintext)
	// log.Infof("原文:%s => 密文:%x\n", plaintext, cipherText)
	return cipherText, nil
}

//DeCrypto 解密
func DeCrypto(cipherText string) ([]byte, error) {
	block, err := initCryptoCipher()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ciphertext, err := hex.DecodeString(cipherText)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	streamCFB := cipher.NewCFBDecrypter(block, commonIV)
	plaintext := make([]byte, len(ciphertext))
	streamCFB.XORKeyStream(plaintext, ciphertext)
	// log.Infof("密文:%x => 原文:%s\n", ciphertext, plaintext)
	return plaintext, nil
}
