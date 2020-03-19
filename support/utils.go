package support

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/project-flogo/core/support/log"
)

// HandlePanic helper method to handle panics
//todo should we remove this
func HandlePanic(name string, err *error) {
	if r := recover(); r != nil {

		log.RootLogger().Warnf("%s: PANIC Occurred  : %v\n", name, r)

		// todo: useful for debugging
		log.RootLogger().Debugf("StackTrace: %s", debug.Stack())

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// URLStringToFilePath convert fileURL to file path
func URLStringToFilePath(fileURL string) (string, bool) {

	if strings.HasPrefix(fileURL, "file://") {

		filePath := fileURL[7:]

		if runtime.GOOS == "windows" {
			if strings.HasPrefix(filePath, "/") {
				filePath = filePath[1:]
			}
			filePath = filepath.FromSlash(filePath)
		}

		filePath = strings.Replace(filePath, "%20", " ", -1)

		return filePath, true
	}

	return "", false
}
