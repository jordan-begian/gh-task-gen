// Package csv
package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"gh-task-gen/pkg/model"
)

type ParserFunc func(filePath string) ([]model.Task, error)

func NewParser() ParserFunc {
	return func(filePath string) ([]model.Task, error) {
		return parse(filePath)
	}
}

func parse(filePath string) ([]model.Task, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open csv file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv headers: %w", err)
	}

	colIndex := buildColumnIndex(headers)
	if _, ok := colIndex["title"]; !ok {
		return nil, fmt.Errorf("csv missing required 'title' column")
	}

	var tasks []model.Task
	lineNum := 1

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv line %d: %w", lineNum, err)
		}

		t, err := recordToTask(record, colIndex)
		if err != nil {
			return nil, fmt.Errorf("parse line %d: %w", lineNum, err)
		}

		if t.Title == "" {
			continue
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func buildColumnIndex(headers []string) map[string]int {
	index := make(map[string]int)
	for i, h := range headers {
		index[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return index
}

func recordToTask(record []string, colIndex map[string]int) (model.Task, error) {
	get := func(name string) string {
		if idx, ok := colIndex[name]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	getPtr := func(name string) *string {
		val := get(name)
		if val == "" {
			return nil
		}
		return &val
	}

	return model.Task{
		Title:     get("title"),
		Body:      getPtr("body"),
		Labels:    splitField(get("labels")),
		Assignees: splitField(get("assignees")),
		Type:      getPtr("type"),
	}, nil
}

func splitField(val string) []string {
	if val == "" {
		return nil
	}

	var result []string

	for v := range strings.SplitSeq(val, ";") {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
