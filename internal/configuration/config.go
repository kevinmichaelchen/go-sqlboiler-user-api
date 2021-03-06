package configuration

import (
	"encoding/json"

	"github.com/kevinmichaelchen/go-sqlboiler-user-api/internal/obs"
	"github.com/rs/zerolog/log"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

const (
	flagForGrpcPort = "grpc_port"
	flagForHTTPPort = "http_port"
)

type Config struct {
	// GrpcPort controls what port our gRPC server runs on.
	GrpcPort int

	// HTTPPort controls what port our HTTP server runs on.
	HTTPPort int

	// SQLConfig is the configuration for SQL database connection.
	SQLConfig SQLConfig

	// RedisConfig is the configuration for Redis connection.
	RedisConfig RedisConfig

	// LoggingConfig is the configuration for logging.
	LoggingConfig obs.LoggingConfig

	// TraceConfig contains config info for how we do tracing.
	TraceConfig obs.TraceConfig
}

func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not marshal config to string")
	}
	return string(b)
}

func LoadConfig() Config {
	c := Config{
		GrpcPort: 8084,
		HTTPPort: 8085,
	}

	c.SQLConfig = LoadSQLConfig()
	c.RedisConfig = LoadRedisConfig()
	c.TraceConfig = obs.LoadTraceConfig()
	c.LoggingConfig = obs.LoadLoggingConfig()

	flag.Int(flagForGrpcPort, c.GrpcPort, "gRPC port")
	flag.Int(flagForHTTPPort, c.HTTPPort, "HTTP port")

	flag.Parse()

	viper.BindPFlag(flagForGrpcPort, flag.Lookup(flagForGrpcPort))
	viper.BindPFlag(flagForHTTPPort, flag.Lookup(flagForHTTPPort))

	viper.AutomaticEnv()

	c.GrpcPort = viper.GetInt(flagForGrpcPort)
	c.HTTPPort = viper.GetInt(flagForHTTPPort)

	return c
}
