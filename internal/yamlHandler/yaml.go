package yamlHandler

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Info struct {
	Section1 []Section1Entry `yaml:"section-1"`
	Section2 []Section2Entry `yaml:"section-2"`
}

type Section1Entry struct {
	Title   string        `yaml:"title"`
	Custom  CustomBlock   `yaml:"custom"`
	Content ContentBlock1 `yaml:"content"`
}

type CustomBlock struct {
	Uptime string `yaml:"uptime"` // parse later
}

type ContentBlock1 struct {
	System    map[string][]string `yaml:"system"`
	Languages map[string][]string `yaml:"languages"`
	Hobbies   map[string][]string `yaml:"hobbies"`
}

type Section2Entry struct {
	Title   string        `yaml:"title"`
	Content ContentBlock2 `yaml:"content"`
}

type ContentBlock2 struct {
	Email   map[string]string `yaml:"email"`
	Socials map[string]string `yaml:"socials"`
}

func LoadYAML(path string) (*Info, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var info Info
	err = yaml.Unmarshal(data, &info)
	return &info, err
}
