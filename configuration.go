package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const (
	defaultConfigurationFileNamePattern = "%sapplication.yaml"
	profileConfigurationFilePattern     = "%sapplication-%s.yaml"
)

var (
	defaultConfig ConfigurationObject
)

func init() {
	defaultConfig = ConfigurationObject{
		configMap: make(map[string]interface{}),
	}
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
	newConfig := make(map[string]interface{})
	newProfileConfig := make(map[string]interface{})

	readFile(fmt.Sprintf(defaultConfigurationFileNamePattern, configDir), &newConfig)
	readFile(fmt.Sprintf(profileConfigurationFilePattern, configDir, profile), &newProfileConfig)

	defaultConfig.processConfig("", newConfig, newProfileConfig)
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
func (cfg ConfigurationObject) AddMapToConfig(prefix string, customCfg map[string]interface{}) {
	//if there is no trailing dot add one
	if !strings.HasSuffix(".", prefix) {
		prefix = prefix + "."
	}
	cfg.processConfig(prefix, customCfg, nil)
}

func (cfg ConfigurationObject) processConfig(prefix string, general, profile map[string]interface{}) {
	for key, value := range general {
		keyString := key
		cfgKey := prefix + keyString
		switch value.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			pV := profile[key]
			if pV == nil {
				pV = make(map[interface{}]interface{})
			}
			cfg.processConfig(cfgKey+".", checkAndConvertMap(value), checkAndConvertMap(pV))
			break
		case string, int, bool:
			if profile != nil {
				if v, ok := profile[key]; ok {
					cfg.configMap[cfgKey] = v
					logMessage(DEBUG, fmt.Sprintf("new profile config value for key %s : %v", cfgKey, v))
					break
				}
			}
			cfg.configMap[cfgKey] = value
			logMessage(DEBUG, fmt.Sprintf("new global config value for key %s : %v", cfgKey, value))
			break
		default:
			logMessage(ERROR, fmt.Sprintf("unregonized configuration type for value: %v", value))
		}
	}
}

func checkAndConvertMap(source interface{}) map[string]interface{} {
	var target map[string]interface{}
	switch source.(type) {
	case map[string]interface{}:
		target = source.(map[string]interface{})
		break
	case map[interface{}]interface{}:
		target = make(map[string]interface{})
		for k, v := range source.(map[interface{}]interface{}) {
			target[k.(string)] = v
		}
		break
	default:
		logMessage(WARN, "passed interface is no valid map")
	}
	return target
}

//GetString get's a config string value from the configuration. if the key was not found we return "".
func (cfg ConfigurationObject) GetString(key string) string {
	return cfg.GetStringWithDefault(key, "")
}

//GetStringWithDefault returns the value for the given key. If the key does not exist we return the defaultValue
func (cfg ConfigurationObject) GetStringWithDefault(key, defaultValue string) string {
	return cfg.GetValueWithDefaultValue(key, defaultValue).(string)
}

//GetInteger get's a integer from the configuration. if the key does not exist 0 is returned
func (cfg ConfigurationObject) GetInteger(key string) int {
	return cfg.GetIntegerWithDefaultValue(key, 0)
}

//GetInteger get's a integer from the configuration. if the key does not exists the default value is returned
func (cfg ConfigurationObject) GetIntegerWithDefaultValue(key string, defaultValue int) int {
	return cfg.GetValueWithDefaultValue(key, defaultValue).(int)
}

//GetBoolean get's a bool from the configuration. if the key is not found false is returned
func (cfg ConfigurationObject) GetBoolean(key string) bool {
	return cfg.GetBooleanWithDefaultValue(key, false)
}

//GetBooleanWithDefaultValue get's a bool from the configuration. if the key was not found the default value is returned
func (cfg ConfigurationObject) GetBooleanWithDefaultValue(key string, defaultValue bool) bool {
	return cfg.GetValueWithDefaultValue(key, defaultValue).(bool)
}

//GetValue returns the generic value of the config map for this key. if the key does not exists nil is returned
func (cfg ConfigurationObject) GetValue(key string) interface{} {
	return cfg.GetValueWithDefaultValue(key, nil)
}

//GetValueWithDefaultValue returns the generic value of the config map for this key it was not found it returns the default value
func (cfg ConfigurationObject) GetValueWithDefaultValue(key string, defaultValue interface{}) interface{} {
	if value, exists := cfg.configMap[key]; exists {
		return value
	}
	return defaultValue
}

//GetSubConfig returns a subConfig object for the given key
func (cfg ConfigurationObject) GetSubConfig(key string) ConfigurationObject {
	if !strings.HasSuffix(key, ".") {
		key = key + "."
	}
	subConfig := ConfigurationObject{configMap: map[string]interface{}{}}
	for k, v := range cfg.configMap {
		if strings.HasPrefix(k, key) {
			subKey := strings.ReplaceAll(k, key, "")
			subConfig.configMap[subKey] = v
		}
	}
	return subConfig
}

//GetString get's a config string value from the configuration. if the key was not found we return "".
func GetString(key string) string {
	return defaultConfig.GetString(key)
}

//GetStringWithDefault returns the value for the given key. If the key does not exist we return the defaultValue
func GetStringWithDefault(key, defaultValue string) string {
	return defaultConfig.GetStringWithDefault(key, defaultValue)
}

//GetInteger get's a integer from the configuration. if the key does not exist 0 is returned
func GetInteger(key string) int {
	return defaultConfig.GetInteger(key)
}

//GetIntegerWithDefaultValue returns the value for the given key. If the key does not exist we return the defaultValue
func GetIntegerWithDefaultValue(key string, defaultValue int) int {
	return defaultConfig.GetValueWithDefaultValue(key, defaultValue).(int)
}

//GetBoolean get's a bool from the configuration. if the key is not found false is returned
func GetBoolean(key string) bool {
	return defaultConfig.GetBoolean(key)
}

//GetBooleanWithDefaultValue get's a bool from the configuration. if the key was not found the default value is returned
func GetBooleanWithDefaultValue(key string, defaultValue bool) bool {
	return defaultConfig.GetBooleanWithDefaultValue(key, defaultValue)
}

//GetValue returns the generic value of the config map for this key. if the key does not exists nil is returned
func GetValue(key string) interface{} {
	return defaultConfig.GetValue(key)
}

//GetValueWithDefaultValue returns the generic value of the config map for this key it was not found it returns the default value
func GetValueWithDefaultValue(key string, defaultValue interface{}) interface{} {
	return defaultConfig.GetValueWithDefaultValue(key, defaultValue)
}

//GetSubConfig returns a subConfig object for the given key
func GetSubConfig(key string) ConfigurationObject {
	return defaultConfig.GetSubConfig(key)
}

//AddMapToConfig adds the given map like a configfile
func AddMapToConfig(prefix string, customCfg map[string]interface{}) {
	defaultConfig.AddMapToConfig(prefix, customCfg)
}
