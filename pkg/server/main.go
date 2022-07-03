package server

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yosadchyi/go-client-server/pkg/message"
)

// NewProcessFn creates new processing function for messages
func NewProcessFn(id int, storage Storage, logFile *os.File) func(*message.Any) {
	name := fmt.Sprintf("processor-%d", id)

	return func(m *message.Any) {
		writeLog(logFile, string(m.Operation))
		switch {
		case m.Add != nil:
			storage.AddItem(Item{
				K: m.Add.Key,
				V: m.Add.Data,
			})
			log.Printf("%s: adding item %s with key %s", name, m.Add.Data, m.Add.Key)
			writeLog(logFile, fmt.Sprintf("%s:%s", m.Add.Key, m.Add.Data))
		case m.Remove != nil:
			key := m.Remove.Key
			err := storage.RemoveItem(key)
			if err != nil {
				log.Printf("%s: can't remove item with key %s", name, key)
			} else {
				log.Printf("%s: removing item with key %s", name, key)
				writeLog(logFile, m.Remove.Key)
			}
		case m.GetItem != nil:
			key := m.GetItem.Key
			item, err := storage.GetItem(key)
			if err != nil {
				log.Printf("%s: can't get item with key %s", name, key)
			}
			log.Printf("%s: Get(%s): %s", name, key, item.V)
			writeLog(logFile, fmt.Sprintf("%s:%s", item.K, item.V))
		case m.GetAllItems != nil:
			items := storage.GetAllItems()
			log.Printf("%s: listing all items:", name)
			for _, item := range items {
				log.Printf("%s: (%s, %s)", name, item.K, item.V)
				writeLog(logFile, fmt.Sprintf("%s:%s", item.K, item.V))
			}
		}
	}
}

func writeLog(logFile *os.File, logEntry string) {
	_, err := io.WriteString(logFile, fmt.Sprintf("%s\n", logEntry))
	if err != nil {
		log.Printf("error writing log %s", err.Error())
	}

	if err := logFile.Sync(); err != nil {
		log.Printf("error writing log %s", err.Error())
	}
}
