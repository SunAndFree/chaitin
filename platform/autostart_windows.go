//go:build windows

package platform

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

// EnableAutoStart registers the app to start on login (Windows)
func EnableAutoStart() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(AppName, execPath); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// DisableAutoStart removes the app from login items (Windows)
func DisableAutoStart() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(AppName); err != nil {
		// If the value doesn't exist, it's not an error
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// IsAutoStartEnabled checks if auto-start is enabled (Windows)
func IsAutoStartEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(AppName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to query registry value: %w", err)
	}

	return true, nil
}
