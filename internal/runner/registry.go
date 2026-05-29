package runner

import (
	"os"

	"github.com/thesouldev/goboxd/internal/types"
	"gopkg.in/yaml.v3"
)

type LanguageConfig struct {
	ID             string        `yaml:"id"`
	Name           string        `yaml:"name"`
	SourceFilename string        `yaml:"source_filename"`
	Artifact       string        `yaml:"artifact"`
	Build          *ActionConfig `yaml:"build"`
	Run            ActionConfig  `yaml:"run"`
}

type ActionConfig struct {
	Cmd           string        `yaml:"cmd"`
	Args          []string      `yaml:"args"`
	Limits        *types.Limits `yaml:"limits"`
	FlagAllowlist []string      `yaml:"flag_allowlist"`
}

type RegistryData struct {
	Languages []LanguageConfig `yaml:"languages"`
}

var registry = make(map[string]LanguageConfig)
var languageList []LanguageConfig

func InitRegistry(yamlPath string) error {
	b, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	var data RegistryData
	if err := yaml.Unmarshal(b, &data); err != nil {
		return err
	}

	newReg := make(map[string]LanguageConfig)
	for _, l := range data.Languages {
		newReg[l.ID] = l
	}
	registry = newReg
	languageList = data.Languages
	return nil
}

func Lookup(lang string) (LanguageConfig, bool) {
	def, ok := registry[lang]
	return def, ok
}

func SupportedLanguages() []LanguageConfig {
	return languageList
}
