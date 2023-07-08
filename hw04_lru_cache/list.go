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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{Value: v}
	if l.front == nil {
		l.front = li
		l.back = li
	} else {
		l.front.Prev = li
		li.Next = l.front
		l.front = li
	}
	l.len++

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{Value: v}
	if l.front == nil {
		l.front = li
		l.back = li
	} else {
		li.Prev = l.back
		l.back.Next = li
		l.back = li
	}
	l.len++

	return li
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.len == 0 || i == nil:
		return
	case l.len == 1:
		// Если список из одного элемента и это наш элемент - обнуляем список
		if i == l.front {
			l.front = nil
			l.back = nil
		} else {
			return
		}
	case i == l.front:
		// Если указатель равен первому элементу
		l.front = i.Next
		l.front.Prev = nil
	case i == l.back:
		// Если указатель равен последнему элементу
		l.back = i.Prev
		l.back.Next = nil
	default:
		// Изменяем связи соседних элементов
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i == nil:
		return
	case l.len == 1 || i == l.front:
		// Если список из одного элемента или это и так первый элемент - ничего не делаем
		return
	case i == l.back:
		// Если элемент - последний в списке
		i.Prev.Next = nil
		l.back = i.Prev
	default:
		// Изменяем связи соседних элементов
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	// Перемещаем текущий элемент в начало списка
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
