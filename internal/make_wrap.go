package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func MakeWrap(args []string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// append log
		args = append([]string{"-Bnkw"}, args...)
		cmd := exec.Command("make", args...)

		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stdoutBuf
		cmd.Run()

		level := log.GetLevel()

		// only print make log
		if ParseConfig.NoBuild == false {
			log.SetLevel(log.PanicLevel)
		}

		buildLog := strings.Split(stdoutBuf.String(), "\n")
		Parse(buildLog)

		// restore log level
		if ParseConfig.NoBuild == false {
			log.SetLevel(level)
		}

		wg.Done()
	}()

	if ParseConfig.NoBuild == false {
		cmd := exec.Command("make", args...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("stdout Error:", err)
			goto out
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("stderr Error:", err)
			goto out
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("start Error:", err)
			goto out
		}

		go TransferPrintScanner(stdout)
		go TransferPrintScanner(stderr)

		if err := cmd.Wait(); err != nil {
			StatusCode = cmd.ProcessState.ExitCode()
			fmt.Printf("make failed! errorCode: %d\n", StatusCode)
		}
	}

out:
	wg.Wait()
}
