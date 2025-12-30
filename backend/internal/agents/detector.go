package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Agent struct {
	Name             string
	ConfigPath       string
	ConfigPaths      []string // Multiple possible config paths
	BinaryNames      []string // Multiple binary names to check
	Installed        bool
	Configured       bool // True if configured with our proxy
	ConfigPathExists bool // True if config file exists
}

// Common binary paths to search (like quotio reference)
var commonBinaryPaths = []string{
	"/usr/local/bin",
	"/opt/homebrew/bin",
	"/usr/bin",
	"~/.local/bin",
	"~/.cargo/bin",
	"~/.bun/bin",
	"~/.deno/bin",
	"~/.npm-global/bin",
	"~/.volta/bin",
	"~/.asdf/shims",
	"~/.local/share/mise/shims",
	"~/.opencode/bin",
}

var KnownAgents = []Agent{
	{
		Name:        "claude-code",
		ConfigPath:  getClaudeCodeConfigPath(),
		ConfigPaths: getClaudeCodeConfigPaths(),
		BinaryNames: []string{"claude"},
	},
	{
		Name:        "codex",
		ConfigPath:  getCodexConfigPath(),
		ConfigPaths: getCodexConfigPaths(),
		BinaryNames: []string{"codex"},
	},
	{
		Name:        "gemini-cli",
		ConfigPath:  getGeminiCLIConfigPath(),
		ConfigPaths: []string{}, // Environment-based, no config file
		BinaryNames: []string{"gemini"},
	},
	{
		Name:        "amp-cli",
		ConfigPath:  getAmpCLIConfigPath(),
		ConfigPaths: getAmpCLIConfigPaths(),
		BinaryNames: []string{"amp"},
	},
	{
		Name:        "opencode",
		ConfigPath:  getOpenCodeConfigPath(),
		ConfigPaths: getOpenCodeConfigPaths(),
		BinaryNames: []string{"opencode", "oc"},
	},
	{
		Name:        "droid",
		ConfigPath:  getDroidConfigPath(),
		ConfigPaths: getDroidConfigPaths(),
		BinaryNames: []string{"droid", "factory-droid", "fd"},
	},
}

func DetectAgents() ([]Agent, error) {
	var detected []Agent

	for _, agent := range KnownAgents {
		// Check if command exists using improved detection
		installed := checkAgentInstalledImproved(agent.BinaryNames)

		// Check if config file exists (check all possible paths)
		configExists := false
		configPath := agent.ConfigPath
		for _, path := range agent.ConfigPaths {
			if exists, _ := checkConfigExists(path); exists {
				configExists = true
				configPath = path
				break
			}
		}
		if !configExists && agent.ConfigPath != "" {
			configExists, _ = checkConfigExists(agent.ConfigPath)
		}

		// Check if already configured with our proxy
		configured := false
		if configExists {
			configured = checkAgentConfigured(configPath)
		}

		detected = append(detected, Agent{
			Name:             agent.Name,
			ConfigPath:       configPath,
			Installed:        installed,
			Configured:       configured,
			ConfigPathExists: configExists,
		})
	}

	return detected, nil
}

func checkConfigExists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	expandedPath := expandPath(path)
	_, err := os.Stat(expandedPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkAgentConfigured(configPath string) bool {
	expandedPath := expandPath(configPath)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return false
	}

	content := string(data)

	// Check for proxy configuration in content (works for JSON, TOML, etc.)
	if strings.Contains(content, "localhost:") || 
	   strings.Contains(content, "127.0.0.1:") ||
	   strings.Contains(content, "cliproxyapi") {
		return true
	}

	// Also try JSON parsing for proxy_url field
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err == nil {
		if proxyURL, ok := config["proxy_url"].(string); ok {
			return strings.Contains(proxyURL, "localhost:") || strings.Contains(proxyURL, "127.0.0.1:")
		}
	}

	return false
}

// checkAgentInstalledImproved checks multiple binary names and common paths
func checkAgentInstalledImproved(binaryNames []string) bool {
	home, _ := os.UserHomeDir()

	for _, name := range binaryNames {
		// 1. Try which/LookPath first (works if PATH is set correctly)
		if _, err := exec.LookPath(name); err == nil {
			return true
		}

		// 2. Check common binary paths
		for _, basePath := range commonBinaryPaths {
			expandedBase := strings.ReplaceAll(basePath, "~", home)
			fullPath := filepath.Join(expandedBase, name)
			if isExecutable(fullPath) {
				return true
			}
		}

		// 3. Check version manager paths (nvm, fnm)
		for _, path := range getVersionManagerPaths(name, home) {
			if isExecutable(path) {
				return true
			}
		}
	}
	return false
}

// isExecutable checks if a file exists and is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if it's a file and has execute permission
	return !info.IsDir() && info.Mode()&0111 != 0
}

// getVersionManagerPaths returns paths from version managers like nvm/fnm
func getVersionManagerPaths(name, home string) []string {
	var paths []string

	// nvm: ~/.nvm/versions/node/v*/bin/
	nvmBase := filepath.Join(home, ".nvm", "versions", "node")
	if versions, err := os.ReadDir(nvmBase); err == nil {
		for i := len(versions) - 1; i >= 0; i-- { // Prefer newer versions
			paths = append(paths, filepath.Join(nvmBase, versions[i].Name(), "bin", name))
		}
	}

	// fnm: ~/.fnm/node-versions/v*/installation/bin/
	fnmBase := filepath.Join(home, ".fnm", "node-versions")
	if versions, err := os.ReadDir(fnmBase); err == nil {
		for i := len(versions) - 1; i >= 0; i-- {
			paths = append(paths, filepath.Join(fnmBase, versions[i].Name(), "installation", "bin", name))
		}
	}

	return paths
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}

