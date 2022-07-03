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
			storage.AddItem(m.Add.Data)
			log.Printf("%s: adding item %s", name, m.Add.Data)
		case m.Remove != nil:
			id := m.Remove.ItemID
			err := storage.RemoveItem(id)
			if err != nil {
				log.Printf("%s: can't remove item %d", name, id)
			} else {
				log.Printf("%s: removing item %d", name, id)
			}
		case m.GetItem != nil:
			id := m.GetItem.ItemID
			item, err := storage.GetItem(id)
			if err != nil {
				log.Printf("%s: can't get item %d", name, id)
			}
			log.Printf("%s: Get(%d): %s", name, id, item)
		case m.GetAllItems != nil:
			items := storage.GetAllItems()
			log.Printf("%s: listing all items:", name)
			for i, item := range items {
				log.Printf("%s: %d: %s", name, i, item)
			}
		}
	}
}
