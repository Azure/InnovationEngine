package fs

import (
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
