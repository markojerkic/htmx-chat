package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	guuid "github.com/google/uuid"
)

type Item interface {
	GetID() string
	SetID(ID string) Item
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
		fmt.Println("Item already has an ID: '", item.GetID(), "'")
		id = item.GetID()
	} else {
		fmt.Println("Item does not have an ID, creating new")
		id = guuid.New().String()
		item = item.SetID(id).(T)
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

func (c *Collection[T]) GetByPredicate(predicate func(T) bool) (T, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, item := range c.values {
		if predicate(item) {
			return item, nil
		}
	}

	return *new(T), fmt.Errorf("Item not found")
}

func (c *Collection[T]) GetAllByPredicate(predicate func(T) bool) []T {
	c.lock.RLock()
	defer c.lock.RUnlock()

	items := make([]T, 0, len(c.values))

	for _, item := range c.values {
		if predicate(item) {
			items = append(items, item)
		}
	}

	return items
}
