package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dhcgn/jxl-for-lightroom/config"
	"github.com/dhcgn/jxl-for-lightroom/converter"
	"github.com/gorilla/mux"
)

//go:embed assets/*
var content embed.FS

type Ui interface {
	ShowDialog([]string) error
}
type ui struct {
	Files             []File
	converter         converter.Converter
	config            config.Config
	isBusy            bool
	progress          chan int
	lastProgressValue int
}

type PageData struct {
	PageTitle        string
	Quality          string
	Effort           string
	TotalValidImages string
	Files            []File
}

type File struct {
	Name       string
	Path       string
	CanConvert bool
}

func (ui ui) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)

	data := PageData{
		PageTitle:        "jxl for lightroom",
		Files:            ui.Files,
		Effort:           fmt.Sprintf("%v", ui.config.GetEffort()),
		Quality:          fmt.Sprintf("%v", ui.config.GetQuality()),
		TotalValidImages: fmt.Sprintf("%d", len(ui.Files)),
	}

	t, err := template.ParseFS(content, "assets/index.html")
	if err != nil {
		log.Println(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func (ui ui) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("quality")
	e := r.FormValue("effort")
	log.Println("SettingsHandler", "Q:", q, "E:", e)

	if i, err := strconv.Atoi(q); err == nil {
		ui.config.SetQuality(i)
	}

	if i, err := strconv.Atoi(e); err == nil {
		ui.config.SetEffort(i)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (ui ui) ConvertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ConvertHandler")

	filenames := []string{}
	for _, f := range ui.Files {
		filenames = append(filenames, f.Path)
	}

	done, error := ui.converter.Convert(filenames, ui.progress)
	log.Println("ConvertHandler", "done:", done, "error:", error)
	if error != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	ui.isBusy = true
	go func() {
		done.Wait()
		ui.isBusy = false
		log.Println("ConvertHandler done")
	}()

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (u *ui) readFiles(files []string) {
	for _, file := range files {
		u.Files = append(u.Files, File{
			Name:       filepath.Base(file),
			Path:       file,
			CanConvert: u.converter.CanConvert(file),
		},
		)
	}
}

// ShowDialog implements Ui
func (ui ui) ShowDialog(files []string) error {

	ui.readFiles(files)

	r := mux.NewRouter()
	r.HandleFunc("/", ui.HomeHandler)
	r.HandleFunc("/settings", ui.SettingsHandler)
	r.HandleFunc("/convert", ui.ConvertHandler)

	htmlContent, _ := fs.Sub(content, "assets")
	fs := http.FileServer(http.FS(htmlContent))
	r.PathPrefix("/").Handler(fs)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:37700",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-time.After(100 * time.Millisecond)

	go func() {
		cmd := "cmd"
		args := []string{"/c", "start", "http://127.0.0.1:37700"}
		exec.Command(cmd, args...).Start()
	}()

	return nil
}

func NewUi(c converter.Converter, cfg config.Config) Ui {
	ui := &ui{
		converter: c,
		config:    cfg,
		progress:  make(chan int),
	}

	go func() {
		for {
			select {
			case p := <-ui.progress:
				ui.lastProgressValue = p
			case <-time.After(1 * time.Second):
				// DEBUG
			}
		}
	}()

	return ui
}
