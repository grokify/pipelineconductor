// Package collector provides interfaces and implementations for collecting
// repository and workflow data from various sources.
package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

// LocalCollector collects repository and workflow data from the local filesystem.
type LocalCollector struct {
	BasePath string
	Workers  int
	Verbose  bool
}

// NewLocalCollector creates a new LocalCollector.
func NewLocalCollector(basePath string) *LocalCollector {
	return &LocalCollector{
		BasePath: basePath,
		Workers:  runtime.GOMAXPROCS(0),
	}
}

// LocalCollectorConfig configures the local collector.
type LocalCollectorConfig struct {
	BasePath string
	Workers  int
	Verbose  bool
}

// NewLocalCollectorWithConfig creates a new LocalCollector with configuration.
func NewLocalCollectorWithConfig(cfg LocalCollectorConfig) *LocalCollector {
	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	return &LocalCollector{
		BasePath: cfg.BasePath,
		Workers:  workers,
		Verbose:  cfg.Verbose,
	}
}

// ListRepos scans local directories and returns repositories matching the filter.
// For local scanning, orgs parameter is treated as subdirectory names.
func (c *LocalCollector) ListRepos(_ context.Context, orgs []string, filter model.RepoFilter) ([]model.Repo, error) {
	var repos []model.Repo

	for _, org := range orgs {
		orgPath := filepath.Join(c.BasePath, org)
		orgRepos, err := c.scanDirectory(orgPath, org, filter)
		if err != nil {
			if c.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to scan %s: %v\n", orgPath, err)
			}
			continue
		}
		repos = append(repos, orgRepos...)
	}

	return repos, nil
}

// ListUserRepos scans local directories for user repositories.
// For local scanning, users parameter is treated as subdirectory names.
func (c *LocalCollector) ListUserRepos(ctx context.Context, users []string, filter model.RepoFilter) ([]model.Repo, error) {
	return c.ListRepos(ctx, users, filter)
}

// ListReposMultiSource combines org and user repository lists.
func (c *LocalCollector) ListReposMultiSource(ctx context.Context, orgs, users []string, filter model.RepoFilter) ([]model.Repo, error) {
	var repos []model.Repo

	if len(orgs) > 0 {
		orgRepos, err := c.ListRepos(ctx, orgs, filter)
		if err != nil {
			return nil, err
		}
		repos = append(repos, orgRepos...)
	}

	if len(users) > 0 {
		userRepos, err := c.ListUserRepos(ctx, users, filter)
		if err != nil {
			return nil, err
		}
		repos = append(repos, userRepos...)
	}

	return repos, nil
}

// scanDirectory scans a directory for repositories.
func (c *LocalCollector) scanDirectory(dirPath, owner string, filter model.RepoFilter) ([]model.Repo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dirPath, err)
	}

	// Filter to directories only
	var dirs []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		dirs = append(dirs, entry)
	}

	total := len(dirs)
	if total == 0 {
		return nil, nil
	}

	// Set up worker pool
	numWorkers := min(c.Workers, total)
	numWorkers = max(numWorkers, 1)

	type workItem struct {
		index int
		entry os.DirEntry
	}
	type resultItem struct {
		index int
		repo  *model.Repo
	}

	workCh := make(chan workItem, total)
	resultCh := make(chan resultItem, total)

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range workCh {
				repoPath := filepath.Join(dirPath, work.entry.Name())
				repo := c.analyzeRepo(repoPath, owner, work.entry.Name(), filter)
				resultCh <- resultItem{index: work.index, repo: repo}
			}
		}()
	}

	// Send work
	go func() {
		for i, entry := range dirs {
			workCh <- workItem{index: i, entry: entry}
		}
		close(workCh)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	repos := make([]*model.Repo, total)
	for item := range resultCh {
		repos[item.index] = item.repo
	}

	// Filter out nils and convert to slice
	var result []model.Repo
	for _, r := range repos {
		if r != nil {
			result = append(result, *r)
		}
	}

	return result, nil
}

