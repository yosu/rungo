package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/go-fsnotify/fsnotify"
)

func main() {
	app := cli.NewApp()
	app.Name = "rungo"
	app.Version = Version
	app.Usage = "Run go files on modify"
	app.Author = "yosu"
	app.Email = "woodstock830@gmail.com"
	app.Action = doMain
	app.Run(os.Args)
}

func doMain(c *cli.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if isGoFile(event.Name) {
						execCommand("go", "run", event.Name)
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func isGoFile(name string) bool {
	match, _ := regexp.MatchString("\\.go$", name)
	return match
}

func execCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
