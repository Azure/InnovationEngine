package lib

import (
	"fmt"
	"os"
	"os/user"
)

func GetHomeDirectory() (string, error) {
	// Try to get home directory from user.Current()
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir, nil
	}

	// Fallback to environment variable
	home, exists := os.LookupEnv("HOME")
	if exists && home != "" {
		return home, nil
	}

	// Fallback for Windows
	homeDrive, driveExists := os.LookupEnv("HOMEDRIVE")
	homePath, pathExists := os.LookupEnv("HOMEPATH")
	if driveExists && pathExists {
		return homeDrive + homePath, nil
	}

	return "", fmt.Errorf("Home directory cannot be determined")
}
