# cachedmap

A Go string keyed object cache with entry timeout and periodic flushing.

Cachedmap provides the equivalent of a cached `map[string]interface{}` store with per record timeouts.

Overall cache growth is managed by a periodic global flush of entire cache. For this reason it may be more suitable for applications where you expect a high hit rate but when a miss is not really so expensive. Or if you know that growth is not going to be a problem. In that case you could set the global flush period to some high value.

# Usage example



```
cm := cachedmap.NewCachedMap(5*time.Second, 5*time.Minute, nil)

cm.Set("cake", &monkey)

...

m, ok := cm.Get("cake")
if ok {
    freshMonkey = (m).(*monkeyType)
}
```

See the examples directory for more.


# Alternatives

- https://github.com/muesli/cache2go is similar and provides more functionality
