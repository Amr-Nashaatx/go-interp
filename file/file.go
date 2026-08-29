package file

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return "", fmt.Errorf("file %q does not exist: %w", path, err)
		case errors.Is(err, fs.ErrPermission):
			return "", fmt.Errorf("no permission to read %q: %w", path, err)
		default:
			return "", fmt.Errorf("reading %q: %w", path, err)
		}
	}

	return string(data), nil
}
