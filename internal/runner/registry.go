package runner

type LanguageDef struct {
	SourceExt        string
	BuildCmd         []string // nil for interpreted languages
	BuildArtifact    string   // artifact filename (empty if no build)
	RunCmd           []string
	RunNeedsArtifact bool
}

var registry = map[string]LanguageDef{
	"py3": {
		SourceExt:        ".py",
		BuildCmd:         nil,
		BuildArtifact:    "",
		RunCmd:           []string{"python3", "{source_dir}/{source_file}"},
		RunNeedsArtifact: false,
	},
	"cpp": {
		SourceExt:        ".cpp",
		BuildCmd:         []string{"g++", "-o", "{artifact_dir}/{artifact_name}", "{source_dir}/{source_file}"},
		BuildArtifact:    "prog",
		RunCmd:           []string{"{artifact_dir}/{artifact_name}"},
		RunNeedsArtifact: true,
	},
}

func Lookup(lang string) (LanguageDef, bool) {
	def, ok := registry[lang]
	return def, ok
}

func SupportedLanguages() []string {
	langs := make([]string, 0, len(registry))
	for l := range registry {
		langs = append(langs, l)
	}
	return langs
}
