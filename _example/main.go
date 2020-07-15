package main

import (
	"log"

	"github.com/crazy-max/gonfig"
	"github.com/pkg/errors"
)

func main() {
	cfg := Config{
		Db: (&Db{}).GetDefaults(),
	}

	// Load from file(s)
	fileLoader := gonfig.NewFileLoader(gonfig.FileLoaderConfig{
		Filename: "/path/to/myapp.yml",
		Finder: gonfig.Finder{
			BasePaths:  []string{"/etc/myapp/myapp", "$XDG_CONFIG_HOME/myapp", "$HOME/.config/myapp", "./myapp"},
			Extensions: []string{"yaml", "yml"},
		},
	})
	if found, err := fileLoader.Load(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to decode configuration from file"))
	} else if !found {
		log.Println(errors.Wrap(err, "No configuration file found"))
	} else {
		log.Printf("Configuration loaded from file: %s", fileLoader.GetFilename())
	}

	// Load from flags
	flagsLoader := gonfig.NewFlagLoader(gonfig.FlagLoaderConfig{
		Args: []string{
			"--timezone=Europe/Paris",
			"--logLevel",
		},
	})
	if found, err := flagsLoader.Load(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to decode configuration from flags"))
	} else if !found {
		log.Println(errors.Wrap(err, "No flags found"))
	} else {
		log.Printf("Configuration loaded from flags")
	}

	// Load from environment variables
	envLoader := gonfig.NewEnvLoader(gonfig.EnvLoaderConfig{
		Prefix: "MYAPP_",
	})
	if found, err := envLoader.Load(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to decode configuration from environment variables"))
	} else if !found {
		log.Println("No MYAPP_* environment variables defined")
	} else {
		log.Printf("Configuration loaded from %d environment variables\n", len(envLoader.GetVars()))
	}
}
