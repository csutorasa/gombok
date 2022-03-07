package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var errorLogger *log.Logger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
var infoLogger *log.Logger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
var debugLogger *log.Logger = log.New(discard{}, "[DEBUG] ", log.LstdFlags)
var filenameFormatter string = "%s_gombok.go"

func main() {
	quiet := flag.Bool("quiet", false, "Sets logger to none")
	info := flag.Bool("info", false, "Sets logger to info")
	debug := flag.Bool("debug", false, "Sets logger to debug")
	filename := flag.String("filename", "%s_gombok.go", "Sets the generated file name formatter")
	flag.Parse()
	if (*quiet && *info) || (*info && *debug) || (*quiet && *debug) {
		fmt.Printf("You cannot use multiple logger levels! Only one or none of these: quiet info debug")
		os.Exit(2)
	}
	if *quiet {
		errorLogger = log.New(discard{}, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(discard{}, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(discard{}, "[DEBUG] ", log.LstdFlags)
	} else if *info {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(discard{}, "[DEBUG] ", log.LstdFlags)
	} else if *debug {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
	} else {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(discard{}, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(discard{}, "[DEBUG] ", log.LstdFlags)
	}
	filenameFormatter = *filename

	d, err := filepath.Abs(".")
	if err != nil {
		errorLogger.Printf("Cannot find absolute path %v", err)
		d = "."
	}
	err = processDirRecursive(d)
	if err != nil {
		errorLogger.Println(err)
		os.Exit(1)
	}
}

type discard struct{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discard) WriteString(s string) (int, error) {
	return len(s), nil
}
