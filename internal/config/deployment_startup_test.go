package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeploymentStartupRejectsInvalidCanonical exercises the actual binaries:
// configuration errors must terminate before any MySQL/Redis initialization.
func TestDeploymentStartupRejectsInvalidCanonical(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, binary := range []string{"server", "admin", "scheduler"} {
		t.Run(binary, func(t *testing.T) {
			executable := filepath.Join(t.TempDir(), binary)
			build := exec.Command("go", "build", "-o", executable, "./cmd/"+binary)
			build.Dir = root
			if output, err := build.CombinedOutput(); err != nil {
				t.Fatalf("build: %v\n%s", err, output)
			}
			for _, mode := range []string{"debug", "test", "release"} {
				for _, raw := range []string{"<unset>", "", "https://example.com?", "https://example.com#", "https://example.com?x=1", "https://example.com#frag", "javascript:bad"} {
					cmd := exec.Command(executable)
					cmd.Dir = t.TempDir()
					cmd.Env = append(withoutCanonicalEnvironment(os.Environ()), "SERVER_MODE="+mode, "ADMIN_MODE="+mode)
					if raw != "<unset>" {
						cmd.Env = append(cmd.Env, "PUBLIC_SITE_URL="+raw)
					}
					output, err := cmd.CombinedOutput()
					if err == nil || !strings.Contains(string(output), "PUBLIC_SITE_URL") || strings.Contains(string(output), "初始化数据库") || strings.Contains(string(output), "数据库连接失败") {
						t.Fatalf("%s/%q: expected config rejection before DB: %v\n%s", mode, raw, err, output)
					}
				}
			}
			if binary == "admin" {
				dir := t.TempDir()
				if err := os.Mkdir(filepath.Join(dir, "config"), 0o700); err != nil {
					t.Fatal(err)
				}
				for name, origin := range map[string]string{"admin": "https://invalid.example?query", "config": "https://valid.example"} {
					if err := os.WriteFile(filepath.Join(dir, "config", name+".yaml"), []byte("public_site_url: "+origin+"\n"), 0o600); err != nil {
						t.Fatal(err)
					}
				}
				cmd := exec.Command(executable)
				cmd.Dir, cmd.Env = dir, withoutCanonicalEnvironment(os.Environ())
				output, err := cmd.CombinedOutput()
				if err == nil || !strings.Contains(string(output), "加载配置失败") || strings.Contains(string(output), "初始化数据库失败") {
					t.Fatalf("invalid admin canonical must not retry valid config.yaml: %v\n%s", err, output)
				}
			}
		})
	}
}

func withoutCanonicalEnvironment(env []string) []string {
	var filtered []string
	for _, value := range env {
		if !strings.HasPrefix(value, "PUBLIC_SITE_URL=") && !strings.HasPrefix(value, "SERVER_MODE=") && !strings.HasPrefix(value, "ADMIN_MODE=") && !strings.HasPrefix(value, "DATABASE_") && !strings.HasPrefix(value, "REDIS_") {
			filtered = append(filtered, value)
		}
	}
	// If validation regresses, never try operator/developer DB credentials.
	return append(filtered, "DATABASE_HOST=127.0.0.1", "DATABASE_PORT=1", "DATABASE_USERNAME=seo-startup-test", "DATABASE_PASSWORD=not-a-secret", "DATABASE_NAME=seo_startup_test", "REDIS_HOST=127.0.0.1", "REDIS_PORT=1")
}
