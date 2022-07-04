package server

import (
	"errors"
	"fmt"
	"sync"
)

// Item represents stored item
type Item struct {
	K string
	V string
}

// Storage defines interface for ordered storage
type Storage interface {
	// AddItem adds Item to storage
	AddItem(item Item)
	// RemoveItem removes Item from storage
	RemoveItem(key string) error
	// GetItem returns Item with given id from storage
	GetItem(key string) (*Item, error)
	// GetAllItems returns items in storage, new slice created every time
	GetAllItems() []Item
	// Iterate allows to iterato over ordered in storage, can be used for processing which does not involve blocking IO
	Iterate(accept func(Item))
}

type rwLockedStorage struct {
	rwLock  sync.RWMutex
	storage Storage
}

// NewRWLockedStorage returns storage which adds read/write mutex to upstream storage
func NewRWLockedStorage(storage Storage) *rwLockedStorage {
	return &rwLockedStorage{
		rwLock:  sync.RWMutex{},
		storage: storage,
	}
}

func (s *rwLockedStorage) AddItem(item Item) {
	s.rwLock.Lock()
	s.storage.AddItem(item)
	s.rwLock.Unlock()
}

func (s *rwLockedStorage) RemoveItem(key string) error {
	s.rwLock.Lock()
	err := s.storage.RemoveItem(key)
	s.rwLock.Unlock()
	return err
}

func (s *rwLockedStorage) GetItem(key string) (*Item, error) {
	s.rwLock.RLock()
	item, err := s.storage.GetItem(key)
	s.rwLock.RUnlock()
	return item, err
}

func (s *rwLockedStorage) GetAllItems() []Item {
	s.rwLock.RLock()
	items := s.storage.GetAllItems()
	s.rwLock.RUnlock()
	return items
}

func (s *rwLockedStorage) Iterate(accept func(Item)) {
	s.rwLock.RLock()
	s.storage.Iterate(accept)
	s.rwLock.RUnlock()
}

type entry struct {
	prev *entry
	next *entry
	item *Item
}

type memoryStorage struct {
	head    *entry
	indexed map[string]*entry
}

// NewMemoryStorage returns storage backed by slice, Item id is an index in slice
func NewMemoryStorage() *memoryStorage {
	head := &entry{}
	head.next = head
	head.prev = head

	return &memoryStorage{
		head:    head,
		indexed: make(map[string]*entry),
	}
}

func (s *memoryStorage) AddItem(item Item) {
	if _, ok := s.indexed[item.K]; ok {
		// we know that index exists
		_ = s.RemoveItem(item.K)
	}
	entry := &entry{
		prev: s.head.prev,
		next: s.head,
		item: &item,
	}
	s.head.prev.next = entry
	s.head.prev = entry

	s.indexed[item.K] = entry
}

func (s *memoryStorage) RemoveItem(key string) error {
	entry, ok := s.indexed[key]

	if !ok {
		return errors.New(fmt.Sprintf("key `%s' not found", key))
	}

	delete(s.indexed, key)
	entry.prev.next = entry.next
	entry.next.prev = entry.prev

	// cleanup references from removed node
	entry.prev = nil
	entry.next = nil
	entry.item = nil

	return nil
}

func (s *memoryStorage) GetItem(key string) (*Item, error) {
	entry, ok := s.indexed[key]

	if !ok {
		return nil, errors.New(fmt.Sprintf("key `%s' not found", key))
	}

	// return item copy
	item := *entry.item

	return &item, nil
}

func (s *memoryStorage) GetAllItems() []Item {
	result := make([]Item, 0, len(s.indexed))
	for e := s.head.next; e != s.head; e = e.next {
		result = append(result, *e.item)
	}
	return result
}

func (s *memoryStorage) Iterate(accept func(Item)) {
	for e := s.head.next; e != s.head; e = e.next {
		accept(*e.item)
	}
}
