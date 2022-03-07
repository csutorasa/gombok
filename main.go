package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var errorLogger *log.Logger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
var infoLogger *log.Logger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
var debugLogger *log.Logger = log.New(io.Discard, "[DEBUG] ", log.LstdFlags)
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
		errorLogger = log.New(io.Discard, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(io.Discard, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(io.Discard, "[DEBUG] ", log.LstdFlags)
	} else if *info {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(io.Discard, "[DEBUG] ", log.LstdFlags)
	} else if *debug {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(os.Stderr, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
	} else {
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		infoLogger = log.New(io.Discard, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(io.Discard, "[DEBUG] ", log.LstdFlags)
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