// Legacy function for backward compatibility
func checkAgentInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Config path functions - Updated based on quotio reference

func getClaudeCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "claude", "settings.json")
	}
	return filepath.Join(home, ".claude", "settings.json")
}

func getClaudeCodeConfigPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return []string{filepath.Join(home, "AppData", "Roaming", "claude", "settings.json")}
	}
	return []string{filepath.Join(home, ".claude", "settings.json")}
}

func getCodexConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "codex", "config.toml")
	}
	return filepath.Join(home, ".codex", "config.toml")
}

func getCodexConfigPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return []string{
			filepath.Join(home, "AppData", "Roaming", "codex", "config.toml"),
			filepath.Join(home, "AppData", "Roaming", "codex", "auth.json"),
		}
	}
	return []string{
		filepath.Join(home, ".codex", "config.toml"),
		filepath.Join(home, ".codex", "auth.json"),
	}
}

func getOpenCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "opencode", "opencode.json")
	}
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

func getOpenCodeConfigPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return []string{filepath.Join(home, "AppData", "Roaming", "opencode", "opencode.json")}
	}
	return []string{filepath.Join(home, ".config", "opencode", "opencode.json")}
}

func getGeminiCLIConfigPath() string {
	// Gemini CLI uses environment variables, not config files
	// Return empty - configuration is done via environment
	return ""
}

func getDroidConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "factory", "config.json")
	}
	return filepath.Join(home, ".factory", "config.json")
}

func getDroidConfigPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return []string{filepath.Join(home, "AppData", "Roaming", "factory", "config.json")}
	}
	return []string{filepath.Join(home, ".factory", "config.json")}
}

func getAmpCLIConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "amp", "settings.json")
	}
	return filepath.Join(home, ".config", "amp", "settings.json")
}

func getAmpCLIConfigPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return []string{
			filepath.Join(home, "AppData", "Roaming", "amp", "settings.json"),
			filepath.Join(home, "AppData", "Local", "amp", "secrets.json"),
		}
	}
	return []string{
		filepath.Join(home, ".config", "amp", "settings.json"),
		filepath.Join(home, ".local", "share", "amp", "secrets.json"),
	}
}

func ConfigureAgent(agentName, proxyURL string) error {
	// Find agent config
	var agent *Agent
	for i := range KnownAgents {
		if KnownAgents[i].Name == agentName {
			agent = &KnownAgents[i]
			break
		}
	}

	if agent == nil {
		return fmt.Errorf("agent not found: %s", agentName)
	}

	// Special handling for gemini-cli (environment-based)
	if agentName == "gemini-cli" {
		// Gemini CLI doesn't use config files, return success
		// User should set environment variables manually
		return nil
	}

	configPath := expandPath(agent.ConfigPath)
	if configPath == "" {
		return fmt.Errorf("no config path for agent: %s", agentName)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Read existing config if exists
	var existingConfig map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &existingConfig)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Update proxy URL
	if existingConfig == nil {
		existingConfig = make(map[string]interface{})
	}
	existingConfig["proxy_url"] = proxyURL

	// Write updated config
	configData, err := json.MarshalIndent(existingConfig, "", "  ")
	if err != nil {
		return err
	}

	// Create a backup of the existing config before writing
	if err := createConfigBackup(configPath); err != nil {
		return fmt.Errorf("failed to create config backup: %w", err)
	}

	// Write the new config
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return err
	}

	return nil
}

// createConfigBackup creates a .bak backup of the config file if it exists
func createConfigBackup(configPath string) error {
	// Check if config file exists
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, no backup needed
			return nil
		}
		return err
	}

	// Read the existing config
	originalData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config for backup: %w", err)
	}

	// Create backup path
	backupPath := configPath + ".bak"

	// Write backup file with same permissions as original
	if err := os.WriteFile(backupPath, originalData, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

func GetAgentStatus(agentName string) (bool, error) {
	for _, agent := range KnownAgents {
		if agent.Name == agentName {
			configPath := expandPath(agent.ConfigPath)
			if configPath == "" {
				return false, nil
			}
			_, err := os.Stat(configPath)
			return err == nil, nil
		}
	}
	return false, fmt.Errorf("agent not found: %s", agentName)
}

func ValidateAgentConfig(configPath string) (bool, string) {
	if configPath == "" {
		// Environment-based config (like gemini-cli) is always "valid"
		return true, ""
	}

	expandedPath := expandPath(configPath)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return false, fmt.Sprintf("Cannot read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		// Not JSON - might be TOML (like codex), check for proxy content
		content := string(data)
		if strings.Contains(content, "localhost") || strings.Contains(content, "127.0.0.1") {
			return true, ""
		}
		return false, fmt.Sprintf("Invalid JSON in config file: %v", err)
	}

	// Check for required fields
	if _, ok := config["proxy_url"]; !ok {
		return false, "proxy_url field is missing"
	}

	return true, ""
}
