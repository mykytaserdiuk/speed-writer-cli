package save

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func getOSSaveDir() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir, nil
	}

	// fall back to something sane on all platforms
	return os.UserCacheDir()
}

func SaveInResult(result string, hard int, persent float32) error {
	dir, err := getOSSaveDir()
	if err != nil {
		return err
	}
	currentTime := time.Now()
	fileName := currentTime.Format("2006.02.01.15.04.05") + ".json"
	path := filepath.Join(dir, "speed-writer")
	save := Save{Result: result, Persent: persent, Hard: hard, Name: fileName}
	file, err := json.Marshal(&save)
	if err != nil {
		return err
	}
	if _, err := os.Open(path); err != nil {
		err = os.MkdirAll(path, os.ModeAppend)
		if err != nil {
			return nil
		}
	}
	err = os.WriteFile(filepath.Join(path, fileName), file, 0600)
	if err != nil {
		return err
	}

	return err
}

func LoadAllSaves() ([]*Save, error) {
	dir, err := getOSSaveDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "speed-writer")
	dirEnts, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	saves := make([]*Save, 0)
	for _, v := range dirEnts {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}

		filePath := filepath.Join(path, v.Name())
		if err != nil {
			return nil, err
		}
		file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModeAppend)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		data, err := ioutil.ReadFile(filePath)
		var save Save
		if err := json.Unmarshal(data, &save); err != nil {
			return nil, fmt.Errorf("unmarshal save error %s: %w", filePath, err)
		}

		saves = append(saves, &save)

	}
	return saves, nil
}

type Save struct {
	Name    string
	Persent float32
	Result  string
	Hard    int
}
