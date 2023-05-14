package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
)

type config struct {
	APIKey string `mapstructure:"sec_api_key"`
}

var once sync.Once
var cfg config

//go:embed local.yml
var defaultLocal []byte

func GetConfig() config {
	once.Do(func() {
		viper.SetConfigType("yml")
		// これでaws上の環境変数を取得しているっぽい
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		def := defaultLocal

		// 設定ファイルを読み込みます
		err := viper.ReadConfig(bytes.NewBuffer(def))
		if err != nil {
			fmt.Println("Failed to read yml file:", err)
			os.Exit(1) //プログラムを終了する関数
		}

		cfg = config{}
		if err := viper.Unmarshal(&cfg); err != nil {
			panic(err)
		}
	})
	return cfg

}
