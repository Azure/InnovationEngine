package parsers

import (
	"github.com/Azure/InnovationEngine/internal/logging"
	"gopkg.in/ini.v1"
)

// Parses an INI file into a flat map of keys mapped to values. This reduces
// the complexity of the INI file to a simple key/value store and ignores the
// sections.
func ParseINIFile(filePath string) map[string]string {

	iniFile, err := ini.Load(filePath)
	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to read the INI file %s because %v", filePath, err)
	}
	data := make(map[string]string)
	for _, section := range iniFile.Sections() {
		for key, value := range section.KeysHash() {
			data[key] = value
		}
	}
	return data
}
