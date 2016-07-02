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

func (p *PDF) toImage(tmpPath string) {
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

func (p PDF) toRed(index int) string {
	return p.toColorImage("Red", p.Images[index])
}

func (p PDF) toBlue(index int) string {
	return p.toColorImage("Blue", p.Images[index])
}

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

	pdfA.toImage(tmpPath)
	pdfB.toImage(tmpPath)

	for i, _ := range pdfA.Images {
		var index int = i + 1
		if page == 0 || page == index { // 0 is all pages.
			fmt.Println("generate Diff Image :", Diff(pdfA.toRed(i), pdfB.toBlue(i), path, index))
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
