package server

import (
	"fmt"
	"log"

	"github.com/yosadchyi/go-client-server/pkg/message"
)

// NewProcessFn creates new processing function for messages
func NewProcessFn(id int, storage Storage) func(*message.Any) {
	name := fmt.Sprintf("processor-%d", id)

	return func(m *message.Any) {
		switch {
		case m.Add != nil:
			storage.AddItem(Item{
				K: m.Add.Key,
				V: m.Add.Data,
			})
			log.Printf("%s: adding item %s with key %s", name, m.Add.Data, m.Add.Key)
		case m.Remove != nil:
			key := m.Remove.Key
			err := storage.RemoveItem(key)
			if err != nil {
				log.Printf("%s: can't remove item with key %s", name, key)
			} else {
				log.Printf("%s: removing item with key %s", name, key)
			}
		case m.GetItem != nil:
			key := m.GetItem.Key
			item, err := storage.GetItem(key)
			if err != nil {
				log.Printf("%s: can't get item with key %s", name, key)
			}
			log.Printf("%s: Get(%s): %s", name, key, item.V)
		case m.GetAllItems != nil:
			items := storage.GetAllItems()
			log.Printf("%s: listing all items:", name)
			for _, item := range items {
				log.Printf("%s: (%s, %s)", name, item.K, item.V)
			}
		}
	}
}
