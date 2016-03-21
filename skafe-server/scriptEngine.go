package main

import (
	"fmt"
)

const (
	MAX_SCRIPT_WORKERS int = 5
)

const (
	SCRIPT_LANG_RUBY string = "ruby"
)

var (
	ScriptLangs []string = []string{
		SCRIPT_LANG_RUBY,
	}
)

type ScriptPool struct {
	execCmd map[string]string
	next    map[string]chan ScriptWorker
	count   map[string]int
	create  map[string]func(string) (ScriptWorker, error)
}

type ScriptWorker interface {
	Run(string, *AuditEvent) (bool, error)
	Lang() string
	HasFunc(string) bool
	Running() bool
}

func SetupScriptPool(conf *ServerConfig) (*ScriptPool, error) {
	pool := &ScriptPool{}

	// setup exec commands
	pool.execCmd = map[string]string{
		SCRIPT_LANG_RUBY: conf.execRuby,
	}

	// setup create functions
	pool.create = map[string]func(string) (ScriptWorker, error){
		SCRIPT_LANG_RUBY: NewRbScriptWorker,
	}

	pool.next = make(map[string]chan ScriptWorker)
	pool.count = make(map[string]int)

	for _, lang := range ScriptLangs {
		// setup worker channels
		pool.next[lang] = make(chan ScriptWorker, MAX_SCRIPT_WORKERS)

		// setup worker counts
		pool.count[lang] = 0
	}

	return pool, nil
}

// Get a script worker of the specified language
func (s *ScriptPool) GetWorker(lang string) (ScriptWorker, error) {
	if _, ok := s.count[lang]; !ok {
		return nil, fmt.Errorf("No such worker language: %s", lang)
	}

	if s.count[lang] == MAX_SCRIPT_WORKERS {
		// not authorized to create new workers, just have to wait
		return <-s.next[lang], nil
	}

	// otherwise we are allowed to create new ones if needbe
	select {
	case worker := <-s.next[lang]:
		return worker, nil
	default:
		s.count[lang] += 1
		return s.create[lang](s.execCmd[lang])
	}
}

// Return a script worker when it's use is finished
func (s *ScriptPool) ReturnWorker(worker ScriptWorker) {
	s.next[worker.Lang()] <- worker
}
