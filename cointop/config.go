package cointop

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var fileperm = os.FileMode(0644)

type config struct {
	Shortcuts map[string]interface{}   `toml:"shortcuts"`
	Favorites map[string][]interface{} `toml:"favorites"`
	Currency  interface{}              `toml:"currency"`
}

func (ct *Cointop) setupConfig() error {
	err := ct.makeConfigDir()
	if err != nil {
		return err
	}
	err = ct.makeConfigFile()
	if err != nil {
		return err
	}
	err = ct.parseConfig()
	if err != nil {
		return err
	}
	err = ct.loadShortcutsFromConfig()
	if err != nil {
		return err
	}
	err = ct.loadFavoritesFromConfig()
	if err != nil {
		return err
	}
	err = ct.loadCurrencyFromConfig()
	if err != nil {
		return err
	}
	return nil
}

func (ct *Cointop) loadFavoritesFromConfig() error {
	for k, arr := range ct.config.Favorites {
		if k == "symbols" {
			for _, ifc := range arr {
				v, ok := ifc.(string)
				if ok {
					ct.favorites[strings.ToUpper(v)] = true
				}
			}
		}
	}
	return nil
}

func (ct *Cointop) configDirPath() string {
	homedir := userHomeDir()
	return fmt.Sprintf("%s%s", homedir, "/.cointop")
}

func (ct *Cointop) configPath() string {
	return fmt.Sprintf("%v%v", ct.configDirPath(), "/config")
}

func (ct *Cointop) makeConfigDir() error {
	path := ct.configDirPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func (ct *Cointop) makeConfigFile() error {
	path := ct.configPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fo, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fo.Close()
		b, err := ct.configToToml()
		if err != nil {
			return err
		}
		if _, err := fo.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func (ct *Cointop) saveConfig() error {
	ct.savemux.Lock()
	defer ct.savemux.Unlock()
	path := ct.configPath()
	if _, err := os.Stat(path); err == nil {
		b, err := ct.configToToml()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, b, fileperm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ct *Cointop) parseConfig() error {
	var conf config
	path := ct.configPath()
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return err
	}

	ct.config = conf
	return nil
}

func (ct *Cointop) configToToml() ([]byte, error) {
	shortcutsIfcs := map[string]interface{}{}
	for k, v := range ct.shortcutkeys {
		var i interface{} = v
		shortcutsIfcs[k] = i
	}

	var favorites []interface{}
	for k, ok := range ct.favorites {
		if ok {
			var i interface{} = strings.ToUpper(k)
			favorites = append(favorites, i)
		}
	}
	favoritesIfcs := map[string][]interface{}{
		"symbols": favorites,
	}

	var currencyIfc interface{} = ct.currencyconversion

	var inputs = &config{
		Shortcuts: shortcutsIfcs,
		Favorites: favoritesIfcs,
		Currency:  currencyIfc,
	}

	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	err := encoder.Encode(inputs)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (ct *Cointop) loadShortcutsFromConfig() error {
	for k, ifc := range ct.config.Shortcuts {
		v, ok := ifc.(string)
		if ok {
			if !ct.actionExists(v) {
				continue
			}
			if ct.shortcutkeys[k] == "" {
				continue
			}
			ct.shortcutkeys[k] = v
		}
	}
	return nil
}

func (ct *Cointop) loadCurrencyFromConfig() error {
	if currency, ok := ct.config.Currency.(string); ok {
		ct.currencyconversion = strings.ToUpper(currency)
	}
	return nil
}
