package cachedmap

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type CachedMap struct {
	Hits       int64
	Misses     int64
	Writes     int64
	Flushes    int64
	MaxLength  int
	entries    map[string]CacheEntry
	lock       sync.RWMutex
	keyTimeout time.Duration
	flushCycle time.Duration

	Log *logrus.Entry
}

type CacheEntry struct {
	Data       interface{}
	RemoveTime time.Time
}

type Stats struct {
	Hits      int64 `json:"hits"`
	Misses    int64 `json:"misses"`
	Writes    int64 `json:"writes"`
	Flushes   int64 `json:"flushes"`
	MaxLength int   `json:"max_length"`
	Length    int   `json:"length"`
}

func NewCachedMap(keyTimeout, flushCycle time.Duration, log *logrus.Entry) *CachedMap {
	cm := &CachedMap{
		keyTimeout: keyTimeout,
		flushCycle: flushCycle,
		entries:    make(map[string]CacheEntry, 100),
		Log:        log,
	}
	cm.flusher()
	return cm
}

// Tidy starts a go routine which will periodically drop the entire cache.
func (c *CachedMap) flusher() {
	go func() {
		for {
			time.Sleep(c.flushCycle)
			c.lock.Lock()
			oldEntries := c.entries
			c.entries = make(map[string]CacheEntry, 10+len(c.entries))
			c.lock.Unlock()
			c.Flushes++
			if len(oldEntries) > c.MaxLength {
				c.MaxLength = len(oldEntries)
			}
			if c.Log != nil {
				c.Log.Info(c.String())
			}
		}
	}()
}

// Set adds an item to the map with a computed remove time.
// It returns the remove time in case you want to use it in a
// nearby SetUntil() call.
func (c *CachedMap) Set(key string, data interface{}) time.Time {
	removeTime := time.Now().Add(c.keyTimeout)
	c.SetUntil(key, data, removeTime)
	return removeTime
}

// SetUntil adds an item to the map with the given remove time.
// This allows you to override the default key timeout.
func (c *CachedMap) SetUntil(key string, data interface{}, removeTime time.Time) {
	e := CacheEntry{
		Data:       data,
		RemoveTime: removeTime,
	}
	c.lock.Lock()
	c.entries[key] = e
	c.Writes++
	c.lock.Unlock()
}

// Len returns the number of items in the map.
func (c *CachedMap) Len() int {
	return len(c.entries)
}

// Get returns a non-expired entry and true.
// If no valid entry is found, it returns (nil, false).
func (c *CachedMap) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	e, ok := c.entries[key]
	c.lock.RUnlock()
	if !ok || time.Now().After(e.RemoveTime) {
		atomic.AddInt64(&c.Misses, 1)
		return nil, false
	}
	atomic.AddInt64(&c.Hits, 1)
	return e.Data, true
}

func (c *CachedMap) String() string {
	return fmt.Sprintf("<len=%d maxlen=%d hits=%d misses=%d writes=%d flushes=%d ttl=%s fc=%s>",
		c.Len(), c.MaxLength,
		c.Hits, c.Misses,
		c.Writes,
		c.Flushes,
		c.keyTimeout, c.flushCycle,
	)
}

func (c *CachedMap) GetStats() Stats {
	// MaxLength is only set at flush time so we might need to update it here.
	l := c.Len()
	m := c.MaxLength
	if l > m {
		m = l
	}
	return Stats{
		Hits:      c.Hits,
		Misses:    c.Misses,
		Writes:    c.Writes,
		Flushes:   c.Flushes,
		MaxLength: m,
		Length:    l,
	}
}
