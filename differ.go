package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Differ struct {
	workingDirectory     string
	tempWorkingDirectory string
	pdfA                 *PDF
	pdfB                 *PDF
	page                 int
}

func newDiffer(pdfA *PDF, pdfB *PDF, page int) *Differ {
	d := new(Differ)
	d.workingDirectory = d.createWorkingDirectory()
	d.tempWorkingDirectory = d.createTempDirectory(d.workingDirectory)
	d.pdfA = pdfA
	d.pdfB = pdfB
	d.page = page
	return d
}

func (d Differ) createWorkingDirectory() string {
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

func (d Differ) createTempDirectory(basePath string) string {
	workingPath := filepath.Join(basePath, "tmp")

	err := os.Mkdir(workingPath, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return workingPath
}

func (d Differ) imageDiff(pdfAImagePath string, pdfBImagePath string, count int) string {
	outputPath := filepath.Join(d.workingDirectory, fmt.Sprintf("diff-%03d.jpg", count))
	out, err := exec.Command("convert", pdfAImagePath, pdfBImagePath, "-compose", "Multiply", "-composite", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func (d Differ) Diff() {
	d.pdfA.ToImage(d.tempWorkingDirectory)
	d.pdfB.ToImage(d.tempWorkingDirectory)

	for i, _ := range d.pdfA.Images {
		var index int = i + 1
		if d.page == 0 || d.page == index { // 0 is all pages.
			fmt.Println("generate Diff Image :", d.imageDiff(d.pdfA.ToRed(i), d.pdfB.ToBlue(i), index))
		}
	}
}
