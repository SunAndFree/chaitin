//go:build darwin

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// launchAgentDir is the directory where user LaunchAgents are stored
func launchAgentDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents")
}

// plistPath returns the path to the LaunchAgent plist file
func plistPath() string {
	return filepath.Join(launchAgentDir(), AppName+".plist")
}

// EnableAutoStart registers the app to start on login (macOS)
func EnableAutoStart() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Ensure LaunchAgents directory exists
	if err := os.MkdirAll(launchAgentDir(), 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	// Create plist content
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>`, AppName, execPath)

	// Write plist file
	if err := os.WriteFile(plistPath(), []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	// Load the LaunchAgent
	uid := os.Getuid()
	cmd := exec.Command("launchctl", "bootstrap", fmt.Sprintf("gui/%d", uid), plistPath())
	if output, err := cmd.CombinedOutput(); err != nil {
		// launchctl might print warnings that don't affect functionality
		if !strings.Contains(string(output), "already bootstrapped") &&
			!strings.Contains(string(output), "service already loaded") {
			return fmt.Errorf("failed to bootstrap launch agent: %s - %w", string(output), err)
		}
	}

	return nil
}

// DisableAutoStart removes the app from login items (macOS)
func DisableAutoStart() error {
	uid := os.Getuid()

	// Bootout the service if it's loaded
	cmd := exec.Command("launchctl", "bootout", fmt.Sprintf("gui/%d", uid), plistPath())
	cmd.Run() // Ignore errors — the service might not be loaded

	// Remove the plist file
	if err := os.Remove(plistPath()); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	return nil
}

// IsAutoStartEnabled checks if auto-start is enabled (macOS)
func IsAutoStartEnabled() (bool, error) {
	// Check if plist file exists
	_, err := os.Stat(plistPath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check plist file: %w", err)
	}

	// Verify it's actually loaded in launchctl
	cmd := exec.Command("launchctl", "list", AppName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Not found in launchctl — file exists but not loaded
		return false, nil
	}

	return !strings.Contains(string(output), "Could not find"), nil
}
