// Package swagger_merge_docs is a middleware plugin that serves inline content from a configuration.
// Paths are matched by patterns that are defined in the configuration.
package swagger_merge_docs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/usalko/swagger-merge-docs/docs"
)

// Config is the plugin configuration.
type Config struct {
	Path string     `json:"path"`
	Docs []*DocPath `json:"docs"`
}

// DocPath is a path configuration.
type DocPath struct {
	// Path is the exact path to match.
	Path string `json:"path"`
	// PathRegex is a regular expression to match.
	PathRegex string `json:"pathRegex"`
	// Content is a go template of content to serve.
	Content string `json:"content"`
	// JSONData is a map of JSON data to return.
	JSONData map[string]any `json:"jsonData"`
	// Indent is the number of spaces to indent the JSON response.
	Indent int `json:"indent"`
	// Status is the HTTP status code to return.
	Status int `json:"status"`

	pathRegex *regexp.Regexp
	template  *template.Template
	jsonData  []byte
}

// compile compiles the path and content.
func (path *DocPath) compile() error {
	if path.Path == "" && path.PathRegex == "" {
		return fmt.Errorf("path or pathRegex must be set")
	}
	if path.Content == "" && len(path.JSONData) == 0 {
		return fmt.Errorf("content or jsonData must be set")
	}
	var err error
	if path.PathRegex != "" {
		path.pathRegex, err = regexp.Compile(path.PathRegex)
		if err != nil {
			return fmt.Errorf("invalid path regexp: %w", err)
		}
	}
	if path.Content != "" {
		// Force a new line at the end of the template
		if !strings.HasSuffix(path.Content, "\n") {
			path.Content += "\n"
		}
		tmplname := path.Path
		if tmplname == "" {
			tmplname = path.PathRegex
		}
		path.template, err = template.New(tmplname).Parse(path.Content)
		if err != nil {
			return fmt.Errorf("invalid content template: %w", err)
		}
	}
	if len(path.JSONData) > 0 {
		if path.Indent == 0 {
			path.jsonData, err = json.Marshal(path.JSONData)
		} else {
			path.jsonData, err = json.MarshalIndent(path.JSONData, "", strings.Repeat(" ", path.Indent))
		}
		if err != nil {
			return fmt.Errorf("invalid json data: %w", err)
		}
		path.jsonData = append(path.jsonData, []byte("\n")...)
	}
	return err
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Docs: make([]*DocPath, 0, 10),
	}
}

// SwaggerMergeDocs is a plugin that merge multiply swagger docs into unified
type SwaggerMergeDocs struct {
	next          http.Handler
	path          string
	pathRegexp    *regexp.Regexp
	refs          []DocPath
	name          string
	staticContent []byte
}

// New creates a new StaticResponse plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	jsonConfig, _ := json.Marshal(config)
	log.Default().Printf("swagger-merge-docs configuration: %v", string(jsonConfig))

	if len(config.Docs) == 0 {
		return nil, fmt.Errorf("docs cannot be empty")
	}
	refs := make([]DocPath, len(config.Docs))
	for i, docPath := range config.Docs {
		ref := docPath
		// if err := ref.compile(); err != nil {
		// 	return nil, fmt.Errorf("invalid path configuration %s: %w", docPath.Path, err)
		// }
		refs[i] = *ref
	}
	pathRegexp, err := regexp.Compile(config.Path)
	if err != nil {
		log.Default().Printf("path is not regexp %v", err)
	}

	return &SwaggerMergeDocs{
		path:          config.Path,
		pathRegexp:    pathRegexp,
		refs:          refs,
		next:          next,
		name:          name,
		staticContent: docs.IndexHtml,
	}, nil
}

// ServeHTTP implements the http.Handler interface.
func (swaggerMerger *SwaggerMergeDocs) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := swaggerMerger.path
	pathRegex := swaggerMerger.pathRegexp
	if path != "" && (path == req.URL.Path) || (pathRegex != nil && pathRegex.MatchString(req.URL.Path)) {
		if len(swaggerMerger.staticContent) > 0 {
			rw.Header().Set("Content-Type", "text/html")
			fmt.Fprint(rw, string(swaggerMerger.staticContent))
			return
		}
		// if err := path.template.Execute(rw, map[string]any{
		// 	"Request": req,
		// }); err != nil {
		// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
		// }
		return
	}
	swaggerMerger.next.ServeHTTP(rw, req)
}
