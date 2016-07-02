package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type PDF struct {
	path   string
	prefix string
	Images []string
}

func newPDF(path string, prefix string) *PDF {
	p := new(PDF)
	p.path = path
	p.prefix = prefix
	return p
}

func (p *PDF) ToImage(tmpPath string) {
	gsOption := []string{
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=jpeg",
		"-r150",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-dMaxStripSize=8192",
	}
	gsOption = append(gsOption, fmt.Sprintf("-sOutputFile=%s/%s_%s.jpg", tmpPath, p.prefix, "%04d"))
	gsOption = append(gsOption, p.path)

	_, err := exec.Command("gs", gsOption...).Output()

	if err != nil {
		fmt.Println("ToImage error", err)
		os.Exit(1)
	}
	p.Images, _ = filepath.Glob(fmt.Sprintf("%s/%s*.jpg", tmpPath, p.prefix))
}

func (p PDF) toGrayScale(path string) string {
	outputPath := fmt.Sprintf("%s.gray.jpg", path)
	out, err := exec.Command("convert", path, "-type", "GrayScale", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("ToGrayScale error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func (p PDF) toColorImage(color string, path string) string {
	grayPath := p.toGrayScale(path)
	outputPath := fmt.Sprintf("%s.%s.jpg", grayPath, color)
	colorOption := fmt.Sprintf("%s,White", color)

	out, err := exec.Command("convert", grayPath, "+level-colors", colorOption, outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("ToRed error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func (p PDF) ToRed(index int) string {
	return p.toColorImage("Red", p.Images[index])
}

func (p PDF) ToBlue(index int) string {
	return p.toColorImage("Blue", p.Images[index])
}
