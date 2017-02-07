// Package prefixsearch implements simple tree-based prefix search that
// i'm using for different web autocomplete services
package prefixsearch

import (
	"strings"
	"unicode"
)

// SearchTree is struct to handle search tree
type SearchTree struct {
	root *node
}

type node struct {
	values    []interface{}
	childnum uint
	childs   map[rune]*node
}

// New creates new search tree
func New() *SearchTree {
	return &SearchTree{
		root: &node{childs: map[rune]*node{}},
	}
}

// Add one leaf to tree
func (tree *SearchTree) Add(key string, value interface{}) {
	current := tree.root

	needUpdate := (nil == tree.Search(key))

	for _, sym := range strings.ToLower(key) {
		if needUpdate {
			current.childnum++
		}
		next, ok := current.childs[sym]
		if !ok {
			newone := &node{childs: map[rune]*node{}}
			current.childs[sym] = newone
			next = newone
		}
		current = next
	}

	if needUpdate {
		current.childnum++
	}
	current.values = append(current.values, value)
}

// AutoComplete returns autocomplete suggestions for given prefix
func (tree *SearchTree) AutoComplete(prefix string) []interface{} {
	// walk thru prefix symbols
	current := tree.root
	for _, sym := range prefix {
		var ok bool
		current, ok = current.childs[unicode.ToLower(sym)]
		if !ok {
			return []interface{}{}
		}
	}

	// we have found, now very stupid tree walk :)
	result := make([]interface{}, 0, current.childnum)
	current.recurse(func(v []interface{}) {
		if nil != v {
			for _, val := range v {
				result = append(result, val)
			}
		}
	})
	return result
}

// Search searches for value of key
func (tree *SearchTree) Search(key string) []interface{} {
	current := tree.root
	for _, sym := range key {
		var ok bool
		current, ok = current.childs[unicode.ToLower(sym)]
		if !ok {
			return nil
		}
	}
	return current.values
}

func (n *node) recurse(callback func([]interface{})) {
	callback(n.values)
	for _, v := range n.childs {
		v.recurse(callback)
	}
}
