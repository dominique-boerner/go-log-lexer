package config

// Config
// Represents the structure of the 'config.toml'
type Config struct {
	Configuration struct {
		LogFile string `toml:"logFile"`
	} `toml:"Configuration"`
}
