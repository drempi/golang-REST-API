package cryptpack

// This package does not import anything from project

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

// Password is the thing you encrypt/decrypt with in cookies
var Password string

// Password2 is the thing with which you encrypt passwords
var Password2 string

// Key its Password but hashed
var Key []byte

// Gcm I dont know
var Gcm cipher.AEAD

// InitializeCrypt initializes all the necessary things
func InitializeCrypt() {
	Password = "J^T*H8y*YBGGg99ikjo)IJIIHYGbjkfsjfb9u(*##"
	Password2 = "9Uaf9sy83u9a8uf9afnH1WHfahu9ha"
	Key = []byte(CreateHash(Password))
	block, _ := aes.NewCipher(Key)
	Gcm, _ = cipher.NewGCM(block)
}

// CreateHash creates hash
func CreateHash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// CreateHash2 creates hash using Password2
func CreateHash2(password string) string {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// Encrypt encrypts a []byte
func Encrypt(data []byte) (bool, []byte) {
	nonce := make([]byte, Gcm.NonceSize())
	ciphertext := Gcm.Seal(nonce, nonce, data, nil)
	return true, ciphertext
}

// Decrypt decrypts a []byte
func Decrypt(data []byte) (bool, []byte) {
	nonceSize := Gcm.NonceSize()
	if len(data) < nonceSize {
		return false, []byte{}
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := Gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, []byte{}
	}
	return true, plaintext
}

// SmallBase convert from 0-255 to 0-15
func SmallBase(B []byte) []byte {
	RES := make([]byte, 2*len(B), 2*len(B))
	for i := range B {
		RES[2*i] = (B[i] / 16) + 65
		RES[2*i+1] = (B[i] % 16) + 65
	}
	return RES
}

// BigBase convert from 0-15 to 0-255
func BigBase(B []byte) []byte {
	RES := make([]byte, len(B)/2, len(B)/2)
	for i := range RES {
		RES[i] = (B[2*i]-65)*16 + B[2*i+1] - 65
	}
	return RES
}

// RandomString generates random string of some length
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
