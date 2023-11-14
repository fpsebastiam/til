Performing basic profiling of CPU usage of a Go program is straightforward.
First you expose some endpoints in your code:

```go
import (
    "log"
    _ "net/http/pprof"
    "net/http"
)

func main(){
    go func(){
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }

    // Program goes here

    select {} // keep it running
}
```

One can then use an interactive shell through `go tool pprof http://localhost:6060/debug/pprof/profile`
to see the profiling. Some useful commands are `top`, `top n`, and `list your-function-name`.

For memory, check the `/debug/pprof/heap` endpoint through `go tool pprof`

A nice summary can be found [here](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)
