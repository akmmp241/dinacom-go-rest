package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func VerifyPassword(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func StringToBase64(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

func Base64ToByte(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)

	return b
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return StringToBase64(cipherText), nil
}

func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Base64ToByte(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
