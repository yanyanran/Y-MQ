package tools

import (
	"fmt"
	"runtime"
)

func String(app string) string {
	return fmt.Sprintf("%s (built w/%s)", app, runtime.Version())
}
