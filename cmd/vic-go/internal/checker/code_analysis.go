package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckStatus represents the status of a check
type CheckStatus string

const (
	StatusPass    CheckStatus = "pass"
	StatusFail    CheckStatus = "fail"
	StatusSkip    CheckStatus = "skip"
	StatusUnknown CheckStatus = "unknown"
)

// CheckResult represents the result of a single alignment check
type CheckResult struct {
	RecordID string      `yaml:"record_id"`
	Category string      `yaml:"category"`
	Decision string      `yaml:"decision"`
	Status   CheckStatus `yaml:"status"`
	Message  string      `yaml:"message"`
	Details  interface{} `yaml:"details,omitempty"`
}

// TechPatterns contains regex patterns for technology detection
type TechPatterns struct {
	Category string
	Tech     string
	Patterns []string
}

// Common tech patterns for detection
var TechPatternList = []TechPatterns{
	// Database
	{"database", "postgresql", []string{"postgres", "psycopg", "pg\\.", "postgresql://", "Prisma.*postgresql"}},
	{"database", "mysql", []string{"mysql", "mysql2", "mariadb", "pymysql", "mysql://"}},
	{"database", "sqlite", []string{"sqlite", "\\.db", "sqlite3", "better-sqlite"}},
	{"database", "mongodb", []string{"mongodb", "mongoose", "MongoClient", "mongodb://"}},
	{"database", "redis", []string{"redis", "ioredis", "redis://"}},
	{"database", "prisma", []string{"prisma", "@prisma/client", "PrismaClient", "schema\\.prisma"}},
	{"database", "typeorm", []string{"typeorm", "TypeOrmModule"}},
	{"database", "sequelize", []string{"sequelize", "Sequelize"}},
	{"database", "sqlalchemy", []string{"sqlalchemy", "flask-sqlalchemy"}},

	// Auth
	{"auth", "jwt", []string{"jwt", "jsonwebtoken", " JJWT", "bcrypt"}},
	{"auth", "oauth", []string{"oauth", "passport", "Auth0", "auth0"}},
	{"auth", "casbin", []string{"casbin"}},
	{"auth", "session", []string{"express-session", "cookie-session", "session"}},

	// Frontend
	{"frontend", "react", []string{"react", "ReactDOM", "create-react-app", "Next\\.js"}},
	{"frontend", "vue", []string{"vue", "Vue\\.", "create-vue", "Nuxt"}},
	{"frontend", "angular", []string{"@angular", "ngModule", "NgModule"}},
	{"frontend", "svelte", []string{"svelte", "SvelteComponent"}},
	{"frontend", "redux", []string{"redux", "react-redux", "@reduxjs"}},
	{"frontend", "pinia", []string{"pinia", "createPinia"}},
	{"frontend", "zustand", []string{"zustand", "create\\(zustand"}},

	// Backend
	{"backend", "express", []string{"express", "Express\\.", "express\\."}},
	{"backend", "fastify", []string{"fastify", "Fastify"}},
	{"backend", "koa", []string{"koa", "Koa"}},
	{"backend", "nestjs", []string{"@nestjs", "NestFactory", "NestModule"}},
	{"backend", "django", []string{"django", "Django", "INSTALLED_APPS"}},
	{"backend", "flask", []string{"flask", "Flask", "@app.route"}},
	{"backend", "fastapi", []string{"fastapi", "FastAPI", "uvicorn"}},
	{"backend", "spring", []string{"@SpringBootApplication", "@RestController", "spring-boot"}},
	{"backend", "gin", []string{"gin-gonic", "gin\\.Engine"}},
	{"backend", "echo", []string{"echo\\.New", "labstack/echo"}},

	// Infrastructure
	{"infra", "docker", []string{"docker", "Dockerfile", "docker-compose"}},
	{"infra", "kubernetes", []string{"kubernetes", "k8s\\.", "\\.yaml"}},
	{"infra", "nginx", []string{"nginx", "nginx\\.conf"}},

	// Testing
	{"testing", "jest", []string{"jest", "describe\\(", "it\\("}},
	{"testing", "pytest", []string{"pytest", "def test_"}},
	{"testing", "mocha", []string{"mocha", "describe\\(", "it\\("}},
	{"testing", "rspec", []string{"rspec", "describe\\ ", "it\\ "}},
	{"testing", "junit", []string{"@Test", "junit"}},
}

// CodeAnalyzer analyzes source code for technology usage
type CodeAnalyzer struct {
	patterns    map[string]map[string][]*regexp.Regexp
	fileMatches map[string]map[string]bool
}

// NewCodeAnalyzer creates a new analyzer
func NewCodeAnalyzer() *CodeAnalyzer {
	patterns := make(map[string]map[string][]*regexp.Regexp)
	for _, tp := range TechPatternList {
		if patterns[tp.Category] == nil {
			patterns[tp.Category] = make(map[string][]*regexp.Regexp)
		}
		var regexes []*regexp.Regexp
		for _, p := range tp.Patterns {
			re, err := regexp.Compile(`(?i)` + p)
			if err == nil {
				regexes = append(regexes, re)
			}
		}
		patterns[tp.Category][tp.Tech] = regexes
	}

	return &CodeAnalyzer{
		patterns:    patterns,
		fileMatches: make(map[string]map[string]bool),
	}
}

