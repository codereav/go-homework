package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type ItemValue struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	mutex    sync.Mutex
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock() // Блокируем c.items[key]
	i, wasInCache := c.items[key]
	if wasInCache { // Если запись была добавлена ранее - обновляем значение, двигаем её в начало списка
		i.Value.(*ItemValue).Value = value
		c.queue.MoveToFront(i)
	} else {
		if c.queue.Len() == c.capacity { // Если достигли максимальной емкости - удаляем последний элемент
			least := c.queue.Back()
			c.queue.Remove(least)
			delete(c.items, least.Value.(*ItemValue).Key)
		}

		li := c.queue.PushFront(&ItemValue{Key: key, Value: value}) // Добавляем значение в список
		c.items[key] = li                                           // Добавляем/обновляем в map

	}
	c.mutex.Unlock() // Разблокируем c.items[key]
	return wasInCache
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock() // Блокируем c.items[key]
	v, ok := c.items[key]
	c.mutex.Unlock() // Разблокируем c.items[key]
	if ok {
		c.queue.MoveToFront(v)
		return v.Value.(*ItemValue).Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mutex.Unlock()
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
