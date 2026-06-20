// Package config provides helpers for loading environment configuration.
package config

import (
	"bufio"
	"os"
	"strings"
)

// Load reads a .env file at the given path and sets any key that is not
// already present in the process environment.  The file format is:
//
//   - Blank lines and lines whose first non-space character is '#' are ignored.
//   - Each non-comment line must be KEY=VALUE (spaces around '=' are trimmed).
//   - The value may optionally be wrapped in double or single quotes, which are
//     stripped.
//
// If the file does not exist Load returns silently — a missing file is not an
// error.  Real environment variables always take precedence over file values.
func Load(path string) {
	f, err := os.Open(path)
	if err != nil {
		// File absent or unreadable — not an error for local-dev use.
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on the first '=' only.
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if key == "" {
			continue
		}

		// Strip optional surrounding quotes (double or single).
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		// Only set if the key is not already in the environment — real env wins.
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}