// analyzeRepo analyzes a single repository directory.
func (c *LocalCollector) analyzeRepo(repoPath, owner, name string, filter model.RepoFilter) *model.Repo {
	// Check if it's a valid repo (has go.mod or .github directory)
	hasGoMod := fileExists(filepath.Join(repoPath, "go.mod"))
	hasGitHub := dirExists(filepath.Join(repoPath, ".github"))

	if !hasGoMod && !hasGitHub {
		return nil
	}

	// Detect languages
	languages := c.detectLanguages(repoPath)

	// Apply language filter
	if len(filter.IncludeLanguages) > 0 {
		found := false
		for _, lang := range languages {
			for _, filterLang := range filter.IncludeLanguages {
				if strings.EqualFold(lang, filterLang) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return nil
		}
	}

	// Check for archived (we can't really detect this locally, skip filter)
	// Check for fork (we can't detect this locally either)

	return &model.Repo{
		Owner:     owner,
		Name:      name,
		FullName:  owner + "/" + name,
		Languages: languages,
		HTMLURL:   "file://" + repoPath,
		LocalPath: repoPath,
	}
}

// detectLanguages detects programming languages in a repository.
func (c *LocalCollector) detectLanguages(repoPath string) []string {
	var languages []string

	// Check for Go
	if fileExists(filepath.Join(repoPath, "go.mod")) {
		languages = append(languages, "Go")
	}

	// Check for TypeScript/JavaScript
	if fileExists(filepath.Join(repoPath, "package.json")) {
		// Check if it's TypeScript
		if fileExists(filepath.Join(repoPath, "tsconfig.json")) {
			languages = append(languages, "TypeScript")
		} else {
			languages = append(languages, "JavaScript")
		}
	}

	// Check for Crystal
	if fileExists(filepath.Join(repoPath, "shard.yml")) {
		languages = append(languages, "Crystal")
	}

	// Check for Python
	if fileExists(filepath.Join(repoPath, "pyproject.toml")) ||
		fileExists(filepath.Join(repoPath, "setup.py")) ||
		fileExists(filepath.Join(repoPath, "requirements.txt")) {
		languages = append(languages, "Python")
	}

	// Check for Rust
	if fileExists(filepath.Join(repoPath, "Cargo.toml")) {
		languages = append(languages, "Rust")
	}

	return languages
}

// GetWorkflows returns CI/CD workflows for a repository.
func (c *LocalCollector) GetWorkflows(_ context.Context, repo model.Repo) ([]model.Workflow, error) {
	repoPath := repo.LocalPath
	if repoPath == "" {
		repoPath = filepath.Join(c.BasePath, repo.Owner, repo.Name)
	}

	workflowsDir := filepath.Join(repoPath, ".github", "workflows")
	if !dirExists(workflowsDir) {
		return nil, nil
	}

	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return nil, fmt.Errorf("reading workflows directory: %w", err)
	}

	var workflows []model.Workflow
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		workflowPath := filepath.Join(workflowsDir, name)
		workflow, err := c.parseWorkflow(workflowPath, name)
		if err != nil {
			if c.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", workflowPath, err)
			}
			continue
		}

		workflows = append(workflows, *workflow)
	}

	return workflows, nil
}

// workflowYAML represents the structure of a GitHub Actions workflow file.
type workflowYAML struct {
	Name string                 `yaml:"name"`
	On   any                    `yaml:"on"`
	Jobs map[string]workflowJob `yaml:"jobs"`
}

type workflowJob struct {
	Name   string         `yaml:"name"`
	Uses   string         `yaml:"uses"`
	RunsOn any            `yaml:"runs-on"`
	Steps  []workflowStep `yaml:"steps"`
	Needs  any            `yaml:"needs"`
}

type workflowStep struct {
	Name string            `yaml:"name"`
	Uses string            `yaml:"uses"`
	Run  string            `yaml:"run"`
	With map[string]string `yaml:"with"`
}

// parseWorkflow parses a workflow YAML file.
func (c *LocalCollector) parseWorkflow(path, filename string) (*model.Workflow, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wf workflowYAML
	if err := yaml.Unmarshal(content, &wf); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	workflow := &model.Workflow{
		Name:    wf.Name,
		Path:    ".github/workflows/" + filename,
		Content: string(content),
	}

	// Extract triggers
	workflow.Triggers = extractTriggers(wf.On)

	// Extract jobs and reusable workflow references
	for jobID, job := range wf.Jobs {
		wfJob := model.WorkflowJob{
			ID:   jobID,
			Name: job.Name,
		}

		// Check if job uses a reusable workflow
		if job.Uses != "" {
			wfJob.UsesReusableWorkflow = true
			ref := model.ParseReusableWorkflowRef(job.Uses)
			wfJob.ReusableWorkflowRef = ref
			workflow.UsesReusableWorkflow = true
			workflow.ReusableWorkflowRefs = append(workflow.ReusableWorkflowRefs, *ref)
		}

		// Extract runs-on
		wfJob.RunsOn = extractRunsOn(job.RunsOn)

		// Extract steps
		for _, step := range job.Steps {
			wfJob.Steps = append(wfJob.Steps, model.WorkflowStep{
				Name: step.Name,
				Uses: step.Uses,
				Run:  step.Run,
				With: step.With,
			})
		}

		workflow.Jobs = append(workflow.Jobs, wfJob)
	}

	return workflow, nil
}

