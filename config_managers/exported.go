package config

var (
	std = New()
)

func StandardConfig() *Config {
	return std
}

func LoadConfig(fileName string) {
	std.LoadConfig(fileName)
}

func GetValue(section, key string) string {
	return std.ConfigInfo.Section(section).Key(key).String()
}
