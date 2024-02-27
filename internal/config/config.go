package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Configuration represents entire service configuration
type Configuration struct {
	Database    configuration.Database
	Debug       configuration.Debug
	Transcoding Transcoding
}

type Transcoding struct {
	Workers uint
	// Minutes
	MaxJobDuration uint `json:"max_job_duration"`
	Directory      string
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
