package engram

import (
	"reflect"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/system"
)

func TestInstallCommandByProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile system.PlatformProfile
		want    [][]string
		wantErr bool
	}{
		{
			name:    "darwin uses brew tap and install",
			profile: system.PlatformProfile{OS: "darwin", PackageManager: "brew"},
			want:    [][]string{{"brew", "tap", "Gentleman-Programming/homebrew-tap"}, {"brew", "install", "engram"}},
		},
		{
			name:    "ubuntu uses go install with correct module path",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroUbuntu, PackageManager: "apt"},
			want: [][]string{
				{"env", "CGO_ENABLED=0", "go", "install", "github.com/Gentleman-Programming/engram/cmd/engram@latest"},
				{"sh", "-c", "sudo ln -sf \"$(go env GOPATH | cut -d: -f1)/bin/engram\" /usr/local/bin/engram || echo '\\n\\033[33mWARNING: Could not symlink engram to /usr/local/bin. Please ensure $(go env GOPATH | cut -d: -f1)/bin is in your PATH\\033[0m\\n'"},
			},
		},
		{
			name:    "arch uses go install with correct module path",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroArch, PackageManager: "pacman"},
			want: [][]string{
				{"env", "CGO_ENABLED=0", "go", "install", "github.com/Gentleman-Programming/engram/cmd/engram@latest"},
				{"sh", "-c", "sudo ln -sf \"$(go env GOPATH | cut -d: -f1)/bin/engram\" /usr/local/bin/engram || echo '\\n\\033[33mWARNING: Could not symlink engram to /usr/local/bin. Please ensure $(go env GOPATH | cut -d: -f1)/bin is in your PATH\\033[0m\\n'"},
			},
		},
		{
			name:    "unsupported package manager errors",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: "fedora", PackageManager: "dnf"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := InstallCommand(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("InstallCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(command, tt.want) {
				t.Fatalf("InstallCommand() = %v, want %v", command, tt.want)
			}
		})
	}
}
