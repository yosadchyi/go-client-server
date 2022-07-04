package server_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yosadchyi/go-client-server/pkg/server"
)

func TestMemoryStorage(t *testing.T) {
	cases := map[string]struct {
		scenario      func(storage server.Storage) (*server.Item, error)
		expectedError error
		expectedItem  *server.Item
		expectedItems []server.Item
	}{
		"Adding item to empty storage": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "K", V: "V"})
				return nil, nil
			},
			expectedItems: []server.Item{{K: "K", V: "V"}},
		},
		"Adding more than one item to empty storage": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				storage.AddItem(server.Item{K: "3", V: "C"})
				return nil, nil
			},
			expectedItems: []server.Item{{K: "1", V: "A"}, {K: "2", V: "B"}, {K: "3", V: "C"}},
		},
		"Replacing item with the same key": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "1", V: "B"})
				return nil, nil
			},
			expectedItems: []server.Item{{K: "1", V: "B"}},
		},
		"Removing single item": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				return nil, storage.RemoveItem("1")
			},
			expectedItems: []server.Item{},
		},
		"Removing last item": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				return nil, storage.RemoveItem("2")
			},
			expectedItems: []server.Item{{K: "1", V: "A"}},
		},
		"Removing first item": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				return nil, storage.RemoveItem("1")
			},
			expectedItems: []server.Item{{K: "2", V: "B"}},
		},
		"Removing non existing item": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				return nil, storage.RemoveItem("3")
			},
			expectedError: errors.New("key `3' not found"),
			expectedItems: []server.Item{{K: "1", V: "A"}, {K: "2", V: "B"}},
		},
		"Getting non existing item, removed before": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				storage.AddItem(server.Item{K: "3", V: "C"})
				err := storage.RemoveItem("3")
				if err != nil {
					return nil, err
				}
				return storage.GetItem("3")
			},
			expectedError: errors.New("key `3' not found"),
			expectedItems: []server.Item{{K: "1", V: "A"}, {K: "2", V: "B"}},
		},
		"Getting existing item": {
			scenario: func(storage server.Storage) (*server.Item, error) {
				storage.AddItem(server.Item{K: "1", V: "A"})
				storage.AddItem(server.Item{K: "2", V: "B"})
				storage.AddItem(server.Item{K: "3", V: "C"})
				return storage.GetItem("3")
			},
			expectedItem:  &server.Item{K: "3", V: "C"},
			expectedItems: []server.Item{{K: "1", V: "A"}, {K: "2", V: "B"}, {K: "3", V: "C"}},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			storage := server.NewMemoryStorage()
			item, err := tc.scenario(storage)
			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err, tc.expectedError) {
					assert.Equal(t, tc.expectedError, err)
				}
			}
			assert.Equal(t, tc.expectedItem, item)
			assert.Equal(t, tc.expectedItems, storage.GetAllItems())
		})
	}
}
