package main

import (
	"github.com/mhvn092/movie-go/internal/util"
)

func main() {
	conn := util.InitDb()
	defer conn.Close()
}
