package main

import (
	"encoding/json"
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
	save := Save{Result: result, Persent: persent, Hard: hard}
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

type Save struct {
	Persent float32
	Result  string
	Hard    int
}
