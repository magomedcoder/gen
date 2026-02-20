package runner

import (
	"sync"

	"github.com/magomedcoder/gen/api/pb/runnerpb"
)

type RunnerState struct {
	Address string
	Enabled bool
}

type Registry struct {
	mu      sync.RWMutex
	runners map[string]bool
}

func NewRegistry(initialAddresses []string) *Registry {
	runners := make(map[string]bool)
	for _, addr := range initialAddresses {
		if addr != "" {
			runners[addr] = true
		}
	}
	return &Registry{runners: runners}
}

func (r *Registry) Register(addr string) {
	if addr == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.runners[addr]; !ok {
		r.runners[addr] = true
	}
}

func (r *Registry) Unregister(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.runners, addr)
}

func (r *Registry) GetRunners() []*runnerpb.RunnerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*runnerpb.RunnerInfo, 0, len(r.runners))
	for addr, enabled := range r.runners {
		out = append(out, &runnerpb.RunnerInfo{
			Address: addr,
			Enabled: enabled,
		})
	}
	return out
}

func (r *Registry) SetEnabled(addr string, enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.runners[addr]; ok {
		r.runners[addr] = enabled
	}
}

func (r *Registry) HasActiveRunners() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, enabled := range r.runners {
		if enabled {
			return true
		}
	}
	return false
}

func (r *Registry) GetEnabledAddresses() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []string
	for addr, enabled := range r.runners {
		if enabled {
			out = append(out, addr)
		}
	}
	return out
}
