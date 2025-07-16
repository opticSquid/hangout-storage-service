package config

import (
	"flag"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func InitAppConfig(k *koanf.Koanf) {
	// Loading base config
	err := k.Load(file.Provider("resources/application.yaml"), yaml.Parser())
	if err != nil {
		log.Fatal().Msg("error loading base configuration")
	}
	// Loading specific config based on flag
	profile := flag.String("profile", "default", "sets a particular profile that overrides base configuration")
	flag.Parse()
	err = k.Load(file.Provider("resources/application-"+*profile+".yaml"), yaml.Parser())
	if err != nil {
		log.Error().Str("active profile", *profile).Msg("error loading profile configuration")
	}
	// Load environment variables and merge into the loaded config.
	// "HCDA" is the prefix to filter the env vars by.
	// "." is the delimiter used to represent the key hierarchy in env vars.
	// The (optional, or can be nil) function can be used to transform
	// the env var names, for instance, to lowercase them.
	//
	// For example, env vars: HSS_TYPE and HSS_PARENT1_CHILD1_NAME
	// will be merged into the "type" and the nested "parent1.child1.name"
	// keys in the config file here as we lowercase the key,
	// replace `_` with `.` and strip the HCDS_ prefix so that
	// only "parent1.child1.name" remains.
	k.Load(env.Provider("HSS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "HSS_")), "_", ".", -1)
	}), nil)
	log.Info().Str("active profile", *profile).Msg("configurations Loaded")
}
