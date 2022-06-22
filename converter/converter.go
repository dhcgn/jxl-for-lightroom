package converter

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	allowedExtensionsMap = map[string]struct{}{".png": {}, ".apng": {}, ".gif": {}, ".jpeg": {}, ".jpg": {}, ".ppm": {}, ".pfm": {}, ".pgx": {}}
)

type Converter interface {
	Convert([]string, chan<- int) (sync.WaitGroup, error)
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
func (c *converter) Convert(files []string, progress chan<- int) (sync.WaitGroup, error) {
	convertFiles := make([]string, 0, len(files))
	for _, file := range files {
		if c.CanConvert(file) {
			convertFiles = append(convertFiles, file)
		}
	}

	if len(convertFiles) == 0 {
		return sync.WaitGroup{}, errors.New("No files to convert")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg sync.WaitGroup, progress chan<- int) {
		for i := 0; i <= 100; i++ {
			log.Println("Convert", "progress:", i)
			progress <- i
			<-time.After(20 * time.Millisecond)
		}

		wg.Done()
	}(wg, progress)

	return wg, nil
}

func NewConvertor() Converter {
	return &converter{}
}
