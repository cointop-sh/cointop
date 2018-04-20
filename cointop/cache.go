package cointop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (ct *Cointop) writeCache(data []byte) error {
	path := ct.cachePath()
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ct *Cointop) readCache(i interface{}) (interface{}, bool, error) {
	path := ct.cachePath()
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

func (ct *Cointop) cachePath() string {
	return fmt.Sprintf("%v%v", ct.configDirPath(), "/.cache")
}
