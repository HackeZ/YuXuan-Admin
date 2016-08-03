package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// PassEncode Encode Password.
func PassEncode(pwd, salt string) string {
	return passEncodepbkdf2([]byte(pwd), []byte(salt))
}

// passEncode Encoding password by pbkdf2.
// For details, please refer to Python and Django documents.
// django.contrib.auth.hashers.make_password
// django.utils.crypto import pbkdf2
// hashlib.sha256
// base64
func passEncodepbkdf2(pwd, salt []byte) string {
	// Step 3 : init encode
	iteration := 15000   // 加密算法的迭代次数，15000 次
	digest := sha256.New // digest 算法，使用 sha256

	// Step 2 : encode by pbkdf2.
	dk := pbkdf2.Key(pwd, salt, iteration, 32, digest)

	// Step 3 : encode Base64.
	str := base64.StdEncoding.EncodeToString(dk)

	return str
}

// GetNewSalt return a new random salt.
func GetNewSalt(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < int(length); i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
