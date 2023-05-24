package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type ImageFile struct {
	Path    string
	ModTime time.Time
}

type PageData struct {
	Images   []ImageFile
	Page     int
	PrevPath string
	NextPath string
	Pages    int
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
}

func main() {
	imagesPerPage := flag.Int("p", 100, "Number of images per page")
	flag.Parse()

	args := flag.Args()
	if len(args) != 3 {
		fmt.Println("Usage: main -p [NUMBER] [TEMPLATE_PATH] [INPUT_DIRECTORY] [OUTPUT_DIRECTORY]")
		os.Exit(1)
	}

	tmplPath := args[0]
	inputDir := args[1]
	outputDir := args[2]

	files, err := os.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	imageFiles := make([]ImageFile, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".jpg" && filepath.Ext(file.Name()) != ".png" {
			continue
		}
		randBytes := make([]byte, 16)
		_, err := rand.Read(randBytes)
		if err != nil {
			panic(err)
		}
		randName := fmt.Sprintf("%x", randBytes) + filepath.Ext(file.Name())

		src := filepath.Join(inputDir, file.Name())
		dst := filepath.Join(outputDir, "images", randName)
		err = copyFile(src, dst)
		if err != nil {
			panic(err)
		}

		fileInfo, err := file.Info()
		if err != nil {
			panic(err)
		}

		imageFiles = append(imageFiles, ImageFile{Path: filepath.Join("images", randName), ModTime: fileInfo.ModTime()})
	}

	sort.Slice(imageFiles, func(i, j int) bool {
		return imageFiles[i].ModTime.After(imageFiles[j].ModTime)
	})

	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	pages := (len(imageFiles) + *imagesPerPage - 1) / *imagesPerPage

	for i := 0; i < pages; i++ {
		start := i * *imagesPerPage
		end := start + *imagesPerPage
		if end > len(imageFiles) {
			end = len(imageFiles)
		}
		pageData := PageData{
			Images:   imageFiles[start:end],
			PrevPath: generateFilename(i),
			Page:     i + 1,
			NextPath: generateFilename(i + 2),
			Pages:    pages,
		}

		filename := fmt.Sprintf("index%d.html", i+1)
		if i == 0 {
			filename = "index.html"
		}
		htmlFile, err := os.Create(filepath.Join(outputDir, filename))
		if err != nil {
			panic(err)
		}
		defer htmlFile.Close()

		err = tmpl.Execute(htmlFile, pageData)
		if err != nil {
			panic(err)
		}
	}
}

func generateFilename(page int) string {
	if page == 1 {
		return "index.html"
	}
	return fmt.Sprintf("index%d.html", page)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	err = os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
	return err
}
