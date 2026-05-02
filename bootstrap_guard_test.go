package bootstrap_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var bootstrapPackageDirs = []string{
	// Keep this list for packages that MUST remain docs-only bootstrap stubs.
	// `app` is intentionally excluded because issue-18 approved its evolution
	// into a real implementation package.
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

func TestErrorsPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("errors")
	if err != nil {
		t.Fatalf("read errors dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("errors package must be allowed to include implementation files beyond doc.go, got %v", goFiles)
	}
}

func TestSecurityPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("security")
	if err != nil {
		t.Fatalf("read security dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("security package must be allowed to include implementation files beyond doc.go, got %v", goFiles)
	}
}

func TestHTTPGinPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("http/gin")
	if err != nil {
		t.Fatalf("read http/gin dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("http/gin package must include implementation files beyond doc.go, got %v", goFiles)
	}
}

func TestConfigPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("config")
	if err != nil {
		t.Fatalf("read config dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("config package must be allowed to include implementation files beyond doc.go, got %v", goFiles)
	}
}

func TestConfigViperPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("config/viper")
	if err != nil {
		t.Fatalf("read config/viper dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("config/viper package must include implementation files beyond doc.go, got %v", goFiles)
	}
}

func TestAppPackageCanEvolveBeyondBootstrapDocs(t *testing.T) {
	t.Helper()

	entries, err := os.ReadDir("app")
	if err != nil {
		t.Fatalf("read app dir: %v", err)
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

	if len(goFiles) <= 1 {
		t.Fatalf("app package must be allowed to include implementation files beyond doc.go, got %v", goFiles)
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

func TestDocsRuntimeLimitsContract(t *testing.T) {
	t.Helper()

	tests := []struct {
		name      string
		path      string
		fragments []string
	}{
		{
			name: "README includes runtime-limit keys defaults env and example",
			path: "README.md",
			fragments: []string{
				"read-timeout",
				"write-timeout",
				"max-header-bytes",
				"10s",
				"1048576",
				"COMMON_FWK_SERVER_READ_TIMEOUT",
				"COMMON_FWK_SERVER_WRITE_TIMEOUT",
				"COMMON_FWK_SERVER_MAX_HEADER_BYTES",
				"```yaml",
				"server:",
			},
		},
		{
			name: "architecture config-core includes runtime-limit keys defaults env and example",
			path: filepath.Join("docs", "architecture", "config-core.md"),
			fragments: []string{
				"read-timeout",
				"write-timeout",
				"max-header-bytes",
				"10s",
				"1048576",
				"COMMON_FWK_SERVER_READ_TIMEOUT",
				"COMMON_FWK_SERVER_WRITE_TIMEOUT",
				"COMMON_FWK_SERVER_MAX_HEADER_BYTES",
				"```yaml",
				"server:",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			content, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("read docs file %q: %v", tc.path, err)
			}

			doc := string(content)
			for _, fragment := range tc.fragments {
				if !strings.Contains(doc, fragment) {
					t.Fatalf("docs file %q missing required runtime-limit contract fragment %q", tc.path, fragment)
				}
			}
		})
	}
}

func TestDocsJWTModeReleaseAndMigrationContracts(t *testing.T) {
	t.Helper()

	tests := []struct {
		name      string
		path      string
		fragments []string
	}{
		{
			name: "architecture security-jwt includes HS256 and RS256 migration checkpoints",
			path: filepath.Join("docs", "architecture", "security-jwt.md"),
			fragments: []string{
				"Validate JWT mode-aware docs include HS256 default behavior and RS256 bootstrap keys",
				"Verify RS256 key-source coverage in tests (`generated`, `public-pem`, `private-pem`).",
				"Verify HS256 backward compatibility tests remain green.",
				"HS256 -> RS256 executable transition sequence",
				"security.auth.jwt.algorithm=RS256",
				"security.auth.jwt.rs256-key-source",
				"security.auth.jwt.rs256-key-id",
				"Protected routes still reject missing token with `401`.",
				"Valid RS256 token for configured issuer passes.",
				"HS256 path remains valid for services that have not migrated yet (default algorithm behavior).",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			content, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("read docs file %q: %v", tc.path, err)
			}

			doc := string(content)
			for _, fragment := range tc.fragments {
				if !strings.Contains(doc, fragment) {
					t.Fatalf("docs file %q missing required JWT mode contract fragment %q", tc.path, fragment)
				}
			}
		})
	}
}
