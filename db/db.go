package db

import (
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"os"
	"sync"
)

type Item interface {
	GetID() string
	SetID(ID string)
}
type Collection[T Item] struct {
	lock   sync.RWMutex
	values map[string]T
	Name   string
}

func NewCollection[T Item](name string) *Collection[T] {
	c := &Collection[T]{
		values: make(map[string]T),
		Name:   name,
	}
	c.values = c.readFromTmp()
	return c
}

func (c *Collection[T]) syncToFile() error {
	values, err := json.Marshal(c.values)
	if err != nil {
		fmt.Printf("Error marshalling values for collection: %s", c.Name)
		return err
	}

	err = os.WriteFile(fmt.Sprintf("/tmp/chat-%s", c.Name), []byte(values), 0777)

	if err != nil {
		fmt.Printf("Error writing to file for collection: %s", c.Name)
		return err
	}
	return nil
}

func (c *Collection[T]) readFromTmp() map[string]T {
	values, err := os.ReadFile(fmt.Sprintf("/tmp/chat-%s", c.Name))
	var valuesMap map[string]T
	if err != nil {
		fmt.Println("No values file found, creating new")
		valuesMap = make(map[string]T)
	} else {
		err = json.Unmarshal(values, &valuesMap)
		fmt.Println("Read values from file", valuesMap)
		if err != nil {
			fmt.Println("Error unmarshalling values")
			panic(err)
		}
	}

	return valuesMap
}

func (c *Collection[T]) Save(item T) (T, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var id string

	if item.GetID() != "" {
		id = item.GetID()
	} else {
		id = guuid.New().String()
		item.SetID(id)
	}

	c.values[id] = item

	err := c.syncToFile()

	return item, err
}

func (c *Collection[T]) Get(id string) (T, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.values[id]

	if !ok {
		return item, fmt.Errorf("Item not found")
	}

	return item, nil
}

func (c *Collection[T]) GetAll() (map[string]T, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.values, nil
}
