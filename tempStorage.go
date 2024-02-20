package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

type TempStorage struct {
	path string
}

func (ts *TempStorage) cleanRoutine() {
	for {
		entries, err := os.ReadDir(ts.path)
		if err != nil {
			ErrorLog.Println(err.Error())
			time.Sleep(time.Second)
			continue
		}
		deadline := time.Now().AddDate(0, 0, -1)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				ErrorLog.Println(err.Error())
				continue
			}
			if info.ModTime().Before(deadline) {
				if os.Remove(path.Join(ts.path, entry.Name())); err != nil {
					ErrorLog.Println(err.Error())
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

func (ts *TempStorage) SaveBencodeFromForwarder(p []byte, hash string, uri string) string {
	f, err := os.CreateTemp(ts.path, fmt.Sprintf("%s_", hash))
	if err != nil {
		ErrorLog.Println(err.Error())
		return ``
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "%s\n%s\n", hash, uri); err != nil {
		ErrorLog.Println(err.Error())
		return f.Name()
	}
	if _, err := f.Write(p); err != nil {
		ErrorLog.Println(err.Error())
	}
	return f.Name()
}

func NewTempStorage(_path string) (*TempStorage, error) {
	ts := TempStorage{
		path: _path,
	}
	if ts.path == `` {
		ts.path = path.Join(os.TempDir(), `retracker`)
	}
	if err := os.MkdirAll(ts.path, 0755); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	go ts.cleanRoutine()
	return &ts, nil
}
