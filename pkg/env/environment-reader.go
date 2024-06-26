package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	DATABASE_URL = "DATABASE_URL"
	HOST         = "HOST"
	PORT         = "PORT"
)

var envValues = make(map[string]string)

func readFile(filename string) string {
	body, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("could not read env file")
		os.Exit(1)
	}
	return string(body)
}

func ReadEnv() {
	pwd, _ := os.Getwd()
	file := readFile(pwd + "/.env")
	keys := strings.Split(file, "\r\n")
	for _, key := range keys {
		if key == "" || strings.HasPrefix(key, "#") {
			continue
		}
		keypair := strings.Split(key, "=")
		envValues[keypair[0]] = keypair[1]
	}
}

func GetEnv(key string) string {
	return envValues[key]
}
