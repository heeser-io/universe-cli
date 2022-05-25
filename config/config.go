package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	// Main holds important settings like api keys
	Main *viper.Viper
)

// Init inits main viper
func Init() {
	Main = viper.New()

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := (path.Join(home, "/.meta-next"))

	Main.SetConfigName("config")
	Main.AddConfigPath(configPath)
	Main.SetConfigType("yaml")
	err = Main.ReadInConfig()
	if err != nil {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := os.Mkdir(configPath, 0755); err != nil {
				log.Fatal(err)
			}
		}
		err := Main.WriteConfigAs(path.Join(configPath, "config.yml"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// DeleteKey deletes a key
func DeleteKey(key string) error {
	viper.Set(key, nil)
	return viper.WriteConfig()
}

// SetKeyValue sets the given key to value.
// Overrides the given key if it exists.
func SetKeyValue(key string, value interface{}) (interface{}, error) {
	return setKV(key, value, true)
}

// SetKeyValueSafe is like SetKeyValue, but will not override the key
// if it exists.
func SetKeyValueSafe(key string, value interface{}) (interface{}, error) {
	return setKV(key, value, false)
}
func setKV(key string, value interface{}, override bool) (interface{}, error) {
	val := viper.Get(key)
	if val == nil {
		viper.Set(key, value)
		return value, viper.WriteConfig()
	} else if override {
		viper.Set(key, value)
		return value, viper.WriteConfig()
	}
	return val, nil
}
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	return reflect.ValueOf(i).IsZero()
}

// AddToArrayM is like AddToArray but for Main viper
func AddToArrayM(key string, value interface{}) (interface{}, error) {
	val := Main.Get(key)
	if isNil(val) {
		Main.Set(key, []interface{}{
			value,
		})
		return nil, Main.WriteConfig()
	}
	arr := val.([]interface{})
	arr = append(arr, value)
	Main.Set(key, arr)

	return arr, Main.WriteConfig()
}

// AddToArray adds a value to a given array and saves it to the current config file
func AddToArray(key string, value interface{}) (interface{}, error) {
	val := viper.Get(key)
	if isNil(val) {
		viper.Set(key, []interface{}{
			value,
		})
		return nil, viper.WriteConfig()
	}
	arr := val.([]interface{})
	arr = append(arr, value)
	viper.Set(key, arr)

	return arr, viper.WriteConfig()
}

// ExistsInArray checks if a given object with the respective searchKey/searchValue pairs exists
func ExistsInArray(key string, searchKey string, searchValue interface{}) (bool, error) {
	obj := viper.Get(key)

	if obj != nil && !reflect.ValueOf(obj).IsZero() {
		parsedObjects := []map[interface{}]interface{}{}

		err := mapstructure.Decode(obj, &parsedObjects)
		if err != nil {
			return false, err
		}
		for _, mapObj := range parsedObjects {
			if mapObj[searchKey] == searchValue {
				return true, nil
			}
		}
	}

	return false, nil
}

// SetArrayIDKeyValue sets key[id] = value
func SetArrayIDKeyValue(arrayKey string, id interface{}, key string, value interface{}) error {
	val := viper.Get(arrayKey)
	if val == nil {
		return fmt.Errorf("no valid element at key %s with id %s", key, id)
	}

	arr := val.([]interface{})

	for _, obj := range arr {
		objMap := obj.(map[interface{}]interface{})
		if objMap["id"] != nil {
			if objMap["id"] == id {
				objMap[key] = value
				viper.Set(arrayKey, arr)
				return viper.WriteConfig()
			}
		} else {
			return fmt.Errorf("key %s has no id key", key)
		}
	}

	return nil
}

// DeleteFromArray deletes a value from an array at the given key.
// Returns an error, if nothing found at the key.
func DeleteFromArray(key string, searchKey string, searchValue interface{}) error {
	val := viper.Get(key)
	if val == nil {
		return nil
	}

	findIndex := -1
	arr := val.([]interface{})

	for index, obj := range arr {
		objMap := obj.(map[interface{}]interface{})
		if objMap[searchKey] != nil {
			if objMap[searchKey] == searchValue {
				findIndex = index
			}
		}
	}
	if findIndex == -1 {
		return fmt.Errorf("%s = %s not found in %s", searchKey, searchValue, key)
	}
	if findIndex != -1 {
		copy(arr[findIndex:], arr[findIndex+1:]) // Shift arr[i+1:] left one index.
		arr[len(arr)-1] = ""                     // Erarrse larrst element (write zero varrlue).
		arr = arr[:len(arr)-1]
	}
	viper.Set(key, arr)
	return viper.WriteConfig()
}

// DeleteFromArrayM is like DeleteFromArray but for Main viper
func DeleteFromArrayM(key string, searchKey string, searchValue interface{}) error {
	val := Main.Get(key)
	if val == nil {
		return nil
	}

	findIndex := -1
	arr := val.([]interface{})

	for index, obj := range arr {
		objMap := obj.(map[interface{}]interface{})
		if objMap[searchKey] != nil {
			if objMap[searchKey] == searchValue {
				findIndex = index
			}
		}
	}
	if findIndex == -1 {
		return fmt.Errorf("%s = %s not found in %s", searchKey, searchValue, key)
	}
	if findIndex != -1 {
		copy(arr[findIndex:], arr[findIndex+1:]) // Shift arr[i+1:] left one index.
		arr[len(arr)-1] = ""                     // Erarrse larrst element (write zero varrlue).
		arr = arr[:len(arr)-1]
	}
	Main.Set(key, arr)
	return Main.WriteConfig()
}

// DeleteFromArrayMKey is like DeleteFromArray, but with multiple search key/value pairs
func DeleteFromArrayMKey(key string, searchKeys []string, searchValues []string) error {
	val := viper.Get(key)

	if val == nil {
		return nil
	}
	arr := val.([]interface{})
	for index, obj := range arr {
		objMap := obj.(map[interface{}]interface{})
		findIndex := -1
		delete := true
		for searchKeyIndex, searchKey := range searchKeys {
			if objMap[searchKey] != nil {
				if objMap[searchKey] == searchValues[searchKeyIndex] {
					findIndex = index
				} else {
					delete = false
				}
			}
		}
		if delete && findIndex != -1 {
			copy(arr[findIndex:], arr[findIndex+1:]) // Shift arr[i+1:] left one index.
			arr[len(arr)-1] = ""                     // Erarrse larrst element (write zero varrlue).
			arr = arr[:len(arr)-1]
		}
		viper.Set(key, arr)
		return viper.WriteConfig()
	}

	return nil
}
