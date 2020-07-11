/*
In Bithose, the connectionstore holds
connectionstore holds all the objects that make up the connectionstore.
*/

package connectionstore

import (
	"encoding/json"
	"errors"
	UuidLib "github.com/google/uuid"
	"sync"
	"time"
)

var (
	ConnectionAlreadyExistsErr = errors.New("connection already exists")

	mapStoreInstance *MapStore
)

type MapStore struct {
	mtx         *sync.RWMutex
	connections map[string]*Connection
	stats       *Statistics
}

func GetMapStore() ConnectionStore {
	if mapStoreInstance == nil {
		mapStoreInstance = &MapStore{
			mtx:         &sync.RWMutex{},
			connections: map[string]*Connection{},
			stats:       NewStatistics(),
		}
	}
	return mapStoreInstance
}

func (m *MapStore) AddConnection(connection *Connection) (string, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// if no uuid is provided, generate one
	u := UuidLib.New()
	uuid := u.String()

	m.connections[uuid] = connection
	m.stats.IncrementConnection()
	return uuid, nil
}

func (m *MapStore) RemoveConnection(uuid string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exists := m.GetConnection(uuid); !exists {
		return
	}

	delete(m.connections, uuid)
	m.stats.DecrementConnection()
}

func (m *MapStore) SendMessage(message Message) (numOfSent int, numOfTimeouts int, err error) {
	// convert message to json
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return 0, numOfTimeouts, err
	}

	// since we are iterating the map, we should lock it from writes
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for uuid, connection := range m.connections {
		accepts, err := connection.AcceptsLabels(message.LabelPairs)
		if err != nil {
			return 0, numOfTimeouts, err
		}
		if accepts {
			select {
			case connection.Ch <- jsonMsg:
				numOfSent++
				m.stats.IncrementMessageSent()
			// if we have to wait, we will assume the connection is stale and delete it
			// NOTE: the wait time has to be as small as possible as this may cause delays
			case <-time.After(time.Millisecond * 100):
				numOfTimeouts++
				m.stats.IncrementMessageTimeout()
				close(connection.Ch)
				m.RemoveConnection(uuid)
			}
		}
	}

	return numOfSent, numOfTimeouts, nil
}

func (m *MapStore) GetConnection(uuid string) (connection *Connection, exists bool) {
	if connection, ok := m.connections[uuid]; ok {
		return connection, ok
	} else {
		return nil, ok
	}
}

func (m *MapStore) Stats() *Statistics {
	return m.stats
}
