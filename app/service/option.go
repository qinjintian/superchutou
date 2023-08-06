package service

type Config struct {
	Authenticate *Authenticate `yaml:"authenticate"`

	Log *Log `yaml:"log"`
}

type Authenticate struct {
	Cookie  string `yaml:"cookie"`
}

type Log struct {
	Name   string `yaml:"name"`
	Output []*Output
}

type Output struct {
	Type  string `yaml:"type"`
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}
