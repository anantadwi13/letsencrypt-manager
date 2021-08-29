package internal

type Config interface {
	PublicStaticPath() string
	ApiPort() int
}

type ConfigParam struct {
	PublicStaticPath string
	ApiPort          int
}

type config struct {
	publicStaticPath string
	apiPort          int
}

func NewConfig(configParam ConfigParam) Config {
	return &config{publicStaticPath: configParam.PublicStaticPath, apiPort: configParam.ApiPort}
}

func (c *config) PublicStaticPath() string {
	return c.publicStaticPath
}

func (c *config) ApiPort() int {
	return c.apiPort
}
