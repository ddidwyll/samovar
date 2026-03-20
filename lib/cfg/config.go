package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct{ name string }

func (c *Config) GetConfigName() string     { return c.name }
func (c *Config) SetConfigName(name string) { c.name = name }

type ConfigBehavior interface {
	GetConfigName() string
	SetConfigName(name string)
}

func ToJson(c ConfigBehavior) ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func FromJson(c ConfigBehavior, data []byte) error {
	return json.Unmarshal(data, c)
}

func Load(c ConfigBehavior) error {
	filename := c.GetConfigName() + ".json"
	filepath := filepath.Join("config", filename)
	bs, err := os.ReadFile(filepath)

	if err != nil {
		return err
	}

	return FromJson(c, bs)
}
