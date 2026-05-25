package conf

import (
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var Dictinary *map[any]any

func LoadLocales(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	m := make(map[any]any)
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}
	Dictinary = &m
	return nil
}

func T(key string) string {
	dic := *Dictinary
	keys := strings.Split(key, ".")
	for index, path := range keys {
		if len(keys) == (index + 1) {
			for k, v := range dic {
				if k, ok := k.(string); ok {
					if k == path {
						if value, ok := v.(string); ok {
							return value
						}
					}
				}
			}
			return path
		}
		for k, v := range dic {
			if ks, ok := k.(string); ok {
				if ks == path {
					if dic, ok = v.(map[any]any); ok == false {
						return path
					}
				}
			} else {
				return ""
			}
		}
	}
	return ""
}
