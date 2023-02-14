package config

import "github.com/spf13/viper"

type Config struct {
	Port  string `mapstructure:"PORT"`
	DBUrl string `mapstructure:"DB_URL"`
}

func LoadConfig() (Config, error) {
	var c Config
	viper.AddConfigPath("./pkg/config/envs")
	viper.AddConfigPath("/GoCommentSvc") //Для docker

	viper.SetConfigName("prod")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		viper.SetConfigName("dev")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		err = viper.ReadInConfig()

		if err != nil {
			return Config{}, err
		}
		err = viper.Unmarshal(&c)

		return c, nil
	}

	err = viper.Unmarshal(&c)

	return c, nil
}
