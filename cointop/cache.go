package cointop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (ct *Cointop) writeHardCache(data interface{}, filename string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/.%s", ct.hardCachePath(), filename)
	of, err := os.Create(path)
	defer of.Close()
	if err != nil {
		return err
	}
	_, err = of.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (ct *Cointop) readHardCache(i interface{}, filename string) (interface{}, bool, error) {
	path := fmt.Sprintf("%s/.%s", ct.hardCachePath(), filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, false, err
	}
	return i, true, nil
}

func (ct *Cointop) hardCachePath() string {
	return fmt.Sprintf("%v%v", ct.configDirPath(), "/.cache")
}

func (ct *Cointop) createCacheDir() error {
	path := ct.hardCachePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
