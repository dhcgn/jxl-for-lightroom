package main

import (
	"fmt"
	"os"

	"github.com/dhcgn/jxl-for-lightroom/config"
	"github.com/dhcgn/jxl-for-lightroom/converter"
	"github.com/dhcgn/jxl-for-lightroom/ui"
)

var (
	Version = "0.0.1-alpha"
)

func main() {
	printIntro()
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	c := converter.NewConvertor()
	cfg := config.NewConfig()

	u := ui.NewUi(c, cfg)
	err := u.ShowDialog(os.Args[1:])

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	select {}
}

func printHelp() {
	fmt.Println("help")
}

func printIntro() {
	fmt.Println("jxl for lightroom", Version)
}
