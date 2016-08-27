package data

// Collection represents a collection of items
type collection struct {
	// items is the result
	items interface{}
	// colkeys is the collection's key for retrieval
	colkeys map[string]interface{}
	// off is the offset
	off int
	// lim is the limit
	lim int
	// colname is the collection name
	colname string
}

// Items gets the items value from the collection
func (c *collection) Items() interface{} {
	return c.items
}

// offset gets the collection's offset
func (c *collection) offset() int {
	return c.off
}

// limit gets the collection's limit
func (c *collection) limit() int {
	return c.lim
}

// cname gets the collection's name
func (c *collection) cname() string {
	return c.colname
}

// keys gets the key set of this collection
func (c *collection) keys() map[string]interface{} {
	return c.colkeys
}