// extractTriggers extracts trigger names from the 'on' field.
func extractTriggers(on any) []string {
	var triggers []string

	switch v := on.(type) {
	case string:
		triggers = append(triggers, v)
	case []any:
		for _, t := range v {
			if s, ok := t.(string); ok {
				triggers = append(triggers, s)
			}
		}
	case map[string]any:
		for k := range v {
			triggers = append(triggers, k)
		}
	}

	return triggers
}

// extractRunsOn extracts runs-on values.
func extractRunsOn(runsOn any) []string {
	var result []string

	switch v := runsOn.(type) {
	case string:
		result = append(result, v)
	case []any:
		for _, r := range v {
			if s, ok := r.(string); ok {
				result = append(result, s)
			}
		}
	}

	return result
}

// GetBranchProtection returns branch protection settings.
// Not applicable for local filesystem.
func (c *LocalCollector) GetBranchProtection(_ context.Context, _ model.Repo, _ string) (*model.BranchProtection, error) {
	return nil, nil
}

// GetLatestWorkflowRun returns the most recent workflow run.
// Not applicable for local filesystem.
func (c *LocalCollector) GetLatestWorkflowRun(_ context.Context, _ model.Repo, _ int64) (*model.WorkflowRun, error) {
	return nil, nil
}

// GetFileContent returns the content of a file from a repository.
func (c *LocalCollector) GetFileContent(_ context.Context, repo model.Repo, path string) (string, error) {
	repoPath := repo.LocalPath
	if repoPath == "" {
		repoPath = filepath.Join(c.BasePath, repo.Owner, repo.Name)
	}

	filePath := filepath.Join(repoPath, path)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetLanguages returns the languages used in a repository.
func (c *LocalCollector) GetLanguages(_ context.Context, repo model.Repo) ([]string, error) {
	repoPath := repo.LocalPath
	if repoPath == "" {
		repoPath = filepath.Join(c.BasePath, repo.Owner, repo.Name)
	}

	return c.detectLanguages(repoPath), nil
}

// ParseGoMod parses a go.mod file and returns module name, replace count, and dependencies.
func ParseGoMod(goModPath string) (moduleName string, replaceCount int, dependencies []string) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", 0, nil
	}
	defer func() { _ = file.Close() }()

	s := bufio.NewScanner(file)
	inReplaceBlock := false
	inRequireBlock := false

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		// Get module name
		if mod, found := strings.CutPrefix(line, "module "); found {
			moduleName = strings.TrimSpace(mod)
		}

		// Count replace directives
		if strings.HasPrefix(line, "replace ") && !strings.HasPrefix(line, "replace (") {
			replaceCount++
		}

		// Handle replace block
		if strings.HasPrefix(line, "replace (") {
			inReplaceBlock = true
			continue
		}
		if inReplaceBlock {
			if line == ")" {
				inReplaceBlock = false
				continue
			}
			if line != "" && !strings.HasPrefix(line, "//") {
				replaceCount++
			}
		}

		// Parse single-line require
		if strings.HasPrefix(line, "require ") && !strings.HasPrefix(line, "require (") {
			if dep := parseRequireLine(strings.TrimPrefix(line, "require ")); dep != "" {
				dependencies = append(dependencies, dep)
			}
		}

		// Handle require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock {
			if line == ")" {
				inRequireBlock = false
				continue
			}
			if dep := parseRequireLine(line); dep != "" {
				dependencies = append(dependencies, dep)
			}
		}
	}

	return moduleName, replaceCount, dependencies
}

// parseRequireLine extracts the module path from a require line.
func parseRequireLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") {
		return ""
	}
	parts := strings.Fields(line)
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// Helper functions

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
