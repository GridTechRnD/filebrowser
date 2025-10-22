package groups


import (
	"sync"
)


type StorageBackend interface {
	Create( group *Group ) error
	Update( group *Group) error
	GetAllGroups() ([]*Group, error)
	DeleteByID( groupId int ) error
}

type Store interface {
	SaveGroup( g *Group ) error
	UpdateGroup( g *Group ) error
	GetAll() ([]*Group, error)
	Delete( groupId int ) error
}

type Storage struct {
	back 		StorageBackend
	updated 	map[uint]int64
	mux			sync.RWMutex
}

func NewStorage(back StorageBackend) *Storage {
	return &Storage{
		back:    back,
		updated: map[uint]int64{},
	}
}


func (s *Storage) SaveGroup( g *Group ) error {
	
	for i, r := range g.Rules {

		if r.Path[ len(r.Path) - 1 ] == '/' {

			g.Rules[i].Path = g.Rules[i].Path[ : len(g.Rules[i].Path) - 1]
		}
	}

	return s.back.Create( g )
}

func (s *Storage) UpdateGroup( g *Group ) error {

	for i, r := range g.Rules {

		if r.Path[ len(r.Path) - 1 ] == '/' {

			g.Rules[i].Path = g.Rules[i].Path[ : len(g.Rules[i].Path) - 1]
		}
	}

	return s.back.Update( g )
}

func (s *Storage) GetAll() ([]*Group, error) {

	return s.back.GetAllGroups()
}

func (s *Storage) Delete( groupId int ) error {

	return s.back.DeleteByID( groupId )
}