package actions

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

type Real struct{}

func (Real) Kill(pid int32, force bool) error {
	if pid <= 0 {
		return fmt.Errorf("no process id available")
	}
	p, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	if force {
		return p.Kill()
	}
	return p.Signal(os.Interrupt)
}

func (Real) Copy(value string) error {
	return clipboard.WriteAll(value)
}

func (Real) Open(url string) error {
	return browser.OpenURL(url)
}
