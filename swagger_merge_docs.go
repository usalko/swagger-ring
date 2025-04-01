// Package swagger_merge_docs is a middleware plugin that serves inline content from a configuration.
// Paths are matched by patterns that are defined in the configuration.
package swagger_merge_docs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/usalko/swagger-merge-docs/docs"
	"gopkg.in/yaml.v3"
)

// Config is the plugin configuration.
type Config struct {
	Path string     `json:"path"`
	Docs []*DocPath `json:"docs"`
}

const (
	DOC_TYPE_YAML = iota
	DOC_TYPE_JSON
)

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
	// jsonConfig, _ := json.Marshal(config)
	// log.Default().Printf("⭕swagger-merge-docs configuration: %v", string(jsonConfig))

	if len(config.Docs) == 0 {
		return nil, fmt.Errorf("⭕docs cannot be empty")
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
		log.Default().Printf("⭕path is not regexp %v", err)
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

func (swaggerMerger *SwaggerMergeDocs) GetMergedSwaggerDoc(docType int) (string, error) {
	// log.Default().Printf("⭕refs are %v", swaggerMerger.refs)
	result := make(map[interface{}]interface{}, 0)
	for _, ref := range swaggerMerger.refs {
		// Get the data
		resp, err := http.Get(ref.Path)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		buf := bytes.NewBufferString("")
		// Writer the body to file
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(ref.Path, ".yml") || strings.HasSuffix(ref.Path, ".yaml") {
			var swagger map[interface{}]interface{}
			err = yaml.Unmarshal(buf.Bytes(), &swagger)
			if err != nil {
				return "", err
			}
			swaggerMerger.deepMergeDocs(result, swagger)
			continue
		}
		if strings.HasSuffix(ref.Path, ".json") {
			var swagger map[interface{}]interface{}
			err = json.Unmarshal(buf.Bytes(), &swagger)
			if err != nil {
				return "", err
			}
			swaggerMerger.deepMergeDocs(result, swagger)
			continue
		}
	}

	if docType == DOC_TYPE_YAML {
		mergedDoc, err := yaml.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(mergedDoc), nil
	}
	if docType == DOC_TYPE_JSON {
		mergedDoc, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(mergedDoc), nil
	}
	return "", fmt.Errorf("unknown document type %v", docType)
}

// deepMergeDocs рекурсивно объединяет два YAML/JSON-объекта
func (swaggerMerger *SwaggerMergeDocs) deepMergeDocs(dst, src map[interface{}]interface{}) {
	for key, srcVal := range src {
		// Если ключ уже есть в dst
		if dstVal, exists := dst[key]; exists {
			// log.Default().Printf("🔥 dstVal for %v is %T", key, dstVal)
			// log.Default().Printf("🔥 srcVal for %v is %T", key, srcVal)
			// Если оба значения — map, рекурсивно объединяем
			if dstMap, ok := dstVal.(map[interface{}]interface{}); ok {
				if srcMap, ok := srcVal.(map[interface{}]interface{}); ok {
					swaggerMerger.deepMergeDocs(dstMap, srcMap)
					dst[key] = dstMap
					continue
				}
			}
			// Если оба значения — slice, объединяем
			if dstSlice, ok := dstVal.([]interface{}); ok {
				if srcSlice, ok := srcVal.([]interface{}); ok {
					slicesUnion := append([]interface{}{}, dstSlice...)
					slicesUnion = append(slicesUnion, srcSlice)
					dst[key] = slicesUnion
					continue
				}
			}
		}
		// Иначе просто перезаписываем
		dst[key] = srcVal
	}
}

// ServeHTTP implements the http.Handler interface.
func (swaggerMerger *SwaggerMergeDocs) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := swaggerMerger.path

	// log.Default().Printf("⭕request.path is %v, path is %v", req.URL.Path, path)

	if path != "" && (path == req.URL.Path) {
		if len(swaggerMerger.staticContent) > 0 {
			rw.Header().Set("Content-Type", "text/html")
			fmt.Fprint(rw, string(swaggerMerger.staticContent))
		}
		return
	}
	if path != "" && (path+"/doc.yaml" == req.URL.Path) {
		rw.Header().Set("Content-Type", "application/yaml")
		mergedSwaggerDocument, err := swaggerMerger.GetMergedSwaggerDoc(DOC_TYPE_YAML)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(rw, mergedSwaggerDocument)
		return
	}
	if path != "" && (path+"/doc.json" == req.URL.Path) {
		rw.Header().Set("Content-Type", "application/yaml")
		mergedSwaggerDocument, err := swaggerMerger.GetMergedSwaggerDoc(DOC_TYPE_JSON)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(rw, mergedSwaggerDocument)
		return
	}
	swaggerMerger.next.ServeHTTP(rw, req)
}
