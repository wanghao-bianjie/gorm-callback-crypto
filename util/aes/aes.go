package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	math_rand "math/rand"
	"time"
)

func CBCPKCS7Encrypt(origData, key []byte) (res []byte, err error) {
	defer func() {
		if i := recover(); i != nil {
			err = fmt.Errorf("%v", i)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func CBCPKCS7Decrypt(crypted, key []byte) (res []byte, err error) {
	defer func() {
		if i := recover(); i != nil {
			err = fmt.Errorf("%v", i)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func CBCPKCS7EncryptToHex(origData, key []byte) (string, error) {
	encrypt, err := CBCPKCS7Encrypt(origData, key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encrypt), nil
}

func CBCPKCS7DecryptFromHex(hexStr string, key []byte) (string, error) {
	crypted, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	decrypt, err := CBCPKCS7Decrypt(crypted, key)
	if err != nil {
		return "", err
	}
	return string(decrypt), nil
}

func CBCPKCS7EncryptToBase64(origData, key []byte) (string, error) {
	encrypt, err := CBCPKCS7Encrypt(origData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypt), nil
}

func CBCPKCS7DecryptFromBase64(base64Str string, key []byte) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	decrypt, err := CBCPKCS7Decrypt(crypted, key)
	if err != nil {
		return "", err
	}
	return string(decrypt), nil
}

func CBCPKCS7EncryptToHexArray(data []string, key []byte) ([]string, error) {
	if len(data) == 0 {
		return data, nil
	}
	var res []string
	for _, v := range data {
		s, err := CBCPKCS7EncryptToHex([]byte(v), key)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// CTRNoPaddingEncrypt AES/CTR/NoPadding
func CTRNoPaddingEncrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCTR(block, iv)

	dst := make([]byte, len(origData), len(origData))
	encrypter.XORKeyStream(dst, origData)
	return dst, nil
}

const (
	keyLen = 32
	ivLen  = 16
)

func GenerateAesKeyAndIV() (key, iv []byte, err error) {
	math_rand.Seed(time.Now().UnixNano())

	// Key
	key = make([]byte, keyLen)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}

	// sizeof uint64
	if ivLen < 8 {
		return nil, nil, fmt.Errorf("ivLen:%d less than 8", ivLen)
	}

	// IV:reserve 8 bytes
	iv = make([]byte, ivLen)
	if _, err := io.ReadFull(rand.Reader, iv[0:ivLen-8]); err != nil {
		return nil, nil, err
	}

	// only use 4 byte,in order not to overflow when SeekIV()
	randNumber := math_rand.Uint32()
	ivLen := len(iv)
	binary.BigEndian.PutUint64(iv[ivLen-8:], uint64(randNumber))

	return key, iv, nil
}
