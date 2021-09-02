package devtest

import (
	"os"
	"os/exec"
	"testing"
)

func TestVet(t *testing.T) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}
}
