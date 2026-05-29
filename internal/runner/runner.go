package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/types"
)

const workDirPrefix = "goboxd-run-"

type Runner struct {
	useNsjail bool
}

func New() *Runner {
	r := &Runner{}
	r.useNsjail = nsjailAvailable()
	return r
}

func nsjailAvailable() bool {
	_, err := exec.LookPath("nsjail")
	return err == nil
}

func validateFlags(provided []string, allowlist []string) bool {
	for _, f := range provided {
		valid := false
		for _, a := range allowlist {
			if strings.HasSuffix(a, "*") {
				if strings.HasPrefix(f, strings.TrimSuffix(a, "*")) {
					valid = true
					break
				}
			} else {
				if f == a {
					valid = true
					break
				}
			}
		}
		if !valid {
			return false
		}
	}
	return true
}

func (r *Runner) Run(req *types.RunRequest) (*types.RunResponse, error) {
	lang, ok := Lookup(req.Language)
	if !ok {
		return nil, &types.APIError{Status: 400, Code: "unsupported_language",
			Message: fmt.Sprintf("unsupported language %q", req.Language),
		}
	}

	var bFlags []string
	if lang.Build != nil {
		if req.Build != nil && req.Build.Flags != nil {
			bFlags = req.Build.Flags
		}
		if !validateFlags(bFlags, lang.Build.FlagAllowlist) {
			return nil, &types.APIError{Status: 400, Code: "disallowed_flag", Message: "build flag not in allowlist"}
		}
	}

	var rFlags []string
	if lang.Run.Cmd != "" {
		if req.Run != nil && req.Run.Flags != nil {
			rFlags = req.Run.Flags
		}
		if !validateFlags(rFlags, lang.Run.FlagAllowlist) {
			return nil, &types.APIError{Status: 400, Code: "disallowed_flag", Message: "run flag not in allowlist"}
		}
	}

	workDir, err := os.MkdirTemp("", workDirPrefix)
	if err != nil {
		return nil, fmt.Errorf("create work dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	sourceFile := lang.SourceFilename
	if req.SourceFilename != "" {
		sourceFile = req.SourceFilename
	}
	if sourceFile == "" {
		sourceFile = "solution.txt"
	}
	sourcePath := filepath.Join(workDir, sourceFile)
	if err := os.WriteFile(sourcePath, []byte(req.Source), 0644); err != nil {
		return nil, fmt.Errorf("write source: %w", err)
	}

	artifactName := lang.Artifact
	if req.ArtifactFilename != "" {
		artifactName = req.ArtifactFilename
	}

	var buildResult types.BuildResult

	if lang.Build != nil {
		buildResult = r.doBuild(req, lang, workDir, sourceFile, artifactName, bFlags)
		if buildResult.Status != "ok" {
			testResults := make([]types.TestResult, len(req.Tests))
			for i := range req.Tests {
				testResults[i] = types.TestResult{Status: "not_executed"}
			}
			return &types.RunResponse{Status: "build_failed", Build: buildResult, Tests: testResults}, nil
		}
	} else {
		buildResult = types.BuildResult{Status: "ok"}
	}

	testResults := make([]types.TestResult, 0, len(req.Tests))
	overallStatus := "accepted"
	for _, tc := range req.Tests {
		tr := r.runTest(req, lang, workDir, sourceFile, artifactName, rFlags, tc)
		testResults = append(testResults, tr)
		if overallStatus == "accepted" && tr.Status != "accepted" {
			overallStatus = tr.Status
		}
	}

	return &types.RunResponse{
		Status: overallStatus,
		Build:  buildResult,
		Tests:  testResults,
	}, nil
}

func (r *Runner) doBuild(req *types.RunRequest, lang LanguageConfig, workDir, sourceFile, artifactName string, flags []string) types.BuildResult {
	buildLimits := defaultLimits()
	if lang.Build.Limits != nil {
		buildLimits = *lang.Build.Limits
	}
	if req.Build != nil && req.Build.Limits != nil {
		if req.Build.Limits.WallTimeS > 0 {
			buildLimits.WallTimeS = req.Build.Limits.WallTimeS
		}
		if req.Build.Limits.MemoryKB > 0 {
			buildLimits.MemoryKB = req.Build.Limits.MemoryKB
		}
		if req.Build.Limits.MaxProcesses > 0 {
			buildLimits.MaxProcesses = req.Build.Limits.MaxProcesses
		}
	}

	tpl := append([]string{lang.Build.Cmd}, lang.Build.Args...)
	cmdLine := expandTemplate(tpl, workDir, sourceFile, artifactName, flags)

	start := time.Now()
	stdout, stderr, runErr := r.execCmd(cmdLine, "", buildLimits, workDir)
	durMs := int(time.Since(start).Milliseconds())

	if runErr != nil {
		return types.BuildResult{Status: "failed", Stdout: stdout, Stderr: stderr, DurationMS: durMs}
	}
	return types.BuildResult{Status: "ok", Stdout: stdout, Stderr: stderr, DurationMS: durMs}
}

func (r *Runner) runTest(req *types.RunRequest, lang LanguageConfig, workDir, sourceFile, artifactName string, flags []string, tc types.TestCase) types.TestResult {
	runLimits := defaultLimits()
	if lang.Run.Limits != nil {
		runLimits = *lang.Run.Limits
	}
	if req.Run != nil && req.Run.Limits != nil {
		if req.Run.Limits.WallTimeS > 0 {
			runLimits.WallTimeS = req.Run.Limits.WallTimeS
		}
		if req.Run.Limits.MemoryKB > 0 {
			runLimits.MemoryKB = req.Run.Limits.MemoryKB
		}
		if req.Run.Limits.MaxProcesses > 0 {
			runLimits.MaxProcesses = req.Run.Limits.MaxProcesses
		}
	}

	tpl := append([]string{lang.Run.Cmd}, lang.Run.Args...)
	cmdLine := expandTemplate(tpl, workDir, sourceFile, artifactName, flags)

	start := time.Now()
	stdout, stderr, runErr := r.execCmd(cmdLine, tc.Stdin, runLimits, workDir)
	durMs := int(time.Since(start).Milliseconds())

	status := "accepted"
	if runErr != nil {
		if isTimeoutError(runErr) {
			status = "time_exceeded"
		} else {
			status = "runtime_error"
		}
	} else if tc.ExpectedStdout != "" {
		trimmedOut := strings.TrimRight(stdout, "\n")
		trimmedExpected := strings.TrimRight(tc.ExpectedStdout, "\n")
		if trimmedOut != trimmedExpected {
			status = "wrong_output"
		}
	} else if tc.ExpectedStdout == "" && stdout != "" {
		status = "wrong_output"
	}

	return types.TestResult{
		Status:       status,
		Stdout:       stdout,
		Stderr:       stderr,
		DurationMS:   durMs,
		MemoryPeakKB: 0, // Mocked for now, parsing cgroup stats is a bonus
	}
}

func (r *Runner) execCmd(cmdLine []string, stdin string, limits types.Limits, workDir string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(limits.WallTimeS)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if r.useNsjail {
		nsjailArgs := []string{
			"--really_quiet",
			"--disable_clone_newuser",
			"--rw",
			"--env", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"--bindmount_ro", "/usr",
			"--bindmount_ro", "/lib",
			"--bindmount_ro", "/lib64",
			"--bindmount_ro", "/bin",
			"--bindmount_ro", "/etc",
			"--bindmount", workDir,
		}
		if limits.WallTimeS > 0 {
			nsjailArgs = append(nsjailArgs, "--time_limit", fmt.Sprintf("%d", limits.WallTimeS))
		}
		if limits.MemoryKB > 0 {
			nsjailArgs = append(nsjailArgs, "--rlimit_as", fmt.Sprintf("%d", limits.MemoryKB/1024))
		}
		if limits.MaxProcesses > 0 {
			nsjailArgs = append(nsjailArgs, "--rlimit_nproc", fmt.Sprintf("%d", limits.MaxProcesses))
		}
		nsjailArgs = append(nsjailArgs, "--", cmdLine[0])
		nsjailArgs = append(nsjailArgs, cmdLine[1:]...)
		cmd = exec.CommandContext(ctx, "nsjail", nsjailArgs...)
	} else {
		cmd = exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	}

	stdoutBuf := &cappedWriter{limit: 1024 * 1024} // 1 MiB limit
	stderrBuf := &cappedWriter{limit: 1024 * 1024}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	err := cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

type cappedWriter struct {
	buf     bytes.Buffer
	limit   int
	written int
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	w.written += len(p)
	if w.buf.Len() < w.limit {
		writeLen := len(p)
		if w.buf.Len()+writeLen > w.limit {
			writeLen = w.limit - w.buf.Len()
		}
		w.buf.Write(p[:writeLen])
		if w.buf.Len() == w.limit {
			w.buf.WriteString("\n[output truncated]")
		}
	}
	return len(p), nil
}
func (w *cappedWriter) String() string {
	return w.buf.String()
}

func defaultLimits() types.Limits {
	return types.Limits{WallTimeS: 30, MemoryKB: 524288, MaxProcesses: 50}
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "signal: killed")
}

func expandTemplate(tpl []string, workDir, sourceFile, artifactName string, flags []string) []string {
	out := make([]string, 0, len(tpl)+len(flags))
	sourcePath := filepath.Join(workDir, sourceFile)
	var artifactPath string
	if artifactName != "" {
		artifactPath = filepath.Join(workDir, artifactName)
	}
	for _, s := range tpl {
		if s == "{{flags}}" {
			out = append(out, flags...)
		} else {
			s = strings.ReplaceAll(s, "{{source}}", sourcePath)
			if artifactPath != "" {
				s = strings.ReplaceAll(s, "{{artifact}}", artifactPath)
			}
			out = append(out, s)
		}
	}
	return out
}

func (r *Runner) Probe() types.Readiness {
	res := types.Readiness{
		Status:    "ok",
		Languages: make(map[string]types.ProbeResult),
	}

	// Probe nsjail
	_, err := exec.LookPath("nsjail")
	if err == nil {
		res.Nsjail = types.ProbeResult{OK: true, Version: "3.4"}
	} else {
		res.Nsjail = types.ProbeResult{OK: false, Error: "nsjail not found or not executable"}
		res.Status = "degraded"
	}

	for _, lang := range SupportedLanguages() {
		var cmd string
		if lang.Build != nil && lang.Build.Cmd != "" {
			cmd = lang.Build.Cmd
		} else if lang.Run.Cmd != "" {
			cmd = lang.Run.Cmd
		}
		
		if cmd == "" {
			res.Languages[lang.ID] = types.ProbeResult{OK: false, Error: "no cmd defined"}
			res.Status = "degraded"
			continue
		}

		out, err := exec.Command(cmd, "--version").CombinedOutput()
		if err != nil {
			// fallback to -version
			out, err = exec.Command(cmd, "-version").CombinedOutput()
		}

		if err == nil {
            lines := strings.Split(strings.TrimSpace(string(out)), "\n")
            ver := lines[0] // take first line
			res.Languages[lang.ID] = types.ProbeResult{OK: true, Version: ver}
		} else {
			res.Languages[lang.ID] = types.ProbeResult{OK: false, Error: err.Error() + ": " + string(out)}
			res.Status = "degraded"
		}
	}
	return res
}

func (r *Runner) Info() types.RunnerInfo {
	nsjailVer := "3.4"
	if _, err := exec.LookPath("nsjail"); err != nil {
		nsjailVer = "unknown"
	}
	
	langs := make([]any, 0, len(languageList))
	for _, l := range languageList {
		defLimits := defaultLimits()
		if l.Run.Limits != nil {
			defLimits = *l.Run.Limits
		}
		langs = append(langs, map[string]any{
			"id": l.ID,
			"name": l.Name,
			"default_run_limits": defLimits,
		})
	}

	return types.RunnerInfo{
		NsjailPath: "/usr/local/bin/nsjail",
		NsjailVersion: nsjailVer,
		Languages: langs,
	}
}
