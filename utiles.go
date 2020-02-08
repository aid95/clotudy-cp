package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// GetSystemEnv 환경변수값을 가져오는데 없다면 기본값 반환
func GetSystemEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// MakePathDir 주어진 경로에 폴더 생성
func MakePathDir(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
	} else {
		return "", err
	}
	return path, nil
}

// Md5HashGen 평문에 대한 md5 hash 생성
func Md5HashGen(plaintxt string) string {
	h := md5.New()
	io.WriteString(h, plaintxt)
	return fmt.Sprintf("%x", h.Sum(nil))
}
