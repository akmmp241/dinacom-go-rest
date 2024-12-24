package helpers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mime/multipart"
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

func Base64ToString(s string) string {
	b, _ := base64.StdEncoding.DecodeString(s)

	return string(b)
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
	cipherText := Base64ToString(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, []byte(cipherText))
	return string(plainText), nil
}

func UploadToGemini(ctx context.Context, client *genai.Client, file multipart.File, mimeType string) (string, error) {
	options := genai.UploadFileOptions{
		DisplayName: "uploaded-image",
		MIMEType:    mimeType,
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		return "", err
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI, nil
}

func UploadS3(ctx context.Context, uploader *manager.Uploader, file multipart.File, fileName string, bucket string) (string, error) {
	uploadedFile, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	log.Printf("Uploaded fileas: %s", uploadedFile.Location)
	return uploadedFile.Location, nil
}
