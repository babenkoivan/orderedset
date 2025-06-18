# orderedset

Package orderedset provides a generic, goroutine-safe ordered set implementation in Go.

It maintains insertion order of unique elements and supports common set operations efficiently.

## Installation

```bash
go get github.com/babenkoivan/orderedset
```

## Overview

The `OrderedSet[T comparable]` type supports:

* Adding and removing unique elements
* Checking element presence
* Accessing elements by index
* Finding index of an element
* Removing elements by value or index
* Cloning and slicing subsets
* Sorting elements by custom comparator
* Set operations: Union, Intersect, Difference
* JSON marshalling/unmarshalling
* Thread-safe operations for concurrent use

## Usage

```go
package main

import (
    "fmt"
    "github.com/babenkoivan/orderedset"
)

func main() {
    s := orderedset.New[int]()
    s.Add(1)
    s.Add(2)
    s.Add(2) // duplicate ignored

    fmt.Println("Values:", s.Values()) // Output: Values: [1 2]

    s.Remove(1)
    fmt.Println("Has 1?", s.Has(1))   // Output: Has 1? false
    fmt.Println("Length:", s.Len())   // Output: Length: 1

    val, ok := s.At(0)
    if ok {
        fmt.Println("Value at index 0:", val) // Output: Value at index 0: 2
    }

    index := s.IndexOf(2)
    fmt.Println("Index of 2:", index) // Output: Index of 2: 0

    removedVal, removed := s.RemoveAt(0)
    if removed {
        fmt.Println("Removed value:", removedVal) // Output: Removed value: 2
    }

    s.Add(3)
    s.Add(4)

    s2 := orderedset.New[int]()
    s2.Add(4)
    s2.Add(5)

    union := s.Union(s2)
    fmt.Println("Union:", union.Values()) // Output: Union: [3 4 5]

    intersect := s.Intersect(s2)
    fmt.Println("Intersect:", intersect.Values()) // Output: Intersect: [4]

    diff := s.Difference(s2)
    fmt.Println("Difference:", diff.Values()) // Output: Difference: [3]

    slice, err := union.Slice(1, 3)
    if err == nil {
        fmt.Println("Slice:", slice.Values()) // Output: Slice: [4 5]
    }

    union.SortBy(func(a, b int) bool {
        return a > b
    })
    fmt.Println("Sorted descending:", union.Values()) // Output: Sorted descending: [5 4 3]
}
```