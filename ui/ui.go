package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dhcgn/jxl-for-lightroom/config"
	"github.com/dhcgn/jxl-for-lightroom/converter"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
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
	encoderResults    chan converter.EncodeResult
	lastProgressValue int
	Log               string
}

type PageData struct {
	PageTitle           string
	Quality             string
	LosslessTranscoding bool
	Effort              string
	TotalValidImages    string
	Files               []File
}

type File struct {
	Name       string
	Path       string
	CanConvert bool
}

func (ui ui) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)

	data := PageData{
		PageTitle:           "jxl for lightroom",
		Files:               ui.Files,
		Effort:              fmt.Sprintf("%v", ui.config.GetEffort()),
		Quality:             fmt.Sprintf("%v", ui.config.GetQuality()),
		LosslessTranscoding: ui.config.GetLosslessTranscoding(),
		TotalValidImages:    fmt.Sprintf("%d", len(ui.Files)),
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
	lt := r.FormValue("losslesstranscoding")
	log.Println("SettingsHandler", "Q:", q, "E:", e, "LT:", lt)

	if i, err := strconv.Atoi(q); err == nil {
		ui.config.SetQuality(i)
	}

	if i, err := strconv.Atoi(e); err == nil {
		ui.config.SetEffort(i)
	}

	if i, err := strconv.ParseBool(lt); err == nil {
		ui.config.SetLosslessTranscoding(i)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (ui ui) ConvertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ConvertHandler")

	filenames := []string{}
	for _, f := range ui.Files {
		filenames = append(filenames, f.Path)
	}

	done, error := ui.converter.Convert(filenames, ui.config.GetEffort(), ui.config.GetQuality(), ui.config.GetLosslessTranscoding(), ui.progress, ui.encoderResults)
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

func getFreeTcpPort() int {
	port := 49152

	for i := 0; i < 24; i++ {

		if a, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%v", port)); err == nil {
			var l *net.TCPListener
			if l, err = net.ListenTCP("tcp", a); err == nil {
				defer l.Close()
				port := l.Addr().(*net.TCPAddr).Port

				return port
			}
		}
		port++
	}

	panic("NOT FREE PORT TO START SERVER")
}

// ShowDialog implements Ui
func (ui *ui) ShowDialog(files []string) error {

	ui.readFiles(files)

	r := mux.NewRouter()
	r.HandleFunc("/", ui.HomeHandler)
	r.HandleFunc("/settings", ui.SettingsHandler)
	r.HandleFunc("/convert", ui.ConvertHandler)
	r.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", ui.lastProgressValue)
	})
	r.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", ui.Log)
	})

	htmlContent, _ := fs.Sub(content, "assets")
	fs := http.FileServer(http.FS(htmlContent))
	r.PathPrefix("/").Handler(fs)

	port := getFreeTcpPort()

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-time.After(100 * time.Millisecond)

	go func() {
		cmd := "cmd"
		args := []string{"/c", "start", fmt.Sprintf("http://127.0.0.1:%v", port)}
		exec.Command(cmd, args...).Start()
	}()

	return nil
}

func NewUi(c converter.Converter, cfg config.Config) Ui {
	u := &ui{
		converter:      c,
		config:         cfg,
		progress:       make(chan int),
		encoderResults: make(chan converter.EncodeResult),
	}

	go func(u *ui) {
		for {
			select {
			case p := <-u.progress:
				log.Println("set lastProgressValue with:", p)
				u.lastProgressValue = p
			case p := <-u.encoderResults:
				data, _ := yaml.Marshal(p)
				u.Log += fmt.Sprintf("%s\n", data)
				log.Println(u.Log)
			case <-time.After(1 * time.Second):
				// DEBUG
			}
		}
	}(u)

	return u
}
