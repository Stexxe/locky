//+build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"log"
	"strings"
	"bytes"
	"go/format"
)

const (
	output_file = "data.go"
	res_dir = "../../res"
)

func main() {
	var conv = map[string]interface{} {"conv": fmtByteSlice}
	tmpl := template.Must(template.New("").Funcs(conv).Parse(`package resources

// Code generated by go generate; DO NOT EDIT.

var resources = make(map[string]*[]byte)

func init() {
    {{- range $name, $file := . }}
		resources["{{ $name }}"] = &[]byte{ {{ conv $file }} }
    {{- end }}
}`),
	)

	result := make(map[string][]byte)

	filepath.Walk(res_dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(res_dir, path)
		bytes, err := ioutil.ReadFile(path)

		if err != nil {
			return nil
		}

		result[relPath] = bytes

		return nil
	})

	builder := &bytes.Buffer{}
	if err := tmpl.Execute(builder, result); err != nil {
		log.Fatal("Error executing template", err)
	}

	data, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("Error formatting generated code", err)
	}

	if err = ioutil.WriteFile(output_file, data, os.ModePerm); err != nil {
		log.Fatal("Error writing blob file", err)
	}
}

func fmtByteSlice(b []byte) string {
	builder := strings.Builder{}

	for _, v := range b {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}

	return builder.String()
}
