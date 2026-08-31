package main

	"path/filepath"
	"regexp"
	"strings"
	"rogchap.com/v8go"
)

	// Create custom variable name for component based on the file path for the layout.
	componentSignature := strings.ReplaceAll(strings.ReplaceAll(layoutPath, "/", "_"), ".", "_")
	// Use signature instead of specific component name (e.g. var Html = create_ssr_component(($$result, $$props, $$bindings, slots) => {)
  ssrStr = reSSRComp.ReplaceAllString(ssrStr, "${1}"+componentSignature+"${2}")

	namedExports := reStaticExport.FindAllStringSubmatch(ssrStr, -1)
	// Loop through all export statements.
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
  "regexp"

	"github.com/plentico/plenti/cmd/build"
	"github.com/plentico/plenti/cmd/serve"
	"github.com/briandowns/spinner"
	"github.com/gerald1248/httpscerts"
	"github.com/spf13/cobra"
  "github.com/go-playground/validator/v10"
	"golang.org/x/net/websocket"
)

// LocalFlag can be set to false to emulate a remote environment
var LocalFlag bool

// Valditor for input validation
var validate *validator.Validate

func checkPortAvailability(port int) bool {
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", address)

		// Start the HTTP webserver
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
    
	},
}

	serveCmd.Flags().StringVarP(&ConfigFileFlag, "config", "c", "plenti.json", "use a custom sitewide configuration file")
}

// Validate user supplied values
type localChange struct {
  Action   string `json:"action" validate:"required,oneof=create update delete"`
  Encoding string `json:"encoding" validate:"required,oneof=base64 text"`
  File     string `json:"file" validate:"file-path"`
  Contents string `json:"contents" validate:"required"`
}

// Custom validation for file path. Only allow files in the layouts and content directories.
func FilePathValidation(fl validator.FieldLevel) bool {
  reFilePath := regexp.MustCompile(`^(layouts|content)[a-zA-Z0-9_\-\/]*(.svelte|.js|json)$`)
  fmt.Println(fl.Field().String())
  return reFilePath.MatchString(fl.Field().String())
}

