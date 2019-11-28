package general

import (
	"container/list"
)

func ListToArray(l *list.List) []interface{} {
	var size int = l.Len()

	var items = make([]interface{}, size)
	var i int
	for e := l.Front(); e != nil; e = e.Next() {
		items[i] = l
		i++
	}
	return items
}
