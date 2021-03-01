package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
)

func PKCS5UnPadding(origData []byte) ([]byte, bool) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if unPadding > length {
		return nil, false
	}
	return origData[:length-unPadding], true
}

func PKCS5Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}

// DESEncrypt 提供 DES 加密的一种简便用法
// 使用 CBC 模式， PKCS5 填充
func DESEncrypt(src, key, iv []byte) (res []byte, err error) {
	if len(key) >= 8 {
		key = key[:8]
	} else {
		return nil, fmt.Errorf("invalid key size %d", len(key))
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) >= block.BlockSize() {
		iv = iv[:block.BlockSize()]
	} else {
		return nil, fmt.Errorf("invalid iv size %d", len(iv))
	}
	paddedData := PKCS5Padding(src, block.BlockSize())
	bm := cipher.NewCBCEncrypter(block, iv)
	res = make([]byte, len(paddedData))
	bm.CryptBlocks(res, paddedData)
	return
}

// DESDecrypt 提供一种 DES 解密的简便用法
// 使用 CBC 模式， PKCS5 填充
func DESDecrypt(src, key, iv []byte) (res []byte, err error) {
	if len(key) >= 8 {
		key = key[:8]
	} else {
		return nil, fmt.Errorf("invalid key size %d", len(key))
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) >= block.BlockSize() {
		iv = iv[:block.BlockSize()]
	} else {
		return nil, errors.New("IV insufficient length")
	}
	bm := cipher.NewCBCDecrypter(block, iv)
	res = make([]byte, len(src))
	bm.CryptBlocks(res, src)
	res, ok := PKCS5UnPadding(res)
	if !ok {
		return nil, errors.New("decryptionFailed")
	}
	return
}
