package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

const version string = "0.0.3"

func ParseArguments(args []string) (string, string, int, error) {
	var pdfA string = ""
	var pdfB string = ""
	var page int = 0

	if len(args) >= 3 {
		pdfA = args[1]
		pdfB = args[2]
		if len(args) >= 4 {
			page, _ = strconv.Atoi(args[3])
		}
	} else {
		command := path.Base(args[0])
		return pdfA, pdfB, page, fmt.Errorf("$ %s foo.pdf bar.pdf [diff page number]\nVersion:%s", command, version)
	}

	files := []string{pdfA, pdfB}
	for _, v := range files {
		_, err := os.Stat(v)
		if err != nil {
			return pdfA, pdfB, page, err
		}
	}

	return pdfA, pdfB, page, nil
}

func main() {
	var pdfAPath string = ""
	var pdfBPath string = ""
	var page int = 0
	pdfAPath, pdfBPath, page, err := ParseArguments(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pdfA *PDF = newPDF(pdfAPath, "pdf_a")
	var pdfB *PDF = newPDF(pdfBPath, "pdf_b")

	var differ *Differ = newDiffer(pdfA, pdfB, page)

	differ.Diff()

	err = os.RemoveAll(differ.tempWorkingDirectory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if runtime.GOOS == "darwin" {
		exec.Command("open", differ.workingDirectory).CombinedOutput()
	}
}
