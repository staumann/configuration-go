package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	defaultConfigurationFileNamePattern = "%sapplication.yaml"
	profileConfigurationFilePattern     = "%sapplication-%s.yaml"
)

var (
	configMap map[string]interface{}
)

func init() {
	configMap = make(map[string]interface{})
}

//Init initializes the configuration
func Init(profile, configDir string, logLevel LogLevel) {

	setLogLevel(logLevel)
	logMessage(INFO, fmt.Sprintf("Profile: %s", profile))
	logMessage(INFO, fmt.Sprintf("Configuration directory: %s", configDir))

	//add a trailing slash to config directory
	if configDir != "" {
		logMessage(WARN, fmt.Sprintf("adding a trailing slash to configDir %s", configDir))
		configDir = configDir + "/"
	}

	readConfig(profile, configDir)
	logMessage(DEBUG, "processing of config finished")
}

func readConfig(profile string, configDir string) {
	newConfig := make(map[interface{}]interface{})
	newProfileConfig := make(map[interface{}]interface{})

	readFile(fmt.Sprintf(defaultConfigurationFileNamePattern, configDir), &newConfig)
	readFile(fmt.Sprintf(profileConfigurationFilePattern, configDir, profile), &newProfileConfig)

	processConfig("", newConfig, newProfileConfig)
}

func readFile(fileName string, pointer interface{}) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logMessage(WARN, fmt.Sprintf("could not open configuration file %s: %s", fileName, err.Error()))
		return
	}

	if e := yaml.Unmarshal(file, pointer); e != nil {
		logMessage(ERROR, fmt.Sprintf("could not unmarshal configuration file %s: %s", fileName, e.Error()))
	}
}

//AddMapToConfig adds the given map like a configfile
func AddMapToConfig(customCfg map[interface{}]interface{}) {
	processConfig("", customCfg, nil)
}

func processConfig(prefix string, general, profile map[interface{}]interface{}) {
	for key, value := range general {
		keyString := key.(string)
		cfgKey := prefix + keyString
		switch value.(type) {
		case map[interface{}]interface{}:
			pV := profile[key]
			if pV == nil {
				pV = make(map[interface{}]interface{})
			}
			processConfig(cfgKey+".", value.(map[interface{}]interface{}), pV.(map[interface{}]interface{}))
			break
		case string, int, bool:
			if profile != nil {
				if v, ok := profile[key]; ok {
					configMap[cfgKey] = v
					logMessage(DEBUG, fmt.Sprintf("new config value for key %s : %v", cfgKey, v))
					break
				}
			}
			configMap[cfgKey] = value
			logMessage(DEBUG, fmt.Sprintf("new config value for key %s : %v", cfgKey, value))
			break
		default:
			logMessage(ERROR, fmt.Sprintf("unregonized configuration type for value: %v", value))
		}
	}
}

//GetString get's a config string value from the configuration. if the key was not found we return "".
func GetString(key string) string {
	return GetStringWithDefault(key, "")
}

//GetStringWithDefault returns the value for the given key. If the key does not exist we return the defaultValue
func GetStringWithDefault(key, defaultValue string) string {
	return GetValueWithDefaultValue(key, defaultValue).(string)
}

//GetInteger get's a integer from the configuration
func GetInteger(key string) int {
	return configMap[key].(int)
}

//GetBoolean get's a bool from the configuration. if the key is not found false is returned
func GetBoolean(key string) bool {
	return GetBooleanWithDefaultValue(key, false)
}

//GetBooleanWithDefaultValue get's a bool from the configuration. if the key was not found the default value is returned
func GetBooleanWithDefaultValue(key string, defaultValue bool) bool {
	return GetValueWithDefaultValue(key, defaultValue).(bool)
}

//GetValue returns the generic value of the config map for this key. if the key does not exists nil is returned
func GetValue(key string) interface{} {
	return GetValueWithDefaultValue(key, nil)
}

//GetValueWithDefaultValue returns the generic value of the config map for this key it was not found it returns the default value
func GetValueWithDefaultValue(key string, defaultValue interface{}) interface{} {
	if value, exists := configMap[key]; exists {
		return value
	}
	return defaultValue
}
