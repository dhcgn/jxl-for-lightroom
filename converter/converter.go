package converter

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	allowedExtensionsMap = map[string]struct{}{".png": {}, ".apng": {}, ".gif": {}, ".jpeg": {}, ".jpg": {}, ".ppm": {}, ".pfm": {}, ".pgx": {}}
)

//go:embed assets/*
var content embed.FS

type Converter interface {
	Convert([]string, int, int, bool, chan<- int, chan<- EncodeResult) (sync.WaitGroup, error)
	CanConvert(string) bool
}
type converter struct{}

// CanConvert implements Converter
func (*converter) CanConvert(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	_, found := allowedExtensionsMap[ext]
	return found
}

// Convert implements Converter
func (c *converter) Convert(files []string, e int, q int, lt bool, progress chan<- int, results chan<- EncodeResult) (sync.WaitGroup, error) {
	convertFiles := make([]string, 0, len(files))
	for _, file := range files {
		if c.CanConvert(file) {
			convertFiles = append(convertFiles, file)
		}
	}

	if len(convertFiles) == 0 {
		return sync.WaitGroup{}, errors.New("No files to convert")
	}

	encoderPath := extractEncoder(filepath.Dir(os.Args[0]))
	log.Println("Using encoder:", encoderPath)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg sync.WaitGroup, progress chan<- int) {

		for i, file := range convertFiles {
			r := encode(encoderPath, file, e, q, lt)
			var p float64 = (float64(i) + 1.0) / float64(len(convertFiles)) * 100.0
			log.Println("Progress", p)
			progress <- int(p)
			results <- r
		}

		// for i := 0; i <= 100; i++ {
		// 	log.Println("Convert", "progress:", i)

		// 	progress <- i
		// 	<-time.After(20 * time.Millisecond)
		// }

		wg.Done()
	}(wg, progress)

	return wg, nil
}

type String string

func (p String) isJpegExtension() bool {
	ext := strings.ToLower(filepath.Ext(string(p)))
	if ext == ".jpeg" || ext == ".jpg" {
		return true
	}
	return false
}

type EncodeResult struct {
	Input   string
	Output  string
	Stdout  string
	Stderr  string
	Error   string
	Elasped string
}

func encode(encoderPath, file string, e, q int, lt bool) EncodeResult {
	jxlPath := filepath.Join(filepath.Dir(file), filepath.Base(file)+".jxl")

	var cmd *exec.Cmd
	if String(file).isJpegExtension() && lt {
		cmd = exec.Command(encoderPath, file, jxlPath)
	} else {
		cmd = exec.Command(encoderPath, file, jxlPath, "-q", strconv.Itoa(q), "-e", strconv.Itoa(e))
	}

	var sdtout bytes.Buffer
	cmd.Stdout = &sdtout
	var sdterr bytes.Buffer
	cmd.Stderr = &sdterr

	start := time.Now()
	err := cmd.Run()
	end := time.Now()

	er := EncodeResult{
		Input:   file,
		Output:  jxlPath,
		Stdout:  sdtout.String(),
		Stderr:  sdterr.String(),
		Elasped: end.Sub(start).String(),
	}

	if err != nil {
		log.Fatal(err)
		er.Error = err.Error()
	}
	return er
}

func extractEncoder(f string) string {
	fs, err := content.Open("assets/cjxl.exe")
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	filepath := filepath.Join(f, "cjxl.exe")
	data, _ := io.ReadAll(fs)

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		tempFolder, err := os.MkdirTemp("", "jxl-for-lightroom")
		if err != nil {
			panic(err)
		}
		return extractEncoder(tempFolder)
	}

	return filepath
}

func NewConvertor() Converter {
	return &converter{}
}
