package conf

type config struct {
	Env        string      `yaml:"env"`
	App        *app        `yaml:"app"`
	DataSource *dataSource `yaml:"data-source"`

	ReCaptchaSecret string `yaml:"re-captcha-secret"`
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

func ReCaptchaSecret() string {
	return conf.ReCaptchaSecret
}

func App() *app {
	return conf.App
}
