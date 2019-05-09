package testaddon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
)

const (
	// ResultDescriptorFileName is the name of the test result descriptor file.
	ResultDescriptorFileName  = "test-info.json"
)

func generateTestInfoFile(dir string, data []byte) error {
	f, err := os.Create(filepath.Join(dir, ResultDescriptorFileName))
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

// ExportArtifacts exports the artifacts in a directory structure rooted at the
// specified directory. The directory where each artifact is exported depends
// on which module and build variant produced it.
func ExportArtifacts(path, baseDir, uniqueDir string) error {
		exportDir := strings.Join([]string{baseDir, uniqueDir}, "/")

		if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
			return fmt.Errorf("skipping artifact (%s): could not ensure unique export dir (%s): %s", path, exportDir, err)
		}

		if _, err := os.Stat(filepath.Join(exportDir, ResultDescriptorFileName)); os.IsNotExist(err) {
			m := map[string]string{"test-name": uniqueDir}
			data, err := json.Marshal(m)
			if err != nil {
				return fmt.Errorf("create test info descriptor: json marshal data (%s): %s", m, err)
			}
			if err := generateTestInfoFile(exportDir, data); err != nil {
				return fmt.Errorf("create test info descriptor: generate file: %s", err)
			}
		}

		name := filepath.Base(path)
		if err := command.CopyFile(path, filepath.Join(exportDir, name)); err != nil {
			return fmt.Errorf("failed to export artifact (%s), error: %v", name, err)
		}
	return nil
}
