package types

type RunRequest struct {
	Language         string      `json:"language"`
	Source           string      `json:"source"`
	SourceFilename   string      `json:"source_filename,omitempty"`
	ArtifactFilename string      `json:"artifact_filename,omitempty"`
	Build            *BuildBlock `json:"build,omitempty"`
	Run              *RunBlock   `json:"run,omitempty"`
	Tests            []TestCase  `json:"tests"`
}

type BuildBlock struct {
	Limits *Limits  `json:"limits,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}

type RunBlock struct {
	Limits *Limits  `json:"limits,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}

type Limits struct {
	WallTimeS    int `json:"wall_time_s,omitempty"`
	MemoryKB     int `json:"memory_kb,omitempty"`
	MaxProcesses int `json:"max_processes,omitempty"`
}

type TestCase struct {
	Stdin          string `json:"stdin,omitempty"`
	ExpectedStdout string `json:"expected_stdout,omitempty"`
}

type RunResponse struct {
	Status string       `json:"status"`
	Build  BuildResult  `json:"build"`
	Tests  []TestResult `json:"tests"`
}

type BuildResult struct {
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMS int    `json:"duration_ms"`
}

type TestResult struct {
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMS   int    `json:"duration_ms"`
	MemoryPeakKB int    `json:"memory_peak_kb"`
}
