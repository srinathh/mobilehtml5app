package todoapp

import "sync"

type mapBackend struct {
	Items map[string]item
	sync.RWMutex
}

func (m *mapBackend) fetchAll() (map[string]item, error) {
	m.RLock()
	ret := make(map[string]item, len(m.Items))
	for k, v := range m.Items {
		ret[k] = v
	}
	m.RUnlock()
	return ret, nil
}

func (m *mapBackend) edit(id string, i item) error {
	m.Lock()
	m.Items[id] = i
	m.Unlock()
	return nil

}

func newMapBackend() *mapBackend {
	return &mapBackend{
		Items: make(map[string]item),
	}
}

func (m *mapBackend) exists(id string) bool {
	m.RLock()
	_, ok := m.Items[id]
	m.RUnlock()
	return ok
}
