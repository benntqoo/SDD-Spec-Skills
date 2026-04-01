package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format represents the output format
type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Result represents a command result
type Result struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Errors    []Error     `json:"errors,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
	StartTime int64       `json:"start_time,omitempty"`
	EndTime   int64       `json:"end_time,omitempty"`
}

// Error represents a structured error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

// Printer handles output formatting
type Printer struct {
	format Format
	writer io.Writer
}

// NewPrinter creates a new printer
func NewPrinter(format Format, writer io.Writer) *Printer {
	return &Printer{format: format, writer: writer}
}

// Print outputs the result
func (p *Printer) Print(result *Result) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(result)
	case FormatYAML:
		return p.printYAML(result)
	default:
		return p.printPlain(result)
	}
}

func (p *Printer) printJSON(result *Result) error {
	encoder := json.NewEncoder(p.writer)
	encoder.SetIndent("", "  ") // Fixed: two parameters (prefix, indent)
	return encoder.Encode(result)
}

func (p *Printer) printYAML(result *Result) error {
	if result.Success {
		fmt.Fprintf(p.writer, "success: true\n")
		if result.Message != "" {
			fmt.Fprintf(p.writer, "message: %q\n", result.Message)
		}
	} else {
		fmt.Fprintf(p.writer, "success: false\n")
		for _, e := range result.Errors {
			fmt.Fprintf(p.writer, "error:\n  code: %s\n  message: %q\n", e.Code, e.Message)
			if e.Hint != "" {
				fmt.Fprintf(p.writer, "  hint: %q\n", e.Hint)
			}
		}
	}
	return nil
}

func (p *Printer) printPlain(result *Result) error {
	if result.Success {
		if result.Message != "" {
			fmt.Fprintf(p.writer, "✅ %s\n", result.Message)
		}
		for _, w := range result.Warnings {
			fmt.Fprintf(p.writer, "⚠️  %s\n", w)
		}
	} else {
		for _, e := range result.Errors {
			fmt.Fprintf(p.writer, "❌ [%s] %s\n", e.Code, e.Message)
			if e.Hint != "" {
				fmt.Fprintf(p.writer, "   💡 %s\n", e.Hint)
			}
		}
	}
	return nil
}

// Success creates a success result
func Success(message string, data ...interface{}) *Result {
	r := &Result{
		Success: true,
		Message: message,
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	return r
}

// Fail creates a failure result
func Fail(code, message string, hint ...string) *Result {
	r := &Result{
		Success: false,
		Errors:  []Error{{Code: code, Message: message}},
	}
	if len(hint) > 0 {
		r.Errors[0].Hint = hint[0]
	}
	return r
}

// WithWarning adds a warning to a result
func WithWarning(r *Result, warning string) *Result {
	r.Warnings = append(r.Warnings, warning)
	return r
}

// WithData adds data to a result
func WithData(r *Result, data interface{}) *Result {
	r.Data = data
	return r
}

// ParseFormat parses a format string
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatPlain
	}
}