func postLocal(w http.ResponseWriter, r *http.Request) {
  // Register custom rules to validator
  validate = validator.New()
  validate.RegisterValidation("file-path", FilePathValidation)
	
  if r.Method == "POST" {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Could not read 'body' from local edit: %v", err)
		}
		if err != nil {
			fmt.Printf("Could not unmarshal JSON data: %v", err)
		}

		var contents []byte
		for _, change := range localChanges {
			      
      // Validate user input, there is any error, return 400 Bad Request
      err := validate.Struct(change)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

			if change.Action == "create" || change.Action == "update" {
				contents = []byte(change.Contents)
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		ErrorLog:       log.New(io.Discard, "", 0),
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      cfg,
	}

// Custom validation for file path. Only allow files in the layouts and content directories.
func FilePathValidation(fl validator.FieldLevel) bool {
  reFilePath := regexp.MustCompile(`^(layouts|content)[a-zA-Z0-9_\-\/]*(.svelte|.js|.json)$`)
  fmt.Println(fl.Field().String())
  return reFilePath.MatchString(fl.Field().String())
}
fixed, only json can be upload under content file.
cmd/build/compile.go                          |   3 +-
cmd/serve.go                                  |  44 ++++++------
test-site/.gitignore                          |   2 -
test-site/.gitlab-ci.yml                      |  12 ----
test-site/content/_index.json                 |  13 ----
test-site/content/blog/_defaults.json         |   6 --
test-site/content/blog/_schema.json           |   8 ---
test-site/content/blog/components.json        |  18 -----
test-site/content/blog/perry.json             |  10 ---
test-site/content/blog/pletiform.json         |  10 ---
test-site/content/blog/stores.json            |  11 ---
test-site/content/pages/_defaults.json        |   5 --
test-site/content/pages/about.json            |  16 -----
test-site/content/pages/admin.json            |  13 ----
test-site/content/pages/contact.json          |  12 ----
test-site/layouts/components/ball.svelte      |  49 -------------
test-site/layouts/components/block.svelte     |  23 -------
.../layouts/components/decrementer.svelte     |  11 ---
test-site/layouts/components/grid.svelte      |  37 ----------
.../layouts/components/incrementer.svelte     |  11 ---
test-site/layouts/components/pager.svelte     |  54 ---------------
test-site/layouts/components/source.svelte    |  65 ------------------
test-site/layouts/content/404.svelte          |   2 -
test-site/layouts/content/_index.svelte       |  35 ----------
test-site/layouts/content/blog.svelte         |  40 -----------
test-site/layouts/content/pages.svelte        |  28 --------
test-site/layouts/content/vartest.svelte      |   1 -
test-site/layouts/global/footer.svelte        |  33 ---------
test-site/layouts/global/head.svelte          |  18 -----
test-site/layouts/global/html.svelte          |  31 ---------
test-site/layouts/global/nav.svelte           |  30 --------
test-site/layouts/global/vartest.svelte       |   1 -
test-site/layouts/scripts/make_title.svelte   |  15 ----
test-site/layouts/scripts/sort_by_date.svelte |  11 ---
test-site/layouts/scripts/stores.svelte       |   4 --
test-site/media/perry.webp                    | Bin 17912 -> 0 bytes
test-site/package.json                        |  11 ---
test-site/plenti.json                         |  20 ------
test-site/static/global.css                   |  47 -------------
test-site/static/logo.svg                     |  17 -----
test-site/static/robots.txt                   |   7 --
41 files changed, 24 insertions(+), 760 deletions(-)
delete mode 100755 test-site/.gitignore
delete mode 100755 test-site/.gitlab-ci.yml
delete mode 100755 test-site/content/_index.json
delete mode 100755 test-site/content/blog/_defaults.json
delete mode 100755 test-site/content/blog/_schema.json
delete mode 100755 test-site/content/blog/components.json
delete mode 100755 test-site/content/blog/perry.json
delete mode 100755 test-site/content/blog/pletiform.json
delete mode 100755 test-site/content/blog/stores.json
delete mode 100755 test-site/content/pages/_defaults.json
delete mode 100755 test-site/content/pages/about.json
delete mode 100755 test-site/content/pages/admin.json
delete mode 100755 test-site/content/pages/contact.json
delete mode 100755 test-site/layouts/components/ball.svelte
delete mode 100755 test-site/layouts/components/block.svelte
delete mode 100755 test-site/layouts/components/decrementer.svelte
delete mode 100755 test-site/layouts/components/grid.svelte
delete mode 100755 test-site/layouts/components/incrementer.svelte
delete mode 100755 test-site/layouts/components/pager.svelte
delete mode 100755 test-site/layouts/components/source.svelte
delete mode 100755 test-site/layouts/content/404.svelte
delete mode 100755 test-site/layouts/content/_index.svelte
delete mode 100755 test-site/layouts/content/blog.svelte
delete mode 100755 test-site/layouts/content/pages.svelte
delete mode 100755 test-site/layouts/content/vartest.svelte
delete mode 100755 test-site/layouts/global/footer.svelte
delete mode 100755 test-site/layouts/global/head.svelte
delete mode 100755 test-site/layouts/global/html.svelte
delete mode 100755 test-site/layouts/global/nav.svelte
delete mode 100755 test-site/layouts/global/vartest.svelte
delete mode 100755 test-site/layouts/scripts/make_title.svelte
delete mode 100755 test-site/layouts/scripts/sort_by_date.svelte
delete mode 100755 test-site/layouts/scripts/stores.svelte
delete mode 100755 test-site/media/perry.webp
delete mode 100755 test-site/package.json
delete mode 100755 test-site/plenti.json
delete mode 100755 test-site/static/global.css
delete mode 100755 test-site/static/logo.svg
delete mode 100755 test-site/static/robots.txt
	"path/filepath"
	"regexp"
	"strings"

	"rogchap.com/v8go"
)

	// Create custom variable name for component based on the file path for the layout.
	componentSignature := strings.ReplaceAll(strings.ReplaceAll(layoutPath, "/", "_"), ".", "_")
	// Use signature instead of specific component name (e.g. var Html = create_ssr_component(($$result, $$props, $$bindings, slots) => {)
	ssrStr = reSSRComp.ReplaceAllString(ssrStr, "${1}"+componentSignature+"${2}")

	namedExports := reStaticExport.FindAllStringSubmatch(ssrStr, -1)
	// Loop through all export statements.
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/plentico/plenti/cmd/build"
	"github.com/plentico/plenti/cmd/serve"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/briandowns/spinner"
	"github.com/gerald1248/httpscerts"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)


		// Start the HTTP webserver
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	},
}


// Validate user supplied values
type localChange struct {
	Action   string `json:"action" validate:"required,oneof=create update delete"`
	Encoding string `json:"encoding" validate:"required,oneof=base64 text"`
	File     string `json:"file" validate:"file-path"`
	Contents string `json:"contents" validate:"required"`
}

// Custom validation for file path. Only allow files in the layouts and content directories.
func FilePathValidation(fl validator.FieldLevel) bool {
	reFilePath := regexp.MustCompile(`^(layouts|content)[a-zA-Z0-9_\-\/]*(.svelte|.js|.json)$`)
	fmt.Println(fl.Field().String())
	return reFilePath.MatchString(fl.Field().String())
}

func postLocal(w http.ResponseWriter, r *http.Request) {
	// Register custom rules to validator
	validate = validator.New()
	validate.RegisterValidation("file-path", FilePathValidation)

	if r.Method == "POST" {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Could not read 'body' from local edit: %v", err)

		var contents []byte
		for _, change := range localChanges {

			// Validate user input, there is any error, return 400 Bad Request
			err := validate.Struct(change)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if change.Action == "create" || change.Action == "update" {
				contents = []byte(change.Contents)

// Custom validation for file path. Only allow files in the layouts and content directories.
func FilePathValidation(fl validator.FieldLevel) bool {
	reFilePath := regexp.MustCompile(`^(content)[a-zA-Z0-9_\-\/]*(.json)$`)
	fmt.Println(fl.Field().String())
	return reFilePath.MatchString(fl.Field().String())
}
