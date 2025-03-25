// Package static_response_plugin is a middleware plugin that serves inline content from a configuration.
// Paths are matched by patterns that are defined in the configuration.
package static_response_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

// Config is the plugin configuration.
type Config struct {
	Path        Path   `json:"path"`
	SwaggerRefs []Path `json:"refs"`
}

// Path is a path configuration.
type Path struct {
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
func (p *Path) compile() error {
	if p.Path == "" && p.PathRegex == "" {
		return fmt.Errorf("path or pathRegex must be set")
	}
	if p.Content == "" && len(p.JSONData) == 0 {
		return fmt.Errorf("content or jsonData must be set")
	}
	var err error
	if p.PathRegex != "" {
		p.pathRegex, err = regexp.Compile(p.PathRegex)
		if err != nil {
			return fmt.Errorf("invalid path regexp: %w", err)
		}
	}
	if p.Content != "" {
		// Force a new line at the end of the template
		if !strings.HasSuffix(p.Content, "\n") {
			p.Content += "\n"
		}
		tmplname := p.Path
		if tmplname == "" {
			tmplname = p.PathRegex
		}
		p.template, err = template.New(tmplname).Parse(p.Content)
		if err != nil {
			return fmt.Errorf("invalid content template: %w", err)
		}
	}
	if len(p.JSONData) > 0 {
		if p.Indent == 0 {
			p.jsonData, err = json.Marshal(p.JSONData)
		} else {
			p.jsonData, err = json.MarshalIndent(p.JSONData, "", strings.Repeat(" ", p.Indent))
		}
		if err != nil {
			return fmt.Errorf("invalid json data: %w", err)
		}
		p.jsonData = append(p.jsonData, []byte("\n")...)
	}
	return err
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		SwaggerRefs: make([]Path, 0),
	}
}

// SwaggerMergeDocs is a plugin that merge multiply swagger docs into unified
type SwaggerMergeDocs struct {
	next http.Handler
	path Path
	refs []Path
	name string
}

// New creates a new StaticResponse plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.SwaggerRefs) == 0 {
		return nil, fmt.Errorf("paths cannot be empty")
	}
	refs := make([]Path, len(config.SwaggerRefs))
	for i, p := range config.SwaggerRefs {
		ref := &p
		if err := ref.compile(); err != nil {
			return nil, fmt.Errorf("invalid path configuration %s: %w", p.Path, err)
		}
		refs[i] = *ref
	}
	return &SwaggerMergeDocs{
		path: config.Path,
		refs: refs,
		next: next,
		name: name,
	}, nil
}

// ServeHTTP implements the http.Handler interface.
func (swaggerMerger *SwaggerMergeDocs) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := swaggerMerger.path
	if (path.Path != "" && path.Path == req.URL.Path) || (path.pathRegex != nil && path.pathRegex.MatchString(req.URL.Path)) {
		if path.Status != 0 {
			rw.WriteHeader(path.Status)
		}
		if len(path.jsonData) > 0 {
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, string(path.jsonData))
			return
		}
		if err := path.template.Execute(rw, map[string]any{
			"Request": req,
		}); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	swaggerMerger.next.ServeHTTP(rw, req)
}
