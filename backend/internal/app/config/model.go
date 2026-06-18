//nolint:lll // struct tags can get long and it's more readable to keep them in one line
package config

type SentryConfig struct {
	DSN                     string  `json:"dsn"                        mapstructure:"dsn"                        validate:"omitempty,url"`
	TraceSampleRate         float64 `json:"trace_sample_rate"          mapstructure:"trace_sample_rate"          validate:"omitempty,gte=0,lte=1"`
	ReplaySessionSampleRate float64 `json:"replay_session_sample_rate" mapstructure:"replay_session_sample_rate" validate:"omitempty,gte=0,lte=1"`
	ReplayErrorSampleRate   float64 `json:"replay_error_sample_rate"   mapstructure:"replay_error_sample_rate"   validate:"omitempty,gte=0,lte=1"`
	Environment             string  `json:"environment"                mapstructure:"environment"                validate:"required"`
	Version                 string  `json:"version"                    mapstructure:"version"                    validate:"required"`
}

type AppConfig struct {
	Version string `mapstructure:"version"`

	GinMode     string `mapstructure:"gin_mode" validate:"required,oneof=debug release test"`
	Environment string `mapstructure:"env"      validate:"required,oneof=development staging production test"`

	LogLevel  string `mapstructure:"log_level"  validate:"required,oneof=debug info warn error"`
	LogFormat string `mapstructure:"log_format" validate:"required,oneof=json text"`

	DbFilename         string                 `mapstructure:"db_filename"         validate:"required"`
	FrontendURL        string                 `mapstructure:"frontend_url"        validate:"required,http_url"`
	CorsAllowOrigins   CorsAllowOriginsConfig `mapstructure:"cors_allow_origins"  validate:"dive,required"`
	EnvironmentMessage string                 `mapstructure:"environment_message"`
	RedisURL           RedisURL               `mapstructure:"redis_url"           validate:"omitempty,url"`

	OtelEndpoint string `mapstructure:"otel_endpoint" validate:"omitempty"`
	Port         int    `mapstructure:"port"          validate:"required,gt=0,lte=65535"`
}

type ExporterConfig struct {
	Name        string            `mapstructure:"name"        validate:"required"`
	Type        string            `mapstructure:"type"        validate:"required,oneof=rss json atom"`
	Destination string            `mapstructure:"destination" validate:"required,oneof=s3 file"`
	Filename    string            `mapstructure:"filename"    validate:"required"`
	Options     map[string]string `mapstructure:"options"`
	Enabled     bool              `mapstructure:"enabled"`
}

type CorsAllowOriginsConfig []string

type S3ClientConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"     validate:"required"`
	SecretAccessKey string `mapstructure:"secret_access_key" validate:"required"`
	Region          string `mapstructure:"region"            validate:"required"`
	Endpoint        string `mapstructure:"endpoint"          validate:"required"`
	UsePathStyle    bool   `mapstructure:"use_path_style"`
}

type FeedConfig struct {
	FeedTitle       string `mapstructure:"feed_title"`
	FeedLink        string `mapstructure:"feed_link"`
	FeedDescription string `mapstructure:"feed_description"`
	AuthorName      string `mapstructure:"author_name"`
	AuthorEmail     string `mapstructure:"author_email"     validate:"omitempty,email"`
}

type DateFormatOptionsConfig map[string]any

type DateFormatConfig struct {
	Locale  string                  `json:"locale"  mapstructure:"locale"  validate:"required"`
	Options DateFormatOptionsConfig `json:"options" mapstructure:"options"`
}

type FormatConfig struct {
	Date DateFormatConfig `json:"date" mapstructure:"date"`
}

type AuthConfig struct {
	Type          string `json:"type"      mapstructure:"type"            validate:"required,oneof=oidc"`
	Name          string `json:"name"      mapstructure:"name"            validate:"required"`
	AuthorityURL  string `json:"authority" mapstructure:"authority"       validate:"required,url"`
	ClientID      string `json:"client_id" mapstructure:"client_id"       validate:"required"`
	SkipTLSVerify bool   `json:"-"         mapstructure:"skip_tls_verify"`
}

type Config struct {
	App      AppConfig        `mapstructure:"app"`
	Sentry   SentryConfig     `mapstructure:"sentry"`
	Format   FormatConfig     `mapstructure:"format"`
	Exporter []ExporterConfig `mapstructure:"exporter"`
	S3Client *S3ClientConfig  `mapstructure:"s3_client"`
	Feed     *FeedConfig      `mapstructure:"feed"`
	Auth     *AuthConfig      `mapstructure:"auth"            validate:"omitempty"`
}
