package util

import (
	"fmt"
	"net/http"
	"os"
)

func ErrorExit(e error, message string) {
	if e != nil {
		err := fmt.Sprintf("%s: the real error(%s)", message, e.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

func HttpError(e error, w http.ResponseWriter, message string) {
	if e != nil {
		err := fmt.Sprintf("%s: the real error(%s)", message, e.Error())
		http.Error(w, err, 500)
	}
}
