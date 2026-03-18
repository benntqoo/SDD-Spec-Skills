package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for vic CLI
type Config struct {
	// Directory settings
	VICDir     string // Override .vic-sdd directory (default: .vic-sdd)
	ProjectDir string // Project directory (default: current directory)

	// File paths (computed)
	EventsFile          string
	StateFile           string
	TechRecordsFile     string
	RiskZonesFile       string
	ProjectFile         string
	DependencyGraphFile string
	SpecRequirements    string
	SpecArchitecture    string
	ProjectState        string

	// Output settings
	OutputFormat string // json, yaml, plain
	Verbose      bool

	// Git settings
	GitBranch string
}

// Load loads configuration from environment variables and config file
func Load() *Config {
	// Set up Viper
	viper.SetEnvPrefix("VIC")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("dir", ".vic-sdd")
	viper.SetDefault("output_format", "plain")
	viper.SetDefault("verbose", false)

	// Get environment variables
	vicDir := viper.GetString("dir")
	projectDir := viper.GetString("project_dir")
	if projectDir == "" {
		// Use current directory as default
		cwd, _ := os.Getwd()
		projectDir = cwd
	}

	// Build config
	cfg := &Config{
		VICDir:       vicDir,
		ProjectDir:   projectDir,
		OutputFormat: viper.GetString("output_format"),
		Verbose:      viper.GetBool("verbose"),
		GitBranch:    getGitBranch(projectDir),
	}

	// Compute file paths
	vicPath := filepath.Join(projectDir, vicDir)
	cfg.EventsFile = filepath.Join(vicPath, "status", "events.yaml")
	cfg.StateFile = filepath.Join(vicPath, "status", "state.yaml")
	cfg.TechRecordsFile = filepath.Join(vicPath, "tech", "tech-records.yaml")
	cfg.RiskZonesFile = filepath.Join(vicPath, "risk-zones.yaml")
	cfg.ProjectFile = filepath.Join(vicPath, "project.yaml")
	cfg.DependencyGraphFile = filepath.Join(vicPath, "dependency-graph.yaml")
	cfg.SpecRequirements = filepath.Join(vicPath, "SPEC-REQUIREMENTS.md")
	cfg.SpecArchitecture = filepath.Join(vicPath, "SPEC-ARCHITECTURE.md")
	cfg.ProjectState = filepath.Join(vicPath, "PROJECT.md")

	return cfg
}

// GetVICDir returns the VIC directory path
func (c *Config) GetVICDir() string {
	return filepath.Join(c.ProjectDir, c.VICDir)
}

// EnsureVICDir creates the VIC directory if it doesn't exist
func (c *Config) EnsureVICDir() error {
	vicDir := c.GetVICDir()
	if _, err := os.Stat(vicDir); os.IsNotExist(err) {
		return os.MkdirAll(vicDir, 0755)
	}
	return nil
}

// EnsureSubDir creates a subdirectory in VIC directory
func (c *Config) EnsureSubDir(name string) error {
	subDir := filepath.Join(c.GetVICDir(), name)
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		return os.MkdirAll(subDir, 0755)
	}
	return nil
}

func getGitBranch(projectDir string) string {
	// Simple git branch detection
	gitDir := filepath.Join(projectDir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return ""
	}

	// Try to read HEAD
	headFile := filepath.Join(gitDir, "HEAD")
	data, err := os.ReadFile(headFile)
	if err != nil {
		return ""
	}

	// Parse "ref: refs/heads/main" or "ref: refs/heads/master"
	var branch string
	_, err = fmt.Sscanf(string(data), "ref: refs/heads/%s", &branch)
	if err != nil {
		return "main" // default
	}
	return branch
}
