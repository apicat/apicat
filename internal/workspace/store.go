// Package workspace stores local environments independently of the OpenAPI spec.
package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/apicat/apicat/v3/internal/problem"
)

type Variable struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
	Secret  bool   `json:"secret"`
}
type Environment struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	BaseURL   string     `json:"baseUrl"`
	Variables []Variable `json:"variables"`
}
type State struct {
	Revision          int64         `json:"revision"`
	ActiveEnvironment string        `json:"activeEnvironment"`
	Variables         []Variable    `json:"variables"`
	Environments      []Environment `json:"environments"`
}

var ErrConflict = problem.New("配置已被其他窗口更新，请重新加载后再保存", nil)

type Store struct {
	mu   sync.Mutex
	path string
}

func New(specDir string) *Store {
	return &Store{path: filepath.Join(specDir, ".apicat", "workspace.json")}
}
func (s *Store) Path() string         { return s.path }
func (s *Store) Load() (State, error) { s.mu.Lock(); defer s.mu.Unlock(); return s.load() }
func (s *Store) load() (State, error) {
	state := State{Variables: []Variable{}, Environments: []Environment{}}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, fmt.Errorf("读取环境配置失败: %w", err)
	}
	normalize(&state)
	return state, Validate(state)
}
func Validate(state State) error {
	if len(state.Environments) > 100 {
		return problem.New("最多支持 100 个环境", nil)
	}
	ids := map[string]bool{}
	for _, env := range state.Environments {
		if strings.TrimSpace(env.ID) == "" || strings.TrimSpace(env.Name) == "" {
			return problem.New("环境 ID 和名称不能为空", nil)
		}
		if ids[env.ID] {
			return problem.New("环境 ID 不能重复", nil)
		}
		ids[env.ID] = true
		if err := validateVariables(env.Variables); err != nil {
			return fmt.Errorf("%s: %w", env.Name, err)
		}
	}
	if state.ActiveEnvironment != "" && !ids[state.ActiveEnvironment] {
		return problem.New("当前环境不存在", nil)
	}
	return validateVariables(state.Variables)
}
func validateVariables(variables []Variable) error {
	if len(variables) > 500 {
		return problem.New("最多支持 500 个变量", nil)
	}
	seen := map[string]bool{}
	for _, v := range variables {
		if !v.Enabled {
			continue
		}
		if strings.TrimSpace(v.Key) == "" || strings.ContainsAny(v.Key, "{}\r\n") {
			return problem.New("变量名不能为空或包含花括号、换行", nil)
		}
		if seen[v.Key] {
			return problem.New("变量 {name} 重复", map[string]any{"name": v.Key})
		}
		seen[v.Key] = true
	}
	return nil
}
func (s *Store) Save(state State) (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	normalize(&state)
	if err := Validate(state); err != nil {
		return state, err
	}
	current, err := s.load()
	if err != nil {
		return state, err
	}
	if current.Revision != state.Revision {
		return state, ErrConflict
	}
	state.Revision++
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return state, err
	}
	if err = os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return state, err
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), "workspace-*.tmp")
	if err != nil {
		return state, err
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err = f.Write(append(data, '\n')); err != nil {
		f.Close()
		return state, err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return state, err
	}
	if err = f.Close(); err != nil {
		return state, err
	}
	if err = os.Rename(name, s.path); err != nil {
		return state, err
	}
	return state, nil
}

func normalize(state *State) {
	if state.Variables == nil {
		state.Variables = []Variable{}
	}
	if state.Environments == nil {
		state.Environments = []Environment{}
	}
	for i := range state.Environments {
		if state.Environments[i].Variables == nil {
			state.Environments[i].Variables = []Variable{}
		}
	}
}
