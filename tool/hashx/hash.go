package hashx

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Option struct {
	Key []byte
}

// Sha256 make hash use sha256
func Sha256(str string, opts ...hashOption) (string, error) {
	return Hash(sha256.New(), str, opts...)
}

// CheckSha256 check if hash and text string are equaled.
func CheckSha256(hashStr string, text string, opts ...hashOption) bool {
	return CheckHash(sha256.New(), hashStr, text, opts...)
}

// WithKey set the key to sha256
func WithKey(str []byte) hashOption {
	return func(option *Option) {
		option.Key = str
	}
}

type hashOption func(option *Option)

// Sha1 sha1算法
func Sha1(str string, opts ...hashOption) (string, error) {
	return Hash(sha1.New(), str, opts...)
}

func CheckSha1(hashStr string, text string, opts ...hashOption) bool {
	return CheckHash(sha1.New(), hashStr, text, opts...)
}

func Md5(str string, opts ...hashOption) (string, error) {
	return Hash(md5.New(), str, opts...)
}

func CheckMd5(hashStr string, text string, opts ...hashOption) bool {
	return CheckHash(md5.New(), hashStr, text, opts...)
}

func Hash(hasher hash.Hash, str string, opts ...hashOption) (string, error) {
	var option Option
	for _, o := range opts {
		o(&option)
	}
	_, err := hasher.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(option.Key)), nil
}

func CheckHash(hasher hash.Hash, hashStr string, text string, opts ...hashOption) bool {
	h, err := Hash(hasher, text, opts...)
	if err != nil {
		return false
	}
	return h == hashStr
}
