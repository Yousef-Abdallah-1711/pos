package main

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

//go:embed routes.manifest
var routeManifest embed.FS

func TestRegisteredRouteManifestParity(t *testing.T) {
	manifestBytes, err := routeManifest.ReadFile("routes.manifest")
	if err != nil {
		t.Fatal(err)
	}
	expected := nonEmptyLines(string(manifestBytes))
	actual := registeredRoutePatterns(t)
	if len(actual) != len(expected) {
		t.Fatalf("registered route count changed: got %d, want %d", len(actual), len(expected))
	}
	for index := range expected {
		if actual[index] != expected[index] {
			t.Fatalf("registered route manifest differs at %d: got %q, want %q", index, actual[index], expected[index])
		}
	}
}

func registeredRoutePatterns(t *testing.T) []string {
	t.Helper()
	fileSet := token.NewFileSet()
	var patterns []string
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || (selector.Sel.Name != "Handle" && selector.Sel.Name != "HandleFunc") {
				return true
			}
			receiver, ok := selector.X.(*ast.Ident)
			if !ok || receiver.Name != "mux" {
				return true
			}
			literal, ok := call.Args[0].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			pattern, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Errorf("invalid route pattern %s: %v", literal.Value, err)
				return true
			}
			patterns = append(patterns, pattern)
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(patterns)
	return patterns
}

func nonEmptyLines(content string) []string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func TestRouteManifestCanBeReadFromWorkingDirectory(t *testing.T) {
	if _, err := os.Stat("routes.manifest"); err != nil {
		t.Fatalf("route manifest missing from package: %v", err)
	}
}
