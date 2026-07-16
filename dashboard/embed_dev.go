//go:build dev

package dashboard

import (
	"fmt"
)

type devDistFS struct{}

func (devDistFS) ReadFile(name string) ([]byte, error) {
	return nil, fmt.Errorf("%s not embedded in dev mode", name)
}

var DistFS devDistFS
