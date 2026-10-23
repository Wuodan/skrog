package wsl

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRunCommandClosesLeakedOutputPipe(t *testing.T) {
	if len(os.Args) > 1 {
		switch os.Args[len(os.Args)-1] {
		case "leak-output":
			child := exec.Command(os.Args[0], "-test.run=TestRunCommandClosesLeakedOutputPipe", "--", "hold-output")
			child.Stdout = os.Stdout
			child.Stderr = os.Stderr
			if err := child.Start(); err != nil {
				os.Exit(2)
			}
			return
		case "hold-output":
			// Longer than the test's WaitDelay, but short enough not to leave an
			// orphaned helper process behind after the test exits.
			time.Sleep(200 * time.Millisecond)
			return
		}
	}

	_, err := runCommand(context.Background(), 25*time.Millisecond, os.Args[0],
		"-test.run=TestRunCommandClosesLeakedOutputPipe", "--", "leak-output")
	if !errors.Is(err, exec.ErrWaitDelay) {
		t.Fatalf("error = %v, want exec.ErrWaitDelay", err)
	}
}
