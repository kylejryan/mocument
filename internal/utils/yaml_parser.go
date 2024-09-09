package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

type FieldSchema struct {
	Type    string   `yaml:"type"`
	Format  string   `yaml:"format,omitempty"`
	Min     float64  `yaml:"min,omitempty"`
	Max     float64  `yaml:"max,omitempty"`
	Options []string `yaml:"options,omitempty"`
}

type TableSchema struct {
	Fields map[string]FieldSchema `yaml:"fields"`
}

type Schema struct {
	Tables map[string]TableSchema `yaml:"tables"`
}

func LoadYAMLConfig(filename string) (*Schema, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var schema Schema
	err = yaml.Unmarshal(data, &schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
