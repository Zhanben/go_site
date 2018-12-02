package config

import(
    "fmt"
    "github.com/spf13/viper"
)

func ParseConfig() error {
    viper.SetConfigName("config") // name of config file (without extension)
    viper.AddConfigPath("./conf") // call multiple times to add many search paths
    viper.AddConfigPath(".")    // optionally look for config in the working directory
    err := viper.ReadInConfig() // Find and read the config file
    if err != nil {             // Handle errors reading the config file
        return  fmt.Errorf("Fatal error config file: %s \n", err)
    }
    viper.WatchConfig()
    return nil
}