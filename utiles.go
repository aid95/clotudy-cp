package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	}
	return path, nil
}

// Md5HashGen 평문에 대한 md5 hash 생성
func Md5HashGen(plaintxt string) string {
	h := md5.New()
	io.WriteString(h, plaintxt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// RunCommandLine 명령행, 실행 인자를 받아 실행 후 표준 입력/출력/에러를 반환하는 함수.
func RunCommandLine(cmdline string, arg []string, input string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(cmdline, arg...)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()

	if err = cmd.Start(); err != nil { //Use start, not run
		fmt.Println("An error occured: ", err) //replace with logger, or anything you want
	}
	if len(input) > 0 {
		io.WriteString(stdin, input)
	}
	cmd.Wait()

	return stdout.String(), stderr.String(), err
}
