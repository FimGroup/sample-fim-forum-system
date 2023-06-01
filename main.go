package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/FimGroup/sample-fim-forum-system/forumcore"
)

//go:embed version
var ver string

func main() {
	fmt.Println("version:", strings.TrimSpace(string(ver)))

	if err := forumcore.StartForum(); err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	fmt.Println("working directory:", wd, err)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT)
	_ = <-c
	fmt.Println("service exit!")
}
