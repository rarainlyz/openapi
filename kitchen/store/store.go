package store

import "errors"

type Menu struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Type  string  `json:"type"`
}

var ErrNotFound = errors.New("menu not found")

type MenuStore interface {
    List(menuType string) []Menu
    Get(id int) (Menu, error)
    Add(m Menu) Menu
    Delete(id int) error
}

type MemoryStore struct {
    menus  []Menu
    nextID int
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        menus: []Menu{
            {ID: 1, Name: "ต้มยำกุ้ง", Price: 120, Type: "soup"},
            {ID: 2, Name: "พิซซ่าฮาวายเอี้ยน", Price: 199, Type: "pizza"},
            {ID: 3, Name: "พิซซ่าเห็ด", Price: 179, Type: "pizza"},
        },
        nextID: 4,
    }
}

func (s *MemoryStore) List(menuType string) []Menu {
    if menuType == "" {
        return s.menus
    }
    filtered := []Menu{}
    for _, m := range s.menus {
        if m.Type == menuType {
            filtered = append(filtered, m)
        }
    }
    return filtered
}

func (s *MemoryStore) Get(id int) (Menu, error) {
    for _, m := range s.menus {
        if m.ID == id {
            return m, nil
        }
    }
    return Menu{}, ErrNotFound
}

func (s *MemoryStore) Add(m Menu) Menu {
    m.ID = s.nextID
    s.nextID++
    s.menus = append(s.menus, m)
    return m
}

func (s *MemoryStore) Delete(id int) error {
    for i, m := range s.menus {
        if m.ID == id {
            s.menus = append(s.menus[:i], s.menus[i+1:]...)
            return nil
        }
    }
    return ErrNotFound
}
