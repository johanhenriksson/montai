package mount

import (
	"fmt"
	"strings"
)

type ProviderMap map[string]Provider
type MountMap map[string]Mount

type Manager struct {
	Root      string
	Providers ProviderMap
	Mounts    MountMap
}

func NewManager(root string) *Manager {
	return &Manager{
		Root:      root,
		Providers: make(ProviderMap),
		Mounts:    make(MountMap),
	}
}

func (m *Manager) AddProvider(p Provider) {
	m.Providers[p.Type()] = p
}

func (m *Manager) Mount(req *Request) (Mount, error) {
	req.Name = strings.ToLower(req.Name)

	provider, ok := m.Providers[req.Type]
	if !ok {
		return nil, ErrUnsupportedProvider
	}

	req.Path = fmt.Sprintf("%s/%s", m.Root, req.Name)
	mnt, err := provider.Mount(req)
	if err != nil {
		return nil, err
	}

	m.Mounts[req.Name] = mnt
	return mnt, nil
}

func (m *Manager) Unmount(name string) (Mount, error) {
	name = strings.ToLower(name)
	mnt, exists := m.Mounts[name]
	if !exists {
		return nil, ErrMountNotFound
	}

	mnt.Unmount()
	delete(m.Mounts, name)
	return mnt, nil
}
