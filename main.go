package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var counter uint

func main() {
	counter = 0

	dir := parseArgs()
	err := processMedia(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs() string {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  COMMAND <dir>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	dir := flag.Args()[0]

	return dir
}

func processMedia(dir string) error {
	numCpu := runtime.NumCPU()

	fpathCh := make(chan string, numCpu)
	counterCh := make(chan bool)

	fileHandlersDone := make([]<-chan bool, cap(fpathCh))

	for i := 0; i < cap(fpathCh); i++ {
		fileHandlersDone[i] = startFileHandler(fpathCh, counterCh)
	}
	counterDone := startCounter(counterCh)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// continue
			return nil
		}

		fpathCh <- path
		return nil
	}

	err := filepath.Walk(dir, walkFn)
	if err != nil && err != io.EOF {
		return err
	}

	close(fpathCh)
	for _, fileHandlerDone := range fileHandlersDone {
		<-fileHandlerDone
	}

	close(counterCh)
	<-counterDone

	fmt.Println("Files parsed:", counter)
	return nil
}

func startCounter(c <-chan bool) <-chan bool {
	done := make(chan bool)
	go func() {
		for {
			_, more := <-c
			if !more {
				done <- true
				return
			}

			counter += 1
		}
	}()
	return done
}

func startFileHandler(c <-chan string, counter chan<- bool) <-chan bool {
	done := make(chan bool)
	go func() {
		for {
			fpath, more := <-c
			if !more {
				done <- true
				return
			}

			m, err := NewMedia(fpath)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("%v: Error %v\n", fpath, err)
				}
				continue
			}

			if stamp := m.DateStamp(); stamp != "" {
				counter <- true
				moveFileToDir(m.fpath, stamp)
			}
		}
	}()
	return done
}

func moveFileToDir(src, stamp string) {
	dir := filepath.Join(".", "results", stamp)
	os.MkdirAll(dir, os.ModePerm)

	base := filepath.Base(src)
	dst := filepath.Join(dir, base)

	err := os.Rename(src, dst)
	if err != nil {
		fmt.Printf("%v -> %v: Error %v\n", src, dir, err)
	}
}
