package configs

import (
	"github.com/spf13/viper"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("http.port", 8080)

	v.SetDefault("log.level", int32(2))
	v.SetDefault("log.json", true)
}
