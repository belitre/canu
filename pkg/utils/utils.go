package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetScriptPath(executablePath, aliasName string) string {
	binaryPath := filepath.Dir(executablePath)

	scriptPath := path.Join(binaryPath, fmt.Sprintf("%s.sh", aliasName))

	return scriptPath
}

func RemoveAliasFromFile(filename string, alias string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("error while opening file %s: %v", filename, err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var bs []byte

	buf := bytes.NewBuffer(bs)

	var text string

	for scanner.Scan() {
		text = scanner.Text()

		if strings.HasPrefix(text, alias) {
			fmt.Printf("found match for alias: %s, removing...\n", text)
			continue
		}

		if _, err := buf.WriteString(text + "\n"); err != nil {
			return fmt.Errorf("error while writing to buffer: %v", err)
		}
	}

	f.Truncate(0)

	f.Seek(0, 0)

	if _, err := buf.WriteTo(f); err != nil {
		return fmt.Errorf("error while writing buffer to shell config file %s: %v", filename, err)
	}

	return nil
}
