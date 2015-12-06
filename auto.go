package main

import(
	"fmt"
	"bytes"
	"log"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/go-fsnotify/fsnotify"
)

var flags = []cli.Flag{}

func start(c *cli.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
			  if event.Op != fsnotify.Chmod {
				  commitByEvent(event)
			  }
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("watch start")
	<-done
}

func operationToString(op fsnotify.Op) string{
	var buffer bytes.Buffer

	if op&fsnotify.Create == fsnotify.Create {
		buffer.WriteString("|CREATE")
	}
	if op&fsnotify.Remove == fsnotify.Remove {
		buffer.WriteString("|REMOVE")
	}
	if op&fsnotify.Write == fsnotify.Write {
		buffer.WriteString("|WRITE")
	}
	if op&fsnotify.Rename == fsnotify.Rename {
		buffer.WriteString("|RENAME")
	}
	if op&fsnotify.Chmod == fsnotify.Chmod {
		buffer.WriteString("|CHMOD")
	}

	return buffer.String()[1:]
}

func commitByEvent(event fsnotify.Event) {
	commitFile(event.Name, fmt.Sprintf("fixup! %s %s", event.Name, operationToString(event.Op)))
}

func commitFile(filepath string, message string) {
	if gitAdd(filepath) {
		gitCommit(message)
	}
}

func gitCommit(message string) {
	_, err := exec.Command("git", "commit", "-m", message).CombinedOutput()
	if err != nil {
		log.Println("git commit fall", err, message)
	}else{
		log.Println("commit ", message)
	}
}

func gitAdd(filepath string) bool {
	_, err := exec.Command("git", "add", filepath).CombinedOutput()
	return err == nil
}