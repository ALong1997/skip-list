# Skip List

This is a simple and thread-safe skip list implemented by [Go](https://go.dev/) with generics.

## Convenience functions

Other than the classic `Get`, `Put` and `Delete`, some more convenience functions are implemented that makes this skiplist implementation very easy and straight forward to use
in real applications. All complexity values are approximates, as skiplist can only approximate runtime complexity.

| Function | Complexity | Description                                                        |
|----------|:----------:|:-------------------------------------------------------------------|
| Level    |    O(1)    | returns the level of the skiplist                                  |
| Cap      |    O(1)    | returns the number of invalid nodes                                |
| Get      | O(log(n))  | returns the value of a given key and whether it is valid           |
| Put      | O(log(n))  | inserts or updates the value of a given key                        |
| Delete   | O(log(n))  | deletes a node for a given key                                     |
| Range    | O(log(n))  | returns kv-pairs of a given key range                              |
| Ceil     | O(log(n))  | returns kv-pairs of the least key greater than or equal to target  |
| Floor    | O(log(n))  | returns kv-pairs of the greatest key less than or equal to target. |


## Getting started

### Prerequisites
- **[Go](https://go.dev/) version 1.18+**

### Getting
With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/ALong1997/skiplist"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `skiplist` package:

```sh
$ go get -u https://github.com/ALong1997/skiplist
```

### Quick Start

```go
package main

import (
    "fmt"

	"github.com/ALong1997/skiplist"
)

func main() {
    sl := skiplist.NewSkipList[int, int](8, true)
    
    for i := 1; i <= 100; i++ {
        sl.Put(i, i)
    }
    
    // randomly generate level, but less than maxLevel
    fmt.Println(sl.Level())
    
    // 100
    fmt.Println(sl.Cap())
    
    // 10
    if v, ok := sl.Get(10); ok {
        fmt.Println(v)
    }
    
    // 10 10
    // 11 11
    for _, v := range sl.Range(10, 11) {
        fmt.Println(v.Key(), v.Val())
    }
    
    // 1
    fmt.Println(sl.Ceil(0))
    
    // 100
    fmt.Println(sl.Floor(1000))
    
    for i := 1; i <= 100; i++ {
        sl.Delete(i)
    }
}

```
