package helpers

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadAndParseYamlFile(filename string, dest interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, dest)
}

func WriteYamlFile(filename string, src interface{}) error {
	raw, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, raw, 0664)
}
