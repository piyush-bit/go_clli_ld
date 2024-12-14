package myutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)


func CompileAndListenCpp(file string, delay *int) error {
	// get absolute path of the file
	absPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	err = watcher.Add(absPath)
	if err != nil {
		return err
	}
	fmt.Println("Watching for changes in ", absPath)
	// run indefinitely


	for event := range watcher.Events {
		if event.Op == fsnotify.Write {
			er := CompileCpp(file)
			if er != nil {
				fmt.Println("Error: ", er)
			}

			if delay != nil {

				time.Sleep(1000 * time.Millisecond)

				// flush all the events in the channel
				go func() {
					for {
						select {
						case <-time.After(1000 * time.Millisecond):
							// return from the for loop
							return

						case event := <-watcher.Events:
							if event.Op == fsnotify.Write {
								continue
							}
						}
					}
				}()

			}
		}
	}

	return nil
}

func CompileCpp(file string) error {
	// clear the console
	fmt.Print("\033[H\033[2J")


	// fmt.Println("Problems :")

	cmdexec := exec.Command("g++", "-o", "out", file)
	cmdexec.Env = append(os.Environ(), "PATH=/usr/bin:/usr/local/bin")

	// print the output  of runnong the command
	cmdexec.Stdout = os.Stdout
	cmdexec.Stderr = os.Stderr

	err := cmdexec.Run()
	if err != nil {
		return err
	}
	return nil
}
