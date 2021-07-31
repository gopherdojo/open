package listen

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"
)

func Test_Listen(t *testing.T) {
	t.Helper()

	cleanFn := func() {}

	doneCh := make(chan struct{})
	osExit = func(code int) { doneCh <- struct{}{} }

	output := new(bytes.Buffer)
	_, cancel := Listen(context.Background(), output, cleanFn)
	defer cancel()

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	select {
	case <-doneCh:
		return
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

}
