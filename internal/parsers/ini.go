package parsers

import (
	"fmt"
	"gopkg.in/ini.v1"
)

// Parses an INI file into a flat map of keys mapped to values. This reduces
// the complexity of the INI file to a simple key/value store and ignores the
// sections.
func ParseINIFile(filePath string) (map[string]string, error) {

	iniFile, err := ini.Load(filePath)

	if err != nil {
		return nil, fmt.Errorf("Failed to read the INI file %s because %v", filePath, err)
	}

	data := make(map[string]string)
	for _, section := range iniFile.Sections() {
		for key, value := range section.KeysHash() {
			data[key] = value
		}
	}
	return data, nil
}
