package expframework

import (
	"encoding/json"
	"os"
)

func writeToJsonFile(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	jsonEncoder := json.NewEncoder(f)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	err = jsonEncoder.Encode(v)
	if err != nil {
		return err
	}

	return nil
}
