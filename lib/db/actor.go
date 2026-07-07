package db

import (
	"samovar/lib/change"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

type data = map[string]string

type DbActor struct {
	act.Actor
	name   string
	prefix string
	buffer data
}

func (da *DbActor) InitDb(prefix, schema string) {
  inter.RegisterActor(da, schema)
	da.buffer = make(data)
	da.prefix = prefix
	da.schedule()
}

func (da *DbActor) SetName(name string) {
	if name != "" {
		da.name = name
	}
}

func (da *DbActor) isReady() bool {
	return da.name != ""
}

func (da *DbActor) isBuffered() bool {
	return len(da.buffer) > 0
}

func (da *DbActor) SetKey(key, val string) {
	da.buffer[key] = val
}

func (da *DbActor) Persist() error {
	if !da.isReady() || !da.isBuffered() {
		return nil
	}
	currentData, err := da.Load()
	if err != nil {
		currentData = make(data)
	}
	for key, val := range da.buffer {
		if val == "" {
			delete(currentData, key)
		} else {
			currentData[key] = val
		}
	}
	return da.Save(currentData)
}

var cleanRe = regexp.MustCompile("[^a-zA-Z0-9_-]+")

func (da *DbActor) fileName() (string, error) {
	cleanName := cleanRe.ReplaceAllString(da.name, "")
	if cleanName == "" {
		return "", errors.New("DB name is required")
	} else {
		fileName := fmt.Sprintf("%s_%s.json", da.prefix, cleanName)
		return "./db/" + fileName, nil
	}
}

func (da *DbActor) Save(d data) error {
	fileName, err := da.fileName()
	if err != nil {
		return err
	}
	json, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, json, 0644)
}

func (da *DbActor) Load() (d data, err error) {
	fileName, err := da.fileName()
	if err != nil {
		return d, err
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}
	return d, err
}

func (da *DbActor) schedule() error {
	_, err := da.SendAfter(da.PID(), "persist", 10*time.Second)
	return err
}

func (da *DbActor) HandleMessage(_ gen.PID, msg any) error {
  da.Log().Debug("DbActor.HandleMessage.msg: %#v", msg)
	switch d := msg.(type) {
	case data:
		for key, val := range d {
			da.SetKey(key, val)
		}
	case change.Report:
		da.SetKey(d.Key, d.NewValue.String())
	case string:
		if err := da.Persist(); err != nil {
			return err
		}
		return da.schedule()
	default:
		return errors.New("Unexpected db message")
	}
	return nil
}

func (da *DbActor) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
  da.Log().Debug("DbActor.HandleCall.req: %#v", req)
	keys, ok := req.([]string)
	if !ok {
		return nil, errors.New("Unexpected db request")
	}
	d, err := da.Load()
	if err != nil {
		return nil, err
	}
	result := make(data, len(keys))
	for _, key := range keys {
		if val, exists := d[key]; exists {
			result[key] = val
		}
	}
	return result, nil
}

func (da *DbActor) Terminate(_ error) {
	if err := da.Persist(); err != nil {
		da.Log().Error("Failed to persist DB[%s_%s]: %s", da.prefix, da.name, err)
	} else {
		da.Log().Info("DB[%s_%s] persisted", da.prefix, da.name)
	}
}
