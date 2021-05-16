// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Period           time.Duration `config:"period"`
	BaseURL          string        `config:"base_url"`
	Username         string        `config:"username"`
	Password         string        `config:"password"`
	MetricConfigPath string        `config:"metric_config_path"`
	MetricsConfig    *MetricsConfig
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}

// MetricsConfig : Struct to represent the config file
type MetricsConfig struct {
	Metrics struct {
		StorageSystems []struct {
			MetricID       int    `yaml:"ibm_spectrum_metric_id"`
			EventFieldName string `yaml:"event_field_name"`
		} `yaml:"storage_systems"`
		StorageSystemsAndVolumes []struct {
			MetricID       int    `yaml:"ibm_spectrum_metric_id"`
			EventFieldName string `yaml:"event_field_name"`
		} `yaml:"storage_systems_and_volumes"`
		Switches []struct {
			MetricID       int    `yaml:"ibm_spectrum_metric_id"`
			EventFieldName string `yaml:"event_field_name"`
		} `yaml:"switches"`
		Pools struct {
			Properties []struct {
				PropertyName   string `yaml:"property_name"`
				EventFieldName string `yaml:"event_field_name"`
			} `yaml:"properties"`
		} `yaml:"pools"`
	} `yaml:"metrics"`
}

// GetConf file from the given path
func (c *Config) GetMetricsConf() error {
	yamlFile, err := ioutil.ReadFile(c.MetricConfigPath)
	if err != nil {
		return err
	}

	c.MetricsConfig = &MetricsConfig{}
	err = yaml.Unmarshal(yamlFile, c.MetricsConfig)
	return err
}
