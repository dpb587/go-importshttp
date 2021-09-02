package devtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	cmdroot := filepath.Join(root, "cmd")
	entries, err := os.ReadDir(cmdroot)
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
			cmd.Dir = filepath.Join(cmdroot, entry.Name())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
		})
	}
}
