package config

type Mode string

const (
	DebugMode   = Mode("debug")
	ReleaseMode = Mode("release")
)

type Config struct {
	Mode                Mode   `json:"mode"                              validate:"required"`
	ServiceName         string `json:"service_name"                      validate:"required"`
	HttpPort            int    `json:"http_port"                         validate:"required"`
	MongodbURL          string `json:"mongodb_url"                       validate:"required"`
	MongodbDatabaseName string `json:"mongodb_database_name"                    validate:"required"`
}

var config *Config

func init() {
	config = &Config{}
}

func GetConfig() Config {
	return *config
}
