package aes

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	aesKeyStr := "1234567890123456"
	data, err := CBCPKCS7EncryptToHex([]byte("qwe"), []byte(aesKeyStr))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data)
	t.Log(len(data))
	bytes, _ := hex.DecodeString(data)
	toString := base64.StdEncoding.EncodeToString(bytes)
	t.Log(toString)
	t.Log(len(toString))
}

func TestAesDecrypt(t *testing.T) {
	aesKeyStr := "1234567890123456"
	data := "3ebc237cbd9fe46c39a0908ccc5dee81"
	data, err := CBCPKCS7DecryptFromHex(data, []byte(aesKeyStr))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data)
}

func TestCTRNoPaddingEncrypt(t *testing.T) {
	key, iv, err := GenerateAesKeyAndIV()
	if err != nil {
		t.Fatal(err)
	}
	data := "123456"
	encrypt, err := CTRNoPaddingEncrypt([]byte(data), key, iv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(encrypt))
	decrypt, err := CTRNoPaddingEncrypt(encrypt, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decrypt))
}

func TestCBCPKCS7EncryptToBase64(t *testing.T) {
	aesKeyStr := "1234567890123456"
	str := "1234512345"
	//str := "qwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfgqwertasdfg@qq.com"
	toBase64, err := CBCPKCS7EncryptToBase64([]byte(str), []byte(aesKeyStr))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(toBase64)
	t.Log(len(toBase64))
	toHex, err := CBCPKCS7EncryptToHex([]byte(str), []byte(aesKeyStr))
	t.Log(toHex)
	t.Log(len(toHex))
	fromBase64, err := CBCPKCS7DecryptFromBase64(toBase64, []byte(aesKeyStr))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fromBase64)
}

func TestLength(t *testing.T) {
	aesKeyStr := "1234567890123456"
	var str string
	for i := 0; i < 1024; i++ {
		str += "@"
		toBase64, _ := CBCPKCS7EncryptToBase64([]byte(str), []byte(aesKeyStr))
		toHex, _ := CBCPKCS7EncryptToHex([]byte(str), []byte(aesKeyStr))
		t.Log(len(str), len(toBase64), len(toHex))
	}
}
