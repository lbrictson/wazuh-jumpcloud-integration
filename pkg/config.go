package pkg

import (
	"encoding/json"
	"os"
	"time"
)

type ConfigurationData struct {
	APIKey  string     `json:"api_key"`
	BaseURL string     `json:"base_url"`
	Last    *time.Time `json:"last"`
	path    string     `json:"-"`
}

func ReadConfigFile(path string) (*ConfigurationData, error) {
	config := ConfigurationData{}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}
	config.path = path
	return &config, nil
}

func (c *ConfigurationData) UpdateLast(newTime time.Time) error {
	c.Last = &newTime
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.path, b, 0644)
	return err
}

func (c *ConfigurationData) GetLastTime() time.Time {
	if c.Last == nil {
		return time.Now().Add(-time.Hour * 1)
	}
	return *c.Last
}
