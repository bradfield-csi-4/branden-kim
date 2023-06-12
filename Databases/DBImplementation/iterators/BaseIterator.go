package iterators

type BaseIterator interface {
	Walk()
}

type BaseIteratorFields struct {
	children []*BaseIterator
}
