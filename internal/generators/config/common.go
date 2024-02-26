package config

type DefaultConfig map[string]string

type File struct {
	Name    string   `yaml:"name"`
	Tmpl    string   `yaml:"tmpl,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
}