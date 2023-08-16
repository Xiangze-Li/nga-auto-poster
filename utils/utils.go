package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func ExitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ReadAndPop(filename string, split string) string {
	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	ExitOnError(err)
	defer f.Close()

	r := bufio.NewReader(f)

	done := false
	var bdContent, bdRest strings.Builder
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			ExitOnError(err)
		}
		if !done {
			if strings.TrimSpace(line) == split {
				done = true
			} else {
				bdContent.WriteString(line)
			}
		} else {
			bdRest.WriteString(line)
		}
	}

	content := bdContent.String()
	rest := bdRest.String()

	f.Seek(0, 0)
	f.WriteString(rest)
	f.Truncate(int64(len(rest)))

	return content
}
