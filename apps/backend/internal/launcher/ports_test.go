package launcher

import (
	"os"
	"strings"
	"testing"
)

func TestResolvePortsUsesServicePort(t *testing.T) {
	clearLauncherPortEnv(t)
	t.Setenv("KANDEV_SERVER_PORT", "18080")
	t.Setenv("KANDEV_PORT", "18079")

	port, err := resolvePorts(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if port != 18080 {
		t.Fatalf("resolved port = %d, want 18080", port)
	}
}

func TestResolvePortsPrefersBackendPortEnv(t *testing.T) {
	clearLauncherPortEnv(t)
	t.Setenv("KANDEV_BACKEND_PORT", "18081")
	t.Setenv("KANDEV_SERVER_PORT", "18080")
	t.Setenv("KANDEV_PORT", "18079")

	port, err := resolvePorts(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if port != 18081 {
		t.Fatalf("resolved port = %d, want 18081", port)
	}
}

func TestResolvePortsRejectsInvalidServicePort(t *testing.T) {
	clearLauncherPortEnv(t)
	t.Setenv("KANDEV_SERVER_PORT", "invalid")

	_, err := resolvePorts(Options{})
	if err == nil || !strings.Contains(err.Error(), "KANDEV_SERVER_PORT") {
		t.Fatalf("resolvePorts() error = %v, want KANDEV_SERVER_PORT validation error", err)
	}
}

func TestPickAvailablePortExceptSkipsUsedPreferredPort(t *testing.T) {
	port, err := pickAvailablePortExcept(defaultBackendPort, map[int]bool{defaultBackendPort: true})
	if err != nil {
		t.Fatal(err)
	}
	if port == defaultBackendPort {
		t.Fatalf("picked reserved preferred port %d", port)
	}
}

func clearLauncherPortEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{"KANDEV_BACKEND_PORT", "KANDEV_SERVER_PORT", "KANDEV_PORT"} {
		value, exists := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("unset %s: %v", name, err)
		}
		t.Cleanup(func() {
			if exists {
				if err := os.Setenv(name, value); err != nil {
					t.Errorf("restore %s: %v", name, err)
				}
				return
			}
			if err := os.Unsetenv(name); err != nil {
				t.Errorf("unset %s during cleanup: %v", name, err)
			}
		})
	}
}
