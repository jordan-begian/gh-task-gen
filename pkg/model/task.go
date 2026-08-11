// Package model
package model

import (
	"fmt"
	"strings"
)

type Task struct {
	Title     string
	Body      *string
	Labels    []string
	Assignees []string
	Type      *string
}

func (t Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("Task title cannot be empty")
	}

	return nil
}
