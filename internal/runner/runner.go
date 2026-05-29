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

func (r *Runner) Run(req *types.RunRequest) (*types.RunResponse, error) {
	lang, ok := Lookup(req.Language)
	if !ok {
		return nil, &types.APIError{Status: 400, Code: "unsupported_language",
			Message: fmt.Sprintf("unsupported language %q", req.Language),
		}
	}

	workDir, err := os.MkdirTemp("", workDirPrefix)
	if err != nil {
		return nil, fmt.Errorf("create work dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	sourceFile := "main" + lang.SourceExt
	if req.SourceFilename != "" {
		sourceFile = req.SourceFilename
	}
	sourcePath := filepath.Join(workDir, sourceFile)
	if err := os.WriteFile(sourcePath, []byte(req.Source), 0644); err != nil {
		return nil, fmt.Errorf("write source: %w", err)
	}

	artifactName := lang.BuildArtifact
	var buildResult types.BuildResult

	if lang.BuildCmd != nil {
		buildResult = r.doBuild(req, lang, workDir, sourceFile, artifactName)
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
		tr := r.runTest(req, lang, workDir, sourceFile, artifactName, tc)
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

func (r *Runner) doBuild(req *types.RunRequest, lang LanguageDef, workDir, sourceFile, artifactName string) types.BuildResult {
	buildLimits := defaultLimits()
	if req.Build != nil && req.Build.Limits != nil {
		buildLimits = *req.Build.Limits
	}

	cmdLine := expandTemplate(lang.BuildCmd, workDir, sourceFile, workDir, artifactName)

	start := time.Now()
	stdout, stderr, runErr := r.execCmd(cmdLine, "", buildLimits, workDir)
	durMs := int(time.Since(start).Milliseconds())

	if runErr != nil {
		return types.BuildResult{Status: "failed", Stdout: stdout, Stderr: stderr, DurationMS: durMs}
	}
	return types.BuildResult{Status: "ok", Stdout: stdout, Stderr: stderr, DurationMS: durMs}
}

func (r *Runner) runTest(req *types.RunRequest, lang LanguageDef, workDir, sourceFile, artifactName string, tc types.TestCase) types.TestResult {
	runLimits := defaultLimits()
	if req.Run != nil && req.Run.Limits != nil {
		runLimits = *req.Run.Limits
	}

	var cmdLine []string
	if lang.RunNeedsArtifact {
		cmdLine = expandTemplate(lang.RunCmd, workDir, sourceFile, workDir, artifactName)
	} else {
		cmdLine = expandTemplate(lang.RunCmd, workDir, sourceFile, workDir, "")
	}

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
		MemoryPeakKB: 0,
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

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	err := cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
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

func expandTemplate(tpl []string, sourceDir, sourceFile, artifactDir, artifactName string) []string {
	out := make([]string, len(tpl))
	for i, s := range tpl {
		s = strings.ReplaceAll(s, "{source_dir}", sourceDir)
		s = strings.ReplaceAll(s, "{source_file}", sourceFile)
		s = strings.ReplaceAll(s, "{artifact_dir}", artifactDir)
		s = strings.ReplaceAll(s, "{artifact_name}", artifactName)
		out[i] = s
	}
	return out
}
