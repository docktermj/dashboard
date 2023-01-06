// The version package simply prints the version of the go-fileindex binary file.
package scan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	log "github.com/docktermj/go-logger/logger"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type FileMetadata struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Volume   string `json:"volume"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
	Modified int64  `json:"modified"`
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func viperAsJson() string {
	viperConfig := viper.AllSettings()
	viperByteArray, err := json.Marshal(viperConfig)
	if err != nil {
		log.Fatalf("Unable to marshal viper config to JSON: %v", err)
	}
	return string(viperByteArray)
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a directory for file metadata.
func Scan(volumeName string, rootPath string, outputFileName string) {

	// Open output file.

	outputFile, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()

	if err != nil {
		log.Errorf("%s output file cannot be created. Error:  %v", outputFileName, err)
		return
	}

	// Recursively descend into directory.

	err = filepath.Walk(rootPath,
		func(path string, fileInfo os.FileInfo, err error) error {

			// If there's a problem, log a warning.

			if err != nil {
				log.Warn(err)
				return nil
			}

			// Don't process directories.

			if fileInfo.IsDir() {
				log.Infof("%s is a directory.", path)
				return nil
			}

			// Don't process anything that isn't a file.

			if !fileInfo.Mode().IsRegular() {
				log.Infof("%s is not a regular file.", path)
				return nil
			}

			// Don't process anything with length 0.

			if fileInfo.Size() == 0 {
				log.Infof("%s is length 0.", path)
				return nil
			}

			// Calculate path relative to volume.

			hostPath := filepath.Dir(path)
			volumePath := strings.TrimPrefix(hostPath, rootPath)

			// Calculate SHA256 value for file.

			file, err := os.Open(path)
			if err != nil {
				log.Warn(err)
				return nil
			}
			defer file.Close()

			shaHash := sha256.New()
			if _, err := io.Copy(shaHash, file); err != nil {
				log.Warn(err)
				return nil
			}

			// Gather information for output.

			fileMetadata := FileMetadata{
				Path:     volumePath,
				Name:     fileInfo.Name(),
				Volume:   volumeName,
				Size:     fileInfo.Size(),
				SHA256:   hex.EncodeToString(shaHash.Sum(nil)),
				Modified: fileInfo.ModTime().Unix(),
			}

			// Create JSON string for output.

			fileMetadataJson, err := json.Marshal(fileMetadata)
			if err != nil {
				log.Fatalf("Unable to marshal file metadata to JSON: %v", err)
			}

			// Write output to file.

			outputFile.WriteString(string(fileMetadataJson) + "\n")

			return nil
		})
	if err != nil {
		log.Warn(err)
	}

}

// ----------------------------------------------------------------------------
// Command pattern "Execute" function.
// ----------------------------------------------------------------------------

// The Command sofware design pattern's Execute() method.
func Execute() {

	// Print context parameters.
	// TODO: Formalize entry parameters

	log.Info(viperAsJson())

	// Get parameters from viper.

	var volumeName = viper.GetString("volume_name")
	var rootPath = viper.GetString("root_path")
	var outputFileName = viper.GetString("output_file_name")

	// Check parameters.

	errors := 0

	if volumeName == "" {
		errors += 1
		log.Warn("--volume-name not set.")
	}

	if rootPath == "" {
		errors += 1
		log.Warn("--root-path not set.")
	}

	if outputFileName == "" {
		errors += 1
		log.Warn("--output-file-name not set.")
	}

	// If any errors, exit.

	if errors > 0 {
		return
	}

	// Perform command.

	Scan(volumeName, rootPath, outputFileName)
}
