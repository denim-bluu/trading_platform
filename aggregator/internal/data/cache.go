package data

import (
	"sync"
	"time"
	pb "trading_platform/aggregator/proto"
)

type cacheItem struct {
	data      []*pb.DataPoint
	startDate time.Time
	endDate   time.Time
	timestamp time.Time
}

type Cache struct {
	data map[string]cacheItem
	ttl  time.Duration
	lock sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
		ttl:  ttl,
	}
}

func (c *Cache) Set(symbol string, startDate, endDate time.Time, data []*pb.DataPoint) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[symbol] = cacheItem{
		data:      data,
		startDate: startDate,
		endDate:   endDate,
		timestamp: time.Now(),
	}
}

func (c *Cache) Get(symbol string, startDate, endDate time.Time) ([]*pb.DataPoint, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, exists := c.data[symbol]
	if !exists || time.Since(item.timestamp) > c.ttl {
		return nil, false
	}

	// Check if the cached data covers the requested date range
	if startDate.Before(item.startDate) || endDate.After(item.endDate) {
		return nil, false
	}

	// Filter the cached data to the requested date range
	var filteredData []*pb.DataPoint
	for _, dp := range item.data {
		dpTime, _ := time.Parse("2006-01-02", dp.Timestamp)
		if (dpTime.Equal(startDate) || dpTime.After(startDate)) && (dpTime.Equal(endDate) || dpTime.Before(endDate)) {
			filteredData = append(filteredData, dp)
		}
	}

	return filteredData, true
}

func (c *Cache) Invalidate(symbol string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, symbol)
}
