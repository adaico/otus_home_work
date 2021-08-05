package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.Lock()
	defer lc.Unlock()

	item, ok := lc.items[key]

	if ok {
		item.Value = cacheItem{key, value}
		lc.queue.MoveToFront(item)

		return true
	}

	listItem := lc.queue.PushFront(cacheItem{key, value})
	lc.items[key] = listItem

	if lc.queue.Len() > lc.capacity {
		backItem := lc.queue.Back()

		delete(lc.items, backItem.Value.(cacheItem).key)
		lc.queue.Remove(backItem)
	}

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.Lock()
	defer lc.Unlock()
	item, ok := lc.items[key]

	if ok {
		lc.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.Lock()
	defer lc.Unlock()

	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
