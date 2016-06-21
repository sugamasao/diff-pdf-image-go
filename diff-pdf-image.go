package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func CreateWorkingDirectory() string {
	currentPath, _ := filepath.Abs(".")
	tmpName := time.Now().Unix()
	workingPath := filepath.Join(currentPath, fmt.Sprintf("pdf-diff-%d", tmpName))

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

func ToImage(targetType string, targetPDF string, outputPath string) []string {
	gsOption := []string{
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=jpeg",
		"-r150",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-dMaxStripSize=8192",
	}
	gsOption = append(gsOption, fmt.Sprintf("-sOutputFile=%s/%s_%s.jpg", outputPath, targetType, "%d"))
	gsOption = append(gsOption, targetPDF)

	_, err := exec.Command("gs", gsOption...).Output()

	if err != nil {
		fmt.Println("ToImage error", err)
		os.Exit(1)
	}

	list, _ := filepath.Glob(fmt.Sprintf("%s/%s*.jpg", outputPath, targetType))

	return list
}

func ToGrayScale(path string) string {
	outputPath := fmt.Sprintf("%s.gray.jpg", path)
	fmt.Println("convert", path, "-type GrayScale", outputPath)
	out, err := exec.Command("convert", path, "-type", "GrayScale", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("ToGrayScale error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func ToRed(path string) string {
	gray_path := ToGrayScale(path)

	outputPath := fmt.Sprintf("%s.red.jpg", gray_path)
	out, err := exec.Command("convert", gray_path, "+level-colors", "Red,White", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("ToRed error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func ToBlue(path string) string {
	gray_path := ToGrayScale(path)

	outputPath := fmt.Sprintf("%s.red.jpg", gray_path)
	out, err := exec.Command("convert", gray_path, "+level-colors", "Blue,White", outputPath).CombinedOutput()

	if err != nil {
		fmt.Println("ToBlue error", err, string(out))
		os.Exit(1)
	}

	return outputPath
}

func Diff(pdf_a string, pdf_b string, workingPath string, count int) string {
	pdf_a_path := ToRed(pdf_a)
	pdf_b_path := ToBlue(pdf_b)

	outputPath := filepath.Join(workingPath, fmt.Sprintf("diff-%03d.jpg", count))
	out, err := exec.Command("convert", pdf_a_path, pdf_b_path, "-compose", "Multiply", "-composite", outputPath).CombinedOutput()

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
		pdfA = os.Args[1]
		pdfB = os.Args[2]
		if os.Args[3] != "" {
			page, _ = strconv.Atoi(os.Args[3])
		}
	} else {
		command := path.Base(args[0])
		return pdfA, pdfB, page, fmt.Errorf("%s pdf_a pdf_b [diff page number]", command)
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
	var pdfA string = ""
	var pdfB string = ""
	var page int = 0
	pdfA, pdfB, page, err := ParseArguments(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var path string = CreateWorkingDirectory()
	var tmpPath string = CreateTempDirectory(path)
	pdfAList := ToImage("pdf_a", pdfA, tmpPath)
	pdfBList := ToImage("pdf_b", pdfB, tmpPath)

	for i, v := range pdfAList {
		var index int = i + 1
		if page == 0 { // all page diff.
			fmt.Println("generate Diff Image :", Diff(v, pdfBList[i], path, index))
		} else if page == index { // diff only in page number.
			fmt.Println("generate Diff Image :", Diff(v, pdfBList[i], path, index))
		}
	}

	err = os.RemoveAll(tmpPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
