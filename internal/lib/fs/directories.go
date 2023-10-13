package fs

import (
	"errors"
	"os"

	"github.com/Azure/InnovationEngine/internal/logging"
)

func SetWorkingDirectory(directory string) error {
	// Change working directory if specified
	if directory != "" {
		err := os.Chdir(directory)
		if err != nil {
			logging.GlobalLogger.Error("Failed to change working directory", err)
			return err
		}

		logging.GlobalLogger.Infof("Changed directory to %s", directory)
	}
	return nil
}

// Executes a function within a given working directory and restores
// the original working directory when the function completes.
func UsingDirectory(directory string, executor func() error) error {
	originalDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	err = SetWorkingDirectory(directory)
	if err != nil {
		return err
	}

	executionError := executor()
	err = SetWorkingDirectory(originalDirectory)

	return errors.Join(executionError, err)
}
