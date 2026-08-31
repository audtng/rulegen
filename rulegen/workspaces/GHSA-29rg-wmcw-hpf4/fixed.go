package main


import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Mzack9999/goja"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/nuclei/v3/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v3/pkg/protocols/common/protocolstate"
	"github.com/projectdiscovery/nuclei/v3/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestNewCompilerConsoleDebug(t *testing.T) {
	}
}

func TestRequireLocalFileAccessDenied(t *testing.T) {
	modulePath := writeModuleFile(t, t.TempDir(), "outside.js", `module.exports = { value: "outside-secret" };`)
	script := fmt.Sprintf(`var helper = require(%q); ExportAs("value", helper.value); true;`, modulePath)

	result, err := executeScript(t, t.Name(), false, script)
	require.Error(t, err)
	require.Contains(t, err.Error(), "-lfa is not enabled")
	require.Equal(t, err.Error(), result["error"])
}

func TestRequireTemplateModuleAllowedWithoutLFA(t *testing.T) {
	originalTemplatesDir := config.DefaultConfig.TemplatesDirectory
	templatesDir := t.TempDir()
	configuredTemplatesDir, moduleBaseDir := templateDirAlias(t, templatesDir)
	config.DefaultConfig.SetTemplatesDir(configuredTemplatesDir)
	t.Cleanup(func() {
		config.DefaultConfig.SetTemplatesDir(originalTemplatesDir)
	})

	modulePath := writeModuleFile(t, moduleBaseDir, filepath.Join("helpers", "allowed.js"), `module.exports = { value: "sandbox-ok" };`)
	script := fmt.Sprintf(`var helper = require(%q); ExportAs("value", helper.value); true;`, modulePath)

	result, err := executeScript(t, t.Name(), false, script)
	require.NoError(t, err)
	require.Equal(t, "sandbox-ok", result["value"])
}

func TestRequireLocalFileAccessAllowed(t *testing.T) {
	modulePath := writeModuleFile(t, t.TempDir(), "outside.js", `module.exports = { value: "outside-ok" };`)
	script := fmt.Sprintf(`var helper = require(%q); ExportAs("value", helper.value); true;`, modulePath)

	result, err := executeScript(t, t.Name(), true, script)
	require.NoError(t, err)
	require.Equal(t, "outside-ok", result["value"])
}

func TestRequireDoesNotReusePrivilegedModuleCacheAcrossExecutions(t *testing.T) {
	modulePath := writeModuleFile(t, t.TempDir(), "outside.js", `module.exports = { value: "outside-ok" };`)
	program, err := goja.Compile("", fmt.Sprintf(`require(%q).value`, modulePath), false)
	require.NoError(t, err)

	allowExecutionID := "allow-" + t.Name()
	denyExecutionID := "deny-" + t.Name()
	protocolstate.SetLfaAllowed(&types.Options{ExecutionId: allowExecutionID, AllowLocalFileAccess: true})
	protocolstate.SetLfaAllowed(&types.Options{ExecutionId: denyExecutionID, AllowLocalFileAccess: false})

	runtime := createNewRuntime()
	firstValue, err := executeWithRuntime(runtime, program, NewExecuteArgs(), &ExecuteOptions{
		ExecutionId: allowExecutionID,
		Context:     context.Background(),
	})
	require.NoError(t, err)
	require.Equal(t, "outside-ok", firstValue.Export())

	_, err = executeWithRuntime(runtime, program, NewExecuteArgs(), &ExecuteOptions{
		ExecutionId: denyExecutionID,
		Context:     context.Background(),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "-lfa is not enabled")
}

func executeScript(t *testing.T, executionID string, allowLocalFileAccess bool, script string) (ExecuteResult, error) {
	t.Helper()
	protocolstate.SetLfaAllowed(&types.Options{ExecutionId: executionID, AllowLocalFileAccess: allowLocalFileAccess})

	compiled, err := SourceAutoMode(script, false)
	require.NoError(t, err)

	compiler := New()
	return compiler.ExecuteWithOptions(compiled, NewExecuteArgs(), &ExecuteOptions{
		ExecutionId: executionID,
		Context:     context.Background(),
		Source:      &script,
		TimeoutVariants: &types.Timeouts{
			JsCompilerExecutionTimeout: 5 * time.Second,
		},
	})
}

func writeModuleFile(t *testing.T, baseDir string, relativePath string, contents string) string {
	t.Helper()
	modulePath := filepath.Join(baseDir, relativePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(modulePath), 0o755))
	require.NoError(t, os.WriteFile(modulePath, []byte(contents), 0o600))
	return modulePath
}

func templateDirAlias(t *testing.T, templateDir string) (string, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		return templateDir, templateDir
	}
	aliasPath := filepath.Join(t.TempDir(), "templates-link")
	if err := os.Symlink(templateDir, aliasPath); err != nil {
		return templateDir, templateDir
	}
	return aliasPath, templateDir
}

type noopWriter struct {
	Callback func(data []byte, level levels.Level)
}
)

var (
	lazyRegistryInit = sync.OnceFunc(func() {
		// autoregister console node module with default printer it uses gologger backend
		require.RegisterNativeModule(console.ModuleName, console.RequireWithPrinter(goconsole.NewGoConsolePrinter()))
	})
	for k, v := range args.Args {
		_ = runtime.Set(k, v)
	}

	runtime.SetContextValue("executionId", opts.ExecutionId)
	runtime.SetContextValue("ctx", opts.Context)
	enableRequire(runtime)

	// register extra callbacks if any
	if opts != nil && opts.Callback != nil {
		if err := opts.Callback(runtime); err != nil {
		}
	}

	// execute the script
	return runtime.RunProgram(p)
}
	return runtime
}

