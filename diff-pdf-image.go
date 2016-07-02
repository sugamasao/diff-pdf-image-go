package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const version string = "0.0.3"

func CreateWorkingDirectory() string {
	currentPath, _ := filepath.Abs(".")
	tmpName := time.Now().Unix()
	workingPath := filepath.Join(currentPath, fmt.Sprintf("diff-pdf-result-%d", tmpName))

	err := os.Mkdir(workingPath, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return workingPath
}

func CreateTempDirectory(basePath string) string {
	workingPath := filepath.Join(basePath, "tmp")

	err := os.Mkdir(workingPath, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return workingPath
}

func Diff(pdf_a string, pdf_b string, workingPath string, count int) string {
	outputPath := filepath.Join(workingPath, fmt.Sprintf("diff-%03d.jpg", count))
	out, err := exec.Command("convert", pdf_a, pdf_b, "-compose", "Multiply", "-composite", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

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

	var path string = CreateWorkingDirectory()
	var tmpPath string = CreateTempDirectory(path)

	pdfA.ToImage(tmpPath)
	pdfB.ToImage(tmpPath)

	for i, _ := range pdfA.Images {
		var index int = i + 1
		if page == 0 || page == index { // 0 is all pages.
			fmt.Println("generate Diff Image :", Diff(pdfA.ToRed(i), pdfB.ToBlue(i), path, index))
		}
	}

	err = os.RemoveAll(tmpPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if runtime.GOOS == "darwin" {
		exec.Command("open", path).CombinedOutput()
	}
}
