package conf

type config struct {
	Env        string      `yaml:"env"`
	Cors       *cors       `yaml:"cors"`
	App        *app        `yaml:"app"`
	DataSource *dataSource `yaml:"data-source"`
}

type cors struct {
	AllowedOrigins string `yaml:"allowed-origins"`
	AllowedMethods string `yaml:"allow-methods"`
	AllowedHeaders string `yaml:"allow-headers"`
	MaxAge         string `yaml:"max-age"`
}

type dataSource struct {
	URL string `yaml:"url"`
}

type app struct {
	Port string `yaml:"port"`
}

var (
	conf *config
)

func DataSource() *dataSource {
	return conf.DataSource
}

func Env() string {
	return conf.Env
}

func Cors() *cors {
	return conf.Cors
}

func App() *app {
	return conf.App
}