// ScanDirectory scans a directory for source files
func (a *CodeAnalyzer) ScanDirectory(dir string) error {
	// Common source code directories
	sourceDirs := []string{"src", "lib", "app", "cmd", "internal", "pkg", "modules"}

	// Check if any source dir exists
	var scanDir string
	for _, s := range sourceDirs {
		testDir := filepath.Join(dir, s)
		if info, err := os.Stat(testDir); err == nil && info.IsDir() {
			scanDir = testDir
			break
		}
	}

	if scanDir == "" {
		scanDir = dir // Scan current directory
	}

	// Walk through directory
	return filepath.Walk(scanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip common non-source directories
		skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", "target", "__pycache__", ".venv", "venv"}
		for _, skip := range skipDirs {
			if strings.Contains(path, skip) {
				return nil
			}
		}

		// Check file extension
		ext := filepath.Ext(path)
		sourceExts := map[string]bool{
			".go": true, ".py": true, ".js": true, ".ts": true, ".tsx": true,
			".jsx": true, ".java": true, ".kt": true, ".cs": true, ".rb": true,
			".php": true, ".rs": true, ".vue": true, ".svelte": true,
		}

		if !sourceExts[ext] {
			return nil
		}

		// Read and analyze file
		a.analyzeFile(path)

		return nil
	})
}

// analyzeFile analyzes a single source file
func (a *CodeAnalyzer) analyzeFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	contentStr := string(content)
	a.fileMatches[path] = make(map[string]bool)

	// Check each pattern
	for category, techs := range a.patterns {
		for tech, regexes := range techs {
			for _, re := range regexes {
				if re.MatchString(contentStr) {
					key := category + ":" + tech
					a.fileMatches[path][key] = true
				}
			}
		}
	}
}

// GetDetectedTech returns all detected technologies
func (a *CodeAnalyzer) GetDetectedTech() map[string][]string {
	detected := make(map[string][]string)

	for _, matches := range a.fileMatches {
		for key := range matches {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				category, tech := parts[0], parts[1]
				found := false
				for _, t := range detected[category] {
					if t == tech {
						found = true
						break
					}
				}
				if !found {
					detected[category] = append(detected[category], tech)
				}
			}
		}
	}

	return detected
}

// CheckDecision checks if a decision aligns with detected code
func (a *CodeAnalyzer) CheckDecision(recordID, category, decision string) CheckResult {
	detected := a.GetDetectedTech()

	// Decision to technology mapping
	decisionMap := map[string][]string{
		"postgresql": {"postgresql", "pg", "postgres"},
		"mysql":      {"mysql", "mariadb"},
		"mongodb":    {"mongodb", "mongoose"},
		"sqlite":     {"sqlite"},
		"redis":      {"redis", "ioredis"},
		"prisma":     {"prisma"},
		"jwt":        {"jwt", "jsonwebtoken"},
		"oauth":      {"oauth", "passport", "auth0"},
		"react":      {"react", "next"},
		"vue":        {"vue", "nuxt"},
		"angular":    {"angular"},
		"express":    {"express"},
		"fastify":    {"fastify"},
		"django":     {"django"},
		"flask":      {"flask"},
		"fastapi":    {"fastapi"},
		"spring":     {"spring", "spring-boot"},
		"docker":     {"docker"},
		"kubernetes": {"kubernetes"},
	}

	// Normalize decision
	decisionLower := strings.ToLower(decision)

	// Find matching technologies
	for key, techs := range decisionMap {
		if strings.Contains(decisionLower, key) {
			// Check if any of these techs are detected
			detectedCategory := ""
			switch {
			case contains(techs, "postgresql", "mysql", "sqlite", "mongodb", "redis", "prisma"):
				detectedCategory = "database"
			case contains(techs, "jwt", "oauth", "oauth2"):
				detectedCategory = "auth"
			case contains(techs, "react", "vue", "angular", "svelte"):
				detectedCategory = "frontend"
			case contains(techs, "express", "fastify", "django", "flask", "spring"):
				detectedCategory = "backend"
			}

			if detectedCategory != "" {
				for _, tech := range techs {
					for _, det := range detected[detectedCategory] {
						if det == tech || strings.Contains(det, tech) {
							return CheckResult{
								RecordID: recordID,
								Category: category,
								Decision: decision,
								Status:   StatusPass,
								Message:  fmt.Sprintf("Detected: %s", det),
							}
						}
					}
				}
			}
		}
	}

	// No match found - could be planned or not implemented yet
	return CheckResult{
		RecordID: recordID,
		Category: category,
		Decision: decision,
		Status:   StatusUnknown,
		Message:  "Could not verify - no matching code patterns found",
	}
}

// contains checks if any of the items are in the slice
func contains(slice []string, items ...string) bool {
	for _, item := range items {
		for _, s := range slice {
			if s == item {
				return true
			}
		}
	}
	return false
}
