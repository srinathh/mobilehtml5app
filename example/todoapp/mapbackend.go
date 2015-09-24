package todoapp

import (
	"fmt"
	"sync"
	"time"
)

type mapBackend struct {
	Items map[string]item
	sync.RWMutex
}

func (m *mapBackend) fetchAll() ([]item, error) {
	m.RLock()
	ret := make([]item, len(m.Items))
	ctr := 0
	for _, v := range m.Items {
		ret[ctr] = v
		ctr++
	}
	m.RUnlock()
	return ret, nil
}

func (m *mapBackend) create(i item) error {
	m.Lock()
	m.Items[i.ID] = i
	m.Unlock()
	return nil
}

func (m *mapBackend) delete(id string) error {
	m.Lock()
	if _, ok := m.Items[id]; !ok {
		return fmt.Errorf("item %s not found", id)
	}
	delete(m.Items, id)
	m.Unlock()
	return nil
}

func newMapBackend() *mapBackend {
	return &mapBackend{
		Items: make(map[string]item),
	}
}

func (m *mapBackend) createSample() {
	m.Lock()
	m.Items = make(map[string]item, 4)
	for j := 0; j < 4; j++ {
		t := time.Now()
		id := t.Format(timestamp)
		m.Items[id] = item{
			ID:       id,
			Text:     fmt.Sprintf("Test Item %d", j),
			Priority: j % 2,
			time:     t,
		}
		time.Sleep(time.Millisecond * 10)
	}
	m.Unlock()
}
