package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config") // Name of the configuration file (without extension)
	viper.AddConfigPath(".")      // Path to look for the configuration file

	err := viper.ReadInConfig() // Find and read the configuration file
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}

	// Access configuration values
	fmt.Println("Database Host:", viper.GetString("database.host"))
	fmt.Println("Port:", viper.GetInt("server.port"))
}
