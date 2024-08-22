package config

type Config struct {
	HttpConfig server    `mapstructure:"server"`
	DB         database  `mapstructure:"database"`
	LogLevel   logConfig `mapstructure:"log"`
	JwtConfig  jwtConfig `mapstructure:"auth"`
	S3Config   s3Config  `mapstructure:"s3"`
}
type database struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	DatabaseName string `mapstructure:"databaseName"`
}
type server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type logConfig struct {
	Level string `mapstructure:"level"`
}

type jwtConfig struct {
	SecretKey       string `mapstructure:"jwtSecretKey"`
	AccessTokenExp  int    `mapstructure:"expAT"`
	RefreshTokenExp int    `mapstructure:"expRT"`
}

type s3Config struct {
	Region   string `mapstructure:"AWS_DEFAULT_REGION"`
	Endpoint string `mapstructure:"ENDPOINT"`
	Bucket   string `mapstructure:"BUCKET"`
	AwsId    string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsKey   string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}
