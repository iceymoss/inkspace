package config

import (
	"os"
	"path/filepath"
	"testing"
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
