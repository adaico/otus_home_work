package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{v, nil, nil}

	if l.front == nil && l.back == nil {
		l.back = &item
	} else {
		bindListItems(&item, l.front)
	}

	l.front = &item
	l.len++

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{v, nil, nil}

	if l.front == nil && l.back == nil {
		l.front = &item
	} else {
		bindListItems(l.back, &item)
	}

	l.back = &item
	l.len++

	return &item
}

func (l *list) Remove(i *ListItem) {
	bindListItems(i.Prev, i.Next)

	if l.front == i {
		l.front = i.Next
	}

	if l.back == i {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}

func bindListItems(i1 *ListItem, i2 *ListItem) {
	if i1 != nil {
		i1.Next = i2
	}

	if i2 != nil {
		i2.Prev = i1
	}
}
