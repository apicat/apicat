package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStorePersistenceAndConflict(t *testing.T) {
	store := New(t.TempDir())
	state, err := store.Load()
	if err != nil || state.Revision != 0 {
		t.Fatalf("initial: %+v %v", state, err)
	}
	state.Environments = []Environment{{ID: "test", Name: "测试", BaseURL: "http://localhost", Variables: []Variable{{Key: "token", Value: "secret", Enabled: true, Secret: true}}}}
	state.ActiveEnvironment = "test"
	saved, err := store.Save(state)
	if err != nil {
		t.Fatal(err)
	}
	reloaded, err := New(filepath.Dir(filepath.Dir(store.Path()))).Load()
	if err != nil || reloaded.Revision != 1 || reloaded.Environments[0].Variables[0].Value != "secret" {
		t.Fatalf("reload: %+v %v", reloaded, err)
	}
	if _, err := store.Save(state); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	info, err := os.Stat(store.Path())
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0077 != 0 {
		t.Fatal("workspace readable by other users")
	}
	saved.Environments = append(saved.Environments, saved.Environments[0])
	if _, err = store.Save(saved); err == nil {
		t.Fatal("accepted duplicate environment")
	}
}
