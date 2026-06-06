package platform

// AutoStart API — platform-specific implementations are in
// autostart_darwin.go and autostart_windows.go
//
// These functions are implemented via build tags:
//   - autostart_darwin.go  (build tag: darwin) — macOS LaunchAgent
//   - autostart_windows.go (build tag: windows) — Windows Registry

const (
	// AppName is used as the identifier for auto-start registration
	AppName = "com.chaitin.workmanager"
	// AppDisplayName is the human-readable name
	AppDisplayName = "工作安排管理系统"
)
