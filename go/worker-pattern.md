This is quite useful for data processing. It limits the number
of go routines and fits nicely with parallel computation paradigms

```go
package main

import (
    "fmt"
    "sync"
)

func processData(n int, wc <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	nSquared := n * n
    fmt.Println(nSquared)
	<-wc
}

func main() {
    data := make([]int, 1000)
    for i:= range data {
        data[i] = i + 1
    }

    var wg sync.WaitGroup
    maxWorkers := 5
    workerChannel := make(chan struct{}, maxWorkers)

    for _, data := range data {
        wg.Add(1)
        workerChannel <- struct{}{}
        go processData(data, workerChannel, &wg)
    }

    wg.Wait()
}
```
