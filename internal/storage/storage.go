// Storage package is module for manage key - value storage
// Outsiders can use following; Get, Put, Delete, Clear, which are self explanatory
package storage

import (
	"sync"
	"time"
)

// Storage structure is struct for holding data structure and misc.
// Cache that I want to build should use hash table, and doubly linked list.
// - Hash Table: since I'm going to implement key-value store, hash table should be good choice since it has O(logN) to insert and search.
// - Double Linked List: lenght of data would be limited, and eviction will be happen in LRU manner(Least Recently Used). To implement this, I will use double linked list here.
type Storage struct {
	table  map[string]*Node
	head   *Node
	tail   *Node
	size   int64
	mutex  *sync.Mutex
	config StorageConfig
}

// StorageConfig structure should be provided when outside code calls New() function. It will set properties of storage such as ttl or capacity
type StorageConfig struct {
	Ttl      time.Duration
	Capacity int64
}

func New(config StorageConfig) *Storage {
	return &Storage{
		table:  make(map[string]*Node),
		head:   nil,
		tail:   nil,
		size:   0,
		mutex:  &sync.Mutex{},
		config: config,
	}
}

type Node struct {
	key  string
	data []byte
	ttl  time.Time
	prev *Node
	next *Node
}

func (s *Storage) Get(key string) ([]byte, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	node, ok := s.table[key]
	if !ok {
		return nil, false
	}

	if node.ttl.Before(time.Now()) {
		s.evict(node)
		s.size--
		return nil, false
	}

	return node.data, true
}

func (s *Storage) Put(key string, data []byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	node, ok := s.table[key]

	if ok {
		s.table[key].data = data
		s.setHead(node)
		return true
	}

	for s.size >= s.config.Capacity {
		s.evict(s.tail)
		s.size--
	}

	ttl := time.Now().Add(s.config.Ttl)
	newNode := &Node{
		key:  key,
		data: data,
		ttl:  ttl,
	}
	s.table[key] = newNode
	s.setHead(newNode)
	s.size++

	return false
}

func (s *Storage) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	node, ok := s.table[key]
	if !ok {
		return false
	}

	s.evict(node)
	s.size--

	return true
}

func (s *Storage) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.head != nil {
		s.evict(s.tail)
	}
	s.size = 0
}

func (s *Storage) Size() int64 {
	return s.size
}

func (s *Storage) Cap() int64 {
	return s.config.Capacity
}

func (s *Storage) RemoveExpired() int64 {
	var count int64 = 0
	for key, value := range s.table {
		if value.ttl.Before(time.Now()) {
			s.Delete(key)
			count++
		}
	}
	return count
}

// setHead function is move node to head of linked list
// if
func (s *Storage) setHead(node *Node) {
	if node == nil {
		return
	}

	if s.head == nil {
		s.head = node
		s.tail = node
		return
	}

	if s.head == node {
		return
	}

	if s.tail == node {
		s.tail = s.tail.prev
		s.tail.next = nil
		node.prev = nil
		s.head.prev = node
		node.next = s.head
		s.head = node
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	s.head.prev = node
	node.next = s.head
	s.head = node
}

// evict is
func (s *Storage) evict(node *Node) {
	if s.head == s.tail && s.head == node {
		s.head = nil
		s.tail = nil
		node.prev = nil
		node.next = nil
		delete(s.table, node.key)
		return
	}

	if s.head == node {
		s.head = s.head.next
		s.head.prev = nil
		node.prev = nil
		node.next = nil
		delete(s.table, node.key)
		return
	}

	if s.tail == node {
		s.tail = s.tail.prev
		s.tail.next = nil
		node.prev = nil
		node.next = nil
		delete(s.table, node.key)
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = nil
	delete(s.table, node.key)
}
