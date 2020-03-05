package conf

type config struct {
	Env  string `yaml:"env"`
	Cors *cors  `yaml:"cors"`
	App  *app   `yaml:"app"`
}

type cors struct {
	AllowedOrigins string `yaml:"allowed-origins"`
	AllowedMethods string `yaml:"allow-methods"`
	AllowedHeaders string `yaml:"allow-headers"`
	MaxAge         string `yaml:"max-age"`
}

type app struct {
	Port string `yaml:"port"`
}

var (
	conf *config
)

func Env() string {
	return conf.Env
}

func Cors() *cors {
	return conf.Cors
}

func App() *app {
	return conf.App
}
