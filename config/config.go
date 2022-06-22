package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config interface {
	GetEffort() int
	GetQuality() int
	SetQuality(int) error
	SetEffort(int) error
	SetLosslessTranscoding(bool) error
	GetLosslessTranscoding() bool
}
type config struct{}

// SetEffort implements Config
func (*config) SetEffort(e int) error {
	if e < 1 || e > 9 {
		return errors.New("Effort must be between 1 and 9")
	}

	cfg := readConfig()
	cfg.Libjxl.Effort = e

	err := writeConfig(cfg)

	return err
}

// SetQuality implements Config
func (*config) SetQuality(q int) error {
	if q < 1 || q > 100 {
		return errors.New("Quality must be between 1 and 100")
	}

	cfg := readConfig()
	cfg.Libjxl.Quality = q

	err := writeConfig(cfg)

	return err
}

// SetLosslessTranscoding implements Config
func (*config) SetLosslessTranscoding(q bool) error {
	cfg := readConfig()
	cfg.Libjxl.LosslessTranscoding = q

	err := writeConfig(cfg)

	return err
}

// GetQuality implements Config
func (*config) GetQuality() int {
	return readConfig().Libjxl.Quality
}

// GetEffort implements Config
func (*config) GetEffort() int {
	return readConfig().Libjxl.Effort
}

// GetEffort implements Config
func (*config) GetLosslessTranscoding() bool {
	return readConfig().Libjxl.LosslessTranscoding
}

var (
	configFolder = filepath.Dir(os.Args[0])
	configPath   = filepath.Join(configFolder, "config.json")
)

func readConfig() Configuration {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		cfg := createStandardConfiguration()
		_ = writeConfig(cfg)
		return cfg
	}

	cfgData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return createStandardConfiguration()
	}

	var cfgUnmarshal Configuration
	json.Unmarshal(cfgData, &cfgUnmarshal)
	// TODO Check if the configuration is valid
	return cfgUnmarshal
}

func writeConfig(cfg Configuration) error {
	j, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, j, 0644)
}

func createStandardConfiguration() Configuration {
	return Configuration{
		Libjxl: Libjxl{
			Effort:  7,
			Quality: 80,
		},
	}
}

type Configuration struct {
	Libjxl Libjxl
}

type Libjxl struct {
	Effort              int
	Quality             int
	LosslessTranscoding bool
}

func NewConfig() Config {
	return &config{}
}
