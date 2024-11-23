package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {

	bitSize := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	fmt.Println(string(privateKeyPEM))
	fmt.Println(string(publicKeyPEM))

	// 获取当前工作目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get caller information")
	}
	dir := filepath.Dir(filename)

	// 保存私钥到文件
	privateFile, err := os.Create(dir + "/private.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer privateFile.Close()

	if _, err := privateFile.Write(privateKeyPEM); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 保存公钥到文件
	publicFile, err := os.Create(dir + "/public.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer publicFile.Close()

	if _, err := publicFile.Write(publicKeyPEM); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
