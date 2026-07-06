package watcher

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch monitors dir recursively for changes to spec files (yaml/yml/json)
// and invokes onChange, debounced, until the process exits. Subdirectories
// created after start are picked up automatically.
func Watch(dir string, debounce time.Duration, onChange func()) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	addDirs := func(root string) {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") && path != root {
					return filepath.SkipDir
				}
				if err := w.Add(path); err != nil {
					log.Printf("watch: cannot watch %s: %v", path, err)
				}
			}
			return nil
		})
	}
	addDirs(dir)

	var timer *time.Timer
	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if ev.Op.Has(fsnotify.Create) {
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					addDirs(ev.Name)
					continue
				}
			}
			if !isSpecFile(ev.Name) {
				continue
			}
			if !ev.Op.Has(fsnotify.Write) && !ev.Op.Has(fsnotify.Create) &&
				!ev.Op.Has(fsnotify.Remove) && !ev.Op.Has(fsnotify.Rename) {
				continue
			}
			if timer == nil {
				timer = time.AfterFunc(debounce, onChange)
			} else {
				timer.Reset(debounce)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			log.Printf("watch error: %v", err)
		}
	}
}

func isSpecFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml", ".json":
		return true
	}
	return false
}
