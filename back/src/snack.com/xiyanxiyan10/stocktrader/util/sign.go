package util

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const (
	// KSIGNTYPERSA2 KSIGNTYPERSA2
	KSIGNTYPERSA2 = "RSA2"
	// KSIGNTYPERSA KSIGNTYPERSA
	KSIGNTYPERSA = "RSA"
	// TimeFormat TimeFormat
	TimeFormat = "2006-01-02 15:04:05"
)

// SinData SinData
func SinData(param url.Values, privateKey string) (s string, err error) {
	pri, err := encoding.ParsePKCS1PrivateKey(encoding.FormatPrivateKey(privateKey))
	if err != nil {
		fmt.Errorf("解析私有错误:%s", err.Error())
		return "", err
	}
	return signWithPKCS1v15(param, pri, crypto.SHA256)

}

func signWithPKCS1v15(param url.Values, privateKey *rsa.PrivateKey, hash crypto.Hash) (s string, err error) {
	if param == nil {
		param = make(url.Values)
	}

	var pList = make([]string, 0)
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")
	fmt.Printf("排序后字段:%s\n", src)
	sig, err := encoding.SignPKCS1v15WithKey([]byte(src), privateKey, hash)
	if err != nil {
		fmt.Errorf("签名失败,err:%s", err.Error())
		return "", errors.New("签名失败")
	}
	s = base64.URLEncoding.EncodeToString(sig)
	fmt.Printf("签名sign=%s\n", s)
	return s, nil
}

// VerifySign VerifySign
func VerifySign(data url.Values, publicKey string) (ok bool, err error) {
	pub, err := encoding.ParsePKCS1PublicKey(encoding.FormatPublicKey(publicKey))
	if err != nil {
		fmt.Errorf("解析公钥错误:%s", err.Error())
		return false, err
	}
	return verifySign(data, pub)
}

func verifySign(data url.Values, key *rsa.PublicKey) (ok bool, err error) {
	sign := data.Get("sign")
	signType := data.Get("sign_type")
	fmt.Printf("验证签名sign=%s,sign_type=%s\n", sign, signType)
	var keys = make([]string, 0)
	for key, value := range data {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	var pList = make([]string, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(data.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")
	fmt.Printf("验证签名排序后字段:%s\n", s)
	return verifyData([]byte(s), signType, sign, key)
}

func verifyData(data []byte, signType, sign string, key *rsa.PublicKey) (ok bool, err error) {
	signBytes, err := base64.URLEncoding.DecodeString(sign)
	if err != nil {
		fmt.Errorf("base64 decode error: %s", err.Error())
		return false, err
	}

	if signType == KSIGNTYPERSA {
		err = encoding.VerifyPKCS1v15WithKey(data, signBytes, key, crypto.SHA1)
	} else {
		err = encoding.VerifyPKCS1v15WithKey(data, signBytes, key, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
