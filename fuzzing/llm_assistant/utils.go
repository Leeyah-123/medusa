package llm_assistant

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateTestContractName generates a test contract name based on the given contract name.
func generateTestContractName(contractName string) string {
	return fmt.Sprintf("%sTest", contractName)
}

// generateTestFilePath generates a test file path based on the given file path.
func generateTestFilePath(filePath string) string {
	dir, file := filepath.Split(filePath)
	ext := filepath.Ext(file)
	filename := strings.TrimSuffix(file, ext)
	testFilename := fmt.Sprintf("%s_fuzz%s", filename, ext)
	return filepath.Join(dir, testFilename)
}
