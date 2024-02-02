# Memory Cache in Go

This is a simple in-memory cache implementation in Go. It provides basic functionalities to set, get and invalidate items in the cache.
The cache is thread-safe and uses a `sync.RWMutex` to store the items.
Automatically removes items from the cache after a specified duration.

## Usage

### Creating a new cache

```go
cache := NewCache[int](time.Second * 60)
```

### Setting a value in the cache

```go
cache.Set("key", 123)
```

You can also set a value with a specific life duration:

```go
cache.SetWithLifeDuration("key", 123, time.Second*10)
```

### Getting a value from the cache

```go
value, exists := cache.Get("key")
if exists {
    fmt.Println(value)
}
```

### Invalidating the cache

```go
cache.Invalidate()
```

## Limitations

This cache implementation has a few limitations:

1. **Reference types are not deep-copied**: When you get a value from the cache, if the value is a reference type (like a slice or a map), and you modify the returned value, the value in the cache will also be modified. This is because the cache returns a reference to the value, not a copy of the value. To avoid this, you need to manually make a copy of the returned value before modifying it.

Here's an example:

```go
cache := NewCache[[]byte](time.Second * 3)

original := []byte{1, 2, 3}
cache.Set("key", original)

// Get the value from the cache
value, _ := cache.Get("key")

// Modify the returned value
value[0] = 9

// The value in the cache is also modified
fmt.Println(original)         // Prints: [9 2 3]
fmt.Println(cache.Get("key")) // Prints: [9 2 3]
```

To avoid this, make a copy of the returned value before modifying it:

```go
value, _ := cache.Get("key")
copyOfValue := make([]byte, len(value))
copy(copyOfValue, value)

// Now you can modify copyOfValue without affecting the value in the cache
```

## Conclusion

This is a basic in-memory cache implementation in Go. It's suitable for simple use cases where you need a temporary, in-memory storage of items for fast access. For more complex use cases, consider using a more feature-rich cache library or a dedicated cache server like Redis or Memcached.