func enableRequire(runtime *goja.Runtime) {
	lazyRegistryInit()
	_ = require.NewRegistry(require.WithLoader(newSourceLoader(runtime))).Enable(runtime)
}

func newSourceLoader(runtime *goja.Runtime) require.SourceLoader {
	return func(path string) ([]byte, error) {
		executionID := ""
		if value, ok := runtime.GetContextValue("executionId"); ok {
			if id, ok := value.(string); ok {
				executionID = id
			}
		}

		normalizedPath, err := protocolstate.NormalizePathWithExecutionId(executionID, path)
		if err != nil {
			return nil, err
		}

		return require.DefaultSourceLoader(normalizedPath)
	}
}

func createNewRuntime() *goja.Runtime {
	runtime := protocolstate.NewJSRuntime()
	enableRequire(runtime)
	// by default import below modules every time
	_ = runtime.Set("console", require.Require(runtime, console.ModuleName))

package protocolstate

import (
	"github.com/projectdiscovery/nuclei/v3/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v3/pkg/types"
	filepathutil "github.com/projectdiscovery/nuclei/v3/pkg/utils/filepath"
	"github.com/projectdiscovery/utils/errkit"
	fileutil "github.com/projectdiscovery/utils/file"
	mapsutil "github.com/projectdiscovery/utils/maps"
	}
	// only allow files inside nuclei-templates directory
	// even current working directory is not allowed
	if filepathutil.IsPathWithinDirectory(cleaned, config.DefaultConfig.GetTemplateDir()) {
		return cleaned, nil
	}
	return "", errkit.Newf("path %v is outside nuclei-template directory and -lfa is not enabled", filePath)
	"github.com/projectdiscovery/nuclei/v3/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v3/pkg/model/types/severity"
	"github.com/projectdiscovery/nuclei/v3/pkg/templates/types"
	filepathutil "github.com/projectdiscovery/nuclei/v3/pkg/utils/filepath"
	"github.com/projectdiscovery/utils/errkit"
	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
	resolvedPath, err := fileutil.ResolveNClean(helperFilePath, config.DefaultConfig.GetTemplateDir())
	if err == nil {
		// As per rule 1, if helper file is present in nuclei-templates directory, allow it
		if filepathutil.IsPathWithinDirectory(resolvedPath, config.DefaultConfig.GetTemplateDir()) {
			return resolvedPath, nil
		}
	}
// Package filepathutil provides utilities for safe filepath operations,
// particularly for sandboxing file access in template environments.
//
// It includes functions to check if a path is contained within a directory,
// with proper canonicalization to handle symlinks and platform-specific
// path differences (such as case sensitivity on Windows).
//
// TODO(dwisiswant0): This package should be moved to the
// [github.com/projectdiscovery/utils/filepath], but let see how it goes first.
package filepathutil
package filepathutil

import (
	"path/filepath"
	"runtime"
	"strings"
)

// IsPathWithinDirectory returns true when path resolves inside directory.
// Both values are canonicalized to handle symlinks and platform-specific case rules.
func IsPathWithinDirectory(path string, directory string) bool {
	canonicalPath := canonicalizePath(path)
	canonicalDirectory := canonicalizePath(directory)

	relativePath, err := filepath.Rel(canonicalDirectory, canonicalPath)
	if err != nil {
		return false
	}
	return relativePath == "." || (relativePath != ".." && !strings.HasPrefix(relativePath, ".."+string(filepath.Separator)))
}

func canonicalizePath(path string) string {
	canonicalPath, err := filepath.Abs(path)
	if err != nil {
		canonicalPath = filepath.Clean(path)
	}
	if resolvedPath, err := filepath.EvalSymlinks(canonicalPath); err == nil {
		canonicalPath = resolvedPath
	}
	canonicalPath = filepath.Clean(canonicalPath)
	if runtime.GOOS == "windows" {
		canonicalPath = strings.ToLower(canonicalPath)
	}
	return canonicalPath
}
package filepathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsPathWithinDirectory(t *testing.T) {
	baseDir := t.TempDir()
	childFile := filepath.Join(baseDir, "nested", "child.txt")
	if err := os.MkdirAll(filepath.Dir(childFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(childFile, []byte("ok"), 0o600); err != nil {
		t.Fatal(err)
	}

	if !IsPathWithinDirectory(childFile, baseDir) {
		t.Fatalf("expected %q to be inside %q", childFile, baseDir)
	}

	outsideFile := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("nope"), 0o600); err != nil {
		t.Fatal(err)
	}
	if IsPathWithinDirectory(outsideFile, baseDir) {
		t.Fatalf("expected %q to be outside %q", outsideFile, baseDir)
	}
}

func TestIsPathWithinDirectoryWithSymlinkedDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation is not reliable on all Windows runners")
	}

	realDir := t.TempDir()
	aliasDir := filepath.Join(t.TempDir(), "templates-link")
	if err := os.Symlink(realDir, aliasDir); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	childFile := filepath.Join(realDir, "helpers", "allowed.js")
	if err := os.MkdirAll(filepath.Dir(childFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(childFile, []byte("module.exports = {};"), 0o600); err != nil {
		t.Fatal(err)
	}

	if !IsPathWithinDirectory(childFile, aliasDir) {
		t.Fatalf("expected %q to be inside symlinked dir %q", childFile, aliasDir)
	}
}
