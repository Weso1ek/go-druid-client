package context

type Config struct {
	DruidHost       string `env:"DRUID_HOST" default:"druid-host"`
	DruidPort       string `env:"DRUID_PORT"`
	DruidDatasource string `env:"DRUID_DS"`
}
