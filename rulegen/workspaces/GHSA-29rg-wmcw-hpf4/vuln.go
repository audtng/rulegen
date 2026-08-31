package main


import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/nuclei/v3/pkg/types"
)

func TestNewCompilerConsoleDebug(t *testing.T) {
	}
}

type noopWriter struct {
	Callback func(data []byte, level levels.Level)
}
)

var (
	r                *require.Registry
	lazyRegistryInit = sync.OnceFunc(func() {
		r = new(require.Registry) // this can be shared by multiple runtimes
		// autoregister console node module with default printer it uses gologger backend
		require.RegisterNativeModule(console.ModuleName, console.RequireWithPrinter(goconsole.NewGoConsolePrinter()))
	})
	for k, v := range args.Args {
		_ = runtime.Set(k, v)
	}
	// register extra callbacks if any
	if opts != nil && opts.Callback != nil {
		if err := opts.Callback(runtime); err != nil {
		}
	}

	// inject execution id and context
	runtime.SetContextValue("executionId", opts.ExecutionId)
	runtime.SetContextValue("ctx", opts.Context)

	// execute the script
	return runtime.RunProgram(p)
}
	return runtime
}

func getRegistry() *require.Registry {
	lazyRegistryInit()
	return r
}

func createNewRuntime() *goja.Runtime {
	runtime := protocolstate.NewJSRuntime()
	_ = getRegistry().Enable(runtime)
	// by default import below modules every time
	_ = runtime.Set("console", require.Require(runtime, console.ModuleName))

package protocolstate

import (
	"strings"

	"github.com/projectdiscovery/nuclei/v3/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v3/pkg/types"
	"github.com/projectdiscovery/utils/errkit"
	fileutil "github.com/projectdiscovery/utils/file"
	mapsutil "github.com/projectdiscovery/utils/maps"
	}
	// only allow files inside nuclei-templates directory
	// even current working directory is not allowed
	if strings.HasPrefix(cleaned, config.DefaultConfig.GetTemplateDir()) {
		return cleaned, nil
	}
	return "", errkit.Newf("path %v is outside nuclei-template directory and -lfa is not enabled", filePath)
	"github.com/projectdiscovery/nuclei/v3/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v3/pkg/model/types/severity"
	"github.com/projectdiscovery/nuclei/v3/pkg/templates/types"
	"github.com/projectdiscovery/utils/errkit"
	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
	resolvedPath, err := fileutil.ResolveNClean(helperFilePath, config.DefaultConfig.GetTemplateDir())
	if err == nil {
		// As per rule 1, if helper file is present in nuclei-templates directory, allow it
		if strings.HasPrefix(resolvedPath, config.DefaultConfig.GetTemplateDir()) {
			return resolvedPath, nil
		}
	}
