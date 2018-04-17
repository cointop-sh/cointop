package cointop

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Shortcuts map[string]interface{} `toml:"shortcuts"`
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
	return nil
}

func (ct *Cointop) configDirPath() string {
	usr, err := user.Current()
	if err != nil {
		return ".cointop"
	}

	return fmt.Sprintf("%s%s", usr.HomeDir, "/.cointop")
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
	s := defaultShortcuts()
	ifcs := map[string]interface{}{}
	for k, v := range s {
		var i interface{} = v
		ifcs[k] = i
	}
	var inputs = &config{
		Shortcuts: ifcs,
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
	actionsmap := actionsMap()
	for k, ifc := range ct.config.Shortcuts {
		k = strings.ToLower(k)
		v, ok := ifc.(string)
		if ok {
			if !actionsmap[v] {
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
