package util

import (
	"os"
	"strings"
)

var EnvKeys = struct {
	DATABASE_URL string
	HOST         string
	PORT         string
}{
	DATABASE_URL: "DATABASE_URL",
	HOST:         "HOST",
	PORT:         "PORT",
}
var envValues = make(map[string]string)

func readFile(filename string) string {
	body, err := os.ReadFile(filename)
	if err != nil {
		ErrorExit(err, "could not read file")
	}
	return string(body)
}

func ReadEnv() {
	pwd, _ := os.Getwd()
	file := readFile(pwd + "/.env")
	keys := strings.Split(file, "\r\n")
	for _, key := range keys {
		keypair := strings.Split(key, "=")
		envValues[keypair[0]] = keypair[1]
	}
}

func GetEnv(key string) string {
	return envValues[key]
}
