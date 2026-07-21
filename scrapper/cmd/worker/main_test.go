package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestReadSitesFileSkipsCommentsAndBlanks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sites.txt")
	content := "# главные площадки\nhttps://a.example\n\n  https://b.example  \n# ещё\nhttps://c.example\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := readSitesFile(path)
	if err != nil {
		t.Fatalf("readSitesFile() returned error: %v", err)
	}

	want := []string{"https://a.example", "https://b.example", "https://c.example"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readSitesFile() = %v, want %v", got, want)
	}
}

func TestReadSitesFileMissingFile(t *testing.T) {
	if _, err := readSitesFile(filepath.Join(t.TempDir(), "nope.txt")); err == nil {
		t.Fatal("readSitesFile() expected error for missing file, got nil")
	}
}

func TestEnvHelpers(t *testing.T) {
	t.Setenv("TEST_STR", "value")
	t.Setenv("TEST_INT", "7")
	t.Setenv("TEST_INT_BAD", "seven")
	t.Setenv("TEST_DUR", "30m")
	t.Setenv("TEST_DUR_BAD", "-5s")

	if got := envOr("TEST_STR", "def"); got != "value" {
		t.Errorf("envOr() = %q, want value", got)
	}
	if got := envOr("TEST_MISSING", "def"); got != "def" {
		t.Errorf("envOr() = %q, want def", got)
	}
	if got := envIntOr("TEST_INT", 1); got != 7 {
		t.Errorf("envIntOr() = %d, want 7", got)
	}
	if got := envIntOr("TEST_INT_BAD", 1); got != 1 {
		t.Errorf("envIntOr() with bad value = %d, want fallback 1", got)
	}
	if got := envDurationOr("TEST_DUR", time.Hour); got != 30*time.Minute {
		t.Errorf("envDurationOr() = %v, want 30m", got)
	}
	if got := envDurationOr("TEST_DUR_BAD", time.Hour); got != time.Hour {
		t.Errorf("envDurationOr() with negative value = %v, want fallback 1h", got)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	cfg := loadConfig()

	if cfg.kafkaTopic != "internships" {
		t.Errorf("kafkaTopic = %q, want internships", cfg.kafkaTopic)
	}
	if cfg.concurrency != 3 {
		t.Errorf("concurrency = %d, want 3", cfg.concurrency)
	}
	if cfg.scrapeInterval != 6*time.Hour {
		t.Errorf("scrapeInterval = %v, want 6h", cfg.scrapeInterval)
	}
}
