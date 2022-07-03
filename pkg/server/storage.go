package server

import (
	"errors"
	"fmt"
	"sync"
)

// Storage defines interface for data storage
type Storage interface {
	// AddItem adds item to storage
	AddItem(item string)
	// RemoveItem removes item from storage
	RemoveItem(itemID int) error
	// GetItem returns item with given id from storage
	GetItem(itemID int) (string, error)
	// GetAllItems returns items in storage, new slice created every time
	GetAllItems() []string
	// Iterate allows to iterato over data in storage, can be used for processing which does not involve blocking IO
	Iterate(accept func(int, string))
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

func (s *rwLockedStorage) AddItem(item string) {
	s.rwLock.Lock()
	s.storage.AddItem(item)
	s.rwLock.Unlock()
}

func (s *rwLockedStorage) RemoveItem(itemID int) error {
	s.rwLock.Lock()
	err := s.storage.RemoveItem(itemID)
	s.rwLock.Unlock()
	return err
}

func (s *rwLockedStorage) GetItem(itemID int) (string, error) {
	s.rwLock.RLock()
	item, err := s.storage.GetItem(itemID)
	s.rwLock.RUnlock()
	return item, err
}

func (s *rwLockedStorage) GetAllItems() []string {
	s.rwLock.RLock()
	items := s.storage.GetAllItems()
	s.rwLock.RUnlock()
	return items
}

func (s *rwLockedStorage) Iterate(accept func(int, string)) {
	s.rwLock.RLock()
	s.storage.Iterate(accept)
	s.rwLock.RUnlock()
}

type arrayStorage struct {
	data []string
}

// NewSliceStorage returns storage backed by slice, item id is an index in slice
func NewSliceStorage(cap uint) *arrayStorage {
	return &arrayStorage{
		data: make([]string, 0, cap),
	}
}

func (s *arrayStorage) AddItem(item string) {
	s.data = append(s.data, item)
}

func (s *arrayStorage) RemoveItem(itemID int) error {
	if itemID < 0 || itemID >= len(s.data) {
		return errors.New(fmt.Sprintf("index out of bounds: %d", itemID))
	}

	s.data = append(s.data[:itemID], s.data[itemID+1:]...)

	return nil
}

func (s *arrayStorage) GetItem(itemID int) (string, error) {
	if itemID < 0 || itemID >= len(s.data) {
		return "", errors.New(fmt.Sprintf("index out of bounds: %d", itemID))
	}

	return s.data[itemID], nil
}

func (s *arrayStorage) GetAllItems() []string {
	result := make([]string, len(s.data))
	copy(result, s.data)
	return result
}

func (s *arrayStorage) Iterate(accept func(int, string)) {
	for i, v := range s.data {
		accept(i, v)
	}
}
