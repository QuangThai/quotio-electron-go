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
	Name       string
	ConfigPath string
	Installed  bool
	Configured bool   // True if configured with our proxy
	ConfigPathExists bool // True if config file exists
}

var KnownAgents = []Agent{
	{Name: "claude-code", ConfigPath: getClaudeCodeConfigPath()},
	{Name: "opencode", ConfigPath: getOpenCodeConfigPath()},
	{Name: "gemini-cli", ConfigPath: getGeminiCLIConfigPath()},
	{Name: "droid", ConfigPath: getDroidConfigPath()},
	{Name: "amp-cli", ConfigPath: getAmpCLIConfigPath()},
}

func DetectAgents() ([]Agent, error) {
	var detected []Agent

	for _, agent := range KnownAgents {
		// Check if command exists
		installed := checkAgentInstalled(agent.Name)

		// Check if config file exists
		configExists, _ := checkConfigExists(agent.ConfigPath)

		// Check if already configured with our proxy
		configured := false
		if configExists {
			configured = checkAgentConfigured(agent.ConfigPath)
		}

		detected = append(detected, Agent{
			Name:            agent.Name,
			ConfigPath:       agent.ConfigPath,
			Installed:        installed,
			Configured:       configured,
			ConfigPathExists: configExists,
		})
	}

	return detected, nil
}

func checkConfigExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkAgentConfigured(configPath string) bool {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return false
	}

	// Check if proxy_url is set and contains localhost (flexible port)
	if proxyURL, ok := config["proxy_url"].(string); ok {
		return strings.Contains(proxyURL, "localhost:") || strings.Contains(proxyURL, "127.0.0.1:")
	}

	return false
}

func checkAgentInstalled(name string) bool {
	// Check if command exists in PATH
	_, err := exec.LookPath(name)
	return err == nil
}

func getClaudeCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "claude-code", "config.json")
	}
	return filepath.Join(home, ".config", "claude-code", "config.json")
}

func getOpenCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "opencode", "config.json")
	}
	return filepath.Join(home, ".config", "opencode", "config.json")
}

func getGeminiCLIConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "gemini-cli", "config.json")
	}
	return filepath.Join(home, ".config", "gemini-cli", "config.json")
}

func getDroidConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "droid", "config.json")
	}
	return filepath.Join(home, ".config", "droid", "config.json")
}

func getAmpCLIConfigPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Roaming", "amp-cli", "config.json")
	}
	return filepath.Join(home, ".config", "amp-cli", "config.json")
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

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(agent.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Read existing config if exists
	var existingConfig map[string]interface{}
	if data, err := os.ReadFile(agent.ConfigPath); err == nil {
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
	if err := createConfigBackup(agent.ConfigPath); err != nil {
		return fmt.Errorf("failed to create config backup: %w", err)
	}

	// Write the new config
	if err := os.WriteFile(agent.ConfigPath, configData, 0644); err != nil {
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
			_, err := os.Stat(agent.ConfigPath)
			return err == nil, nil
		}
	}
	return false, fmt.Errorf("agent not found: %s", agentName)
}

func ValidateAgentConfig(configPath string) (bool, string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false, fmt.Sprintf("Cannot read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return false, fmt.Sprintf("Invalid JSON in config file: %v", err)
	}

	// Check for required fields
	if _, ok := config["proxy_url"]; !ok {
		return false, "proxy_url field is missing"
	}

	return true, ""
}

