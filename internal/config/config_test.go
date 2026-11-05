package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadEnvFileDoesNotOverrideProcessEnvironment(t *testing.T) {
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "injected-secret")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("INKSPACE_BOOTSTRAP_SECRET=stale-file-secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	oldWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	loadEnvFile()

	if got := os.Getenv("INKSPACE_BOOTSTRAP_SECRET"); got != "injected-secret" {
		t.Fatalf("bootstrap secret=%q, want injected process value", got)
	}
}

func TestPublicSiteURLRequiredInEveryEnvironment(t *testing.T) {
	old := AppConfig
	t.Cleanup(func() { AppConfig = old; viper.Reset() })
	t.Chdir(t.TempDir())
	for _, mode := range []string{"release", "debug", "test"} {
		for _, name := range []string{"config", "admin"} {
			for _, raw := range []string{"", "https://example.com/path", "https://example.com?", "https://example.com#", "https://example.com?q=x", "https://example.com#fragment", "https://example.com/"} {
				viper.Reset()
				t.Setenv("SERVER_MODE", mode)
				t.Setenv("PUBLIC_SITE_URL", raw)
				err := InitWithFile(name)
				if raw != "https://example.com/" {
					if err == nil {
						t.Fatalf("accepted %q in %s/%s", raw, name, mode)
					}
				} else if err != nil || AppConfig.PublicSiteURL != "https://example.com" {
					t.Fatalf("config %s/%s: %v %+v", name, mode, err, AppConfig)
				}
			}
		}
	}
}

func TestAdminFallbackStillValidatesPublicSiteURL(t *testing.T) {
	old := AppConfig
	t.Cleanup(func() { AppConfig = old; viper.Reset() })
	t.Chdir(t.TempDir())
	if err := os.Mkdir("config", 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("config/config.yaml", []byte("public_site_url: https://yaml.example/\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"", "https://env.example/", "javascript:bad", "https://example.com?", "https://example.com#fragment"} {
		viper.Reset()
		t.Setenv("PUBLIC_SITE_URL", raw)
		err := InitWithFile("admin")
		if raw != "" && raw != "https://env.example/" {
			if err == nil {
				t.Fatal("fallback skipped validation")
			}
			continue
		}
		want := "https://yaml.example"
		if raw != "" {
			want = "https://env.example"
		}
		if err != nil || AppConfig.PublicSiteURL != want {
			t.Fatalf("%v %+v", err, AppConfig)
		}
	}
}

func TestBootstrapSecretOnlyUsesEnvironment(t *testing.T) {
	cfg := &Config{Bootstrap: BootstrapConfig{Secret: "yaml-secret-must-not-apply"}}
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "environment-secret")
	applyEnvironmentOnlyConfig(cfg)
	if cfg.Bootstrap.Secret != "environment-secret" {
		t.Fatalf("bootstrap secret=%q, want environment value", cfg.Bootstrap.Secret)
	}

	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "")
	applyEnvironmentOnlyConfig(cfg)
	if cfg.Bootstrap.Secret != "" {
		t.Fatalf("bootstrap secret=%q, want empty without environment value", cfg.Bootstrap.Secret)
	}
}

func TestAdminFallbackRejectsMissingOrInvalidCanonicalInAllModes(t *testing.T) {
	old := AppConfig
	t.Cleanup(func() { AppConfig = old; viper.Reset() })
	t.Chdir(t.TempDir())
	t.Setenv("PUBLIC_SITE_URL", "")
	if err := os.Mkdir("config", 0o700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"release", "debug", "test"} {
		for _, raw := range []string{"", "https://invalid.example?", "https://invalid.example#fragment"} {
			viper.Reset()
			body := "server:\n  mode: " + mode + "\npublic_site_url: '" + raw + "'\n"
			if err := os.WriteFile("config/config.yaml", []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := InitWithFile("admin"); err == nil {
				t.Fatalf("admin fallback accepted %q in mode %s", raw, mode)
			}
		}
	}
}
