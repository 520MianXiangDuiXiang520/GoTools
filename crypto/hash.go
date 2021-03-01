package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

// 使用 SHA256 哈希
func SHA256(sList []string) (res string) {
	s := strings.Join(sList, "")
	hash := sha256.New()
	hash.Write([]byte(s))
	hash.Sum(nil)
	bytes := hash.Sum(nil)
	res = hex.EncodeToString(bytes)
	return
}

func SHA1(sList []string) (res string) {
	s := strings.Join(sList, "")
	hash := sha1.New()
	hash.Write([]byte(s))
	hash.Sum(nil)
	bytes := hash.Sum(nil)
	res = hex.EncodeToString(bytes)
	return
}

func SHA512(sList []string) (res string) {
	s := strings.Join(sList, "")
	hash := sha512.New()
	hash.Write([]byte(s))
	hash.Sum(nil)
	bytes := hash.Sum(nil)
	res = hex.EncodeToString(bytes)
	return
}

// 使用 MD5 做哈希
func MD5(strList []string) (h string) {
	r := strings.Join(strList, "")
	hash := md5.New()
	hash.Write([]byte(r))
	return hex.EncodeToString(hash.Sum(nil))
}
