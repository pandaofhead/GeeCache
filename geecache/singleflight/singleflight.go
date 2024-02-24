package singleflight

import "sync"

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup //A WaitGroup waits for a collection of goroutines to finish
	val interface{}
	err error
}

// Group handles duplicate requests for the same key
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// if the key is already in-flight, wait for it
	if c, ok := g.m[key]; ok {
		g.mu.Unlock() // unlock before wg.Wait()
		c.wg.Wait()   // wait for the call to complete
		return c.val, c.err
	}

	// the key is not in-flight; make the fn call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c

	// unlock and get key from remote
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()

	// remove the key from the map
	// prevents memory leak and updates the map
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
