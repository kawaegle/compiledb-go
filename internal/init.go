package internal

import (
	"bufio"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	InputFile    string
	OutputFile   string
	BuildDir     string
	Exclude      string
	Macros       string
	RegexCompile string
	RegexFile    string
	CommandStyle bool
	FullPath     bool
	NoBuild      bool
	NoStrict     bool
}

var ParseConfig Config
var StatusCode int = 0

func WriteJSON(filename string, cmdCnt int, data *[]Command) {
	if cmdCnt == 0 {
		return
	}

	// format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding JSON:%v", err)
	}

	// write file
	if filename == "-" {
		println(string(jsonData))
	} else {
		outfile, err := os.Create(filename)
		if err != nil {
			log.Fatalf("create %v failed! err:%v", filename, err)
		}
		defer outfile.Close()

		_, err = outfile.Write(jsonData)
		if err != nil {
			log.Fatalf("write %v failed! err:%v", filename, err)
		}
		log.Infof("write %d entries to %s", cmdCnt, filename)
	}
}

func Generate() {
	var (
		buildLog []string
		scnner   *bufio.Scanner
		file     *os.File
		err      error
	)
	defer file.Close()

	if ParseConfig.InputFile != "stdin" {
		file, err = os.OpenFile(ParseConfig.InputFile, os.O_RDONLY, 0444)
		if err != nil {
			log.Fatalf("open %v failed!", ParseConfig.InputFile)
		}
		scnner = bufio.NewScanner(file)
		log.Debugf("Build from file")
	} else {
		scnner = bufio.NewScanner(os.Stdin)
		log.Debugf("Build from stdin")
	}

	scnner.Buffer(make([]byte, 1024*1024), 1024*1024*100)
	for scnner.Scan() {
		buildLog = append(buildLog, scnner.Text())
	}
	Parse(buildLog)
}
