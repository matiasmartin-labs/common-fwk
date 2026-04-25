package bootstrap_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var bootstrapPackageDirs = []string{
	"app",
	"config",
	"config/viper",
	"security",
	"http/gin",
	"errors",
}

func TestBootstrapPackagesRemainStructuralOnly(t *testing.T) {
	t.Helper()

	for _, dir := range bootstrapPackageDirs {
		dir := dir
		t.Run(dir, func(t *testing.T) {
			t.Helper()

			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatalf("read bootstrap dir %q: %v", dir, err)
			}

			goFiles := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				if strings.HasSuffix(entry.Name(), ".go") {
					goFiles = append(goFiles, entry.Name())
				}
			}

			if len(goFiles) != 1 || goFiles[0] != "doc.go" {
				t.Fatalf("bootstrap package %q must contain only doc.go, got %v", dir, goFiles)
			}

			docPath := filepath.Join(dir, "doc.go")
			content, err := os.ReadFile(docPath)
			if err != nil {
				t.Fatalf("read %q: %v", docPath, err)
			}

			doc := string(content)
			if strings.Contains(doc, "func ") {
				t.Fatalf("business behavior detected in %q: functions are not allowed during bootstrap", docPath)
			}
		})
	}
}

func TestCIBaselineIncludesPRTriggerAndGoTestCommand(t *testing.T) {
	t.Helper()

	workflowPath := filepath.Join(".github", "workflows", "ci.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow %q: %v", workflowPath, err)
	}

	workflow := string(content)

	mustContain := []string{
		"push:",
		"pull_request:",
		"run: go test ./...",
	}

	for _, fragment := range mustContain {
		if !strings.Contains(workflow, fragment) {
			t.Fatalf("workflow %q missing required fragment %q", workflowPath, fragment)
		}
	}

	if strings.Contains(workflow, "continue-on-error: true") {
		t.Fatalf("workflow %q disables fail-on-error semantics via continue-on-error", workflowPath)
	}

	if strings.Contains(workflow, "run: go test ./... || true") {
		t.Fatalf("workflow %q disables fail-on-error semantics via shell bypass", workflowPath)
	}
}

func TestCIBaselineRemainsBootstrapMinimal(t *testing.T) {
	t.Helper()

	workflowPath := filepath.Join(".github", "workflows", "ci.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow %q: %v", workflowPath, err)
	}

	workflow := string(content)

	if strings.Count(workflow, "run: go test ./...") != 1 {
		t.Fatalf("workflow %q must keep a single baseline go test command", workflowPath)
	}

	for _, forbidden := range []string{"golangci-lint", "coverage", "release"} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("workflow %q includes out-of-scope gate %q", workflowPath, forbidden)
		}
	}
}
