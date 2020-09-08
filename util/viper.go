package util

import (
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"strings"
)

func SetupViperHomeConfig(configFile, configName, extension, envPrefix string) {
	// https://github.com/spf13/viper/issues/390#issuecomment-511177615
	viper.SetConfigType(extension) // or viper.SetConfigType("YAML")

	homeDir, _ := homedir.Dir()
	viper.AddConfigPath(homeDir)

	// https://github.com/spf13/viper/issues/390#issuecomment-336464039
	// the name of the config file
	if string(configName[0]) != "." {
		configName = "." + configName
	}
	viper.SetConfigName(configName)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// the prefix of env vars to override the configs
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
}

/**
 * The path to the file used for configuration
 */
func GetConfigFilePath() string {
	return viper.ConfigFileUsed()
}

/**
 * Read the configuration and update the reference of a provided struct represented by the given config path
 */
func ReadConfig(configPath string, result interface{}) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil
	}
	log.Infof("Loading the config object '%s' from '%s'", configPath, configFilePath)

	// Looking for the key of configuration in the config
	log.Tracef("The keys in the config: %v", viper.AllKeys())
	containsKey := false
	for _, fullKey := range viper.AllKeys() {
		if funk.Contains(fullKey, configPath) {
			containsKey = true
			break
		}
	}

	// Verify if the key exists in the
	if !containsKey {
		log.Debugf("The config file does NOT contain the key '%s'", configPath)
		return nil
	}

	// Load the MAP from viper config, just have the map
	configMap := viper.Get(configPath)
	log.Tracef("The loaded map is %v", configMap)

	if err := fromMapToStruct(configPath, configMap, result); err != nil {
		return err
	}
	return nil
}

/**
 * Convert the map loaded from a config to a Struct provided
 * https://stackoverflow.com/questions/26744873/converting-map-to-struct/26746461
 */
func fromMapToStruct(configPath string, configMap interface{}, configStruct interface{}) error {
	// Use https://godoc.org/github.com/mitchellh/mapstructure to read properties from tagged structs
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &configStruct,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	err = decoder.Decode(configMap)
	if err != nil {
		log.Errorf("An error occurred while decoding key %s: %v", configPath, err)
		return err
	}

	log.Tracef("Successfully transformed the map to struct: %v", configStruct)
	return nil
}

func getConfigFilePath() (string, error) {
	err := viper.ReadInConfig()
	if err == nil {
		// viper.ConfigFileUsed() only returns properties after it is read (filePathConfigured)
		configFile := viper.ConfigFileUsed()

		// viper.ConfigFileUsed() only returns properties after it is read
		log.Debug("Using config file: ", configFile)
		return configFile, nil

	} else {
		log.Debugf("can't parse configuration: %s", err)
		return "", err
	}
}
