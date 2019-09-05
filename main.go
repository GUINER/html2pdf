package main

import (
	"bytes"
	"flag"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	htmlSubfix = "html"
	pdfSubfix  = "pdf"
	outputPath = "./pdf/"
)

// todo: 需要先安装wkhtmltopdf工具
// $sudo apt install wkhtmltopdf

type Pdfer struct {
	FileName string // 源文件名,html
	Output   string // 输出文件,pdf
}

func NewPdfer(file, out string) *Pdfer {
	return &Pdfer{FileName: file, Output: out}
}

// 读文件的方法
func (p *Pdfer) Generate() {
	log.Println("html2pdf file: " + p.FileName)

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	html, err := ioutil.ReadFile(p.FileName)
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(html)))
	if err = pdfg.Create(); err != nil {
		log.Fatal(err)
	}

	if err = pdfg.WriteFile(p.Output); err != nil {
		log.Fatal(err)
	}

	log.Println("Done file: " + p.Output)
}

// 读目录/文件/url
func readFile(path string) string {
	f, err := os.Stat(path)
	if err != nil {
		panic("openfile error" + err.Error())
	}

	if f.IsDir() {
		if path[len(path)-1:] == "/" {
			return path
		}
		return path + "/"
	}

	return path
}

func getEnv() string {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		panic("param is empty")
	}
	return args[0]
}

func main() {

	path := readFile(getEnv())

	dir, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dir.Close()

	fileInfo, _ := dir.Readdir(-1)
	//操作系统指定的路径分隔符
	//separator := string(os.PathSeparator)

	os.Mkdir(outputPath, os.ModePerm)

	for _, f := range fileInfo {
		baseName := strings.Replace(f.Name(), htmlSubfix, pdfSubfix, 1)
		output := outputPath + baseName
		NewPdfer(path+f.Name(), output).Generate()
	}
}
