package main
import (
    "log"

    "github.com/spf13/viper"
  )
type Env struct {
  DBHost     string `mapstructure:"DB_HOST"`
  DBPort     int    `mapstructure:"DB_PORT"`
  DBUser     string `mapstructure:"DB_USER"`
  DBPassword string `mapstructure:"DB_PASSWORD"`
  DBName     string `mapstructure:"DB_NAME"`

  APIKey    string `mapstructure:"API_KEY"`
  SecretKey string `mapstructure:"SECRET_KEY"`

  Debug     bool   `mapstructure:"DEBUG"`
  LogLevel  string `mapstructure:"LOG_LEVEL"`

  SMTPHost  string `mapstructure:"SMTP_HOST"`
  SMTPPort  int    `mapstructure:"SMTP_PORT"`
  SMTPUser  string `mapstructure:"SMTP_USER"`
  SMTPPass  string `mapstructure:"SMTP_PASSWORD"`

  RedisURL     string `mapstructure:"REDIS_URL"`
  CacheTimeout int    `mapstructure:"CACHE_TIMEOUT"`
}
func NewEnv() *Env {
  env := Env{}
  viper.SetConfigFile(".env")
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal("Unable to load environment file :", err)
  }
  err = viper.Unmarshal(&env)
  if err != nil {
    log.Fatal("Unable to decode into struct :", err)
  }
  if env.AppEnv == "development" {
    log.Println("The App is running in development env")
  }
  return &env
}
