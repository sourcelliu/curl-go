package tool

// StringList is a translation of the C struct `slist_wc` from
// curl-src/src/slist_wc.h.
//
// The C implementation uses a singly-linked list with a cached tail pointer
// for efficient appends. The idiomatic and more performant Go equivalent is
// a slice of strings, which provides amortized O(1) appends.
type StringList struct {
	data []string
}

// NewStringList creates and returns a new, empty StringList.
func NewStringList() *StringList {
	return &StringList{
		data: []string{},
	}
}

// Append adds a string to the list. This is a translation of the C function
// `slist_wc_append` from curl-src/src/slist_wc.c.
func (l *StringList) Append(s string) {
	l.data = append(l.data, s)
}

// Strings returns the underlying slice of strings.
func (l *StringList) Strings() []string {
	return l.data
}

// The C file also contains `slist_wc_free_all` to deallocate the linked list.
// This is not needed in Go due to automatic garbage collection. The StringList
// and its underlying data will be automatically garbage-collected when they
// are no longer referenced.