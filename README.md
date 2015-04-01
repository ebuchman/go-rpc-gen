# RPC-Gen

`go-rpc-gen` auto generates an rpc client using the same API as the core functions running on the server.

So if you have a function called by your rpc handlers that looks like:

```
package core

type Bar struct{
	Name string
	Age int
}

type Baz struct{
	Id []byte
	Gopher bool
}

func Foo(b *Bar) (*Baz, error){
	// do something
}
```

then `go-rpc-gen` will generate a client that looks like 

```
package client

func Foo(b *core.Bar) (*core.Baz, error){
	// call core.Foo(b) over rpc
}

```

where the functionality for calling over rpc is specified in comments using an extremely simplified templating language.

A more thorough example is provided in the `example` directory. See `example/client.go` for `go-rpc-gen` directives and the templates for the client functions.
The generated methods are in `client_methods.go`.

eg. `go-rpc-gen -interface Client -pkg core -type *ClientHTTP,*ClientJSON -exclude pipe.go -out-pkg rpc`

will make a new interface `Client`, with all the exported methods from the package `core` but excluding the files `pipe.go`. 
Two implementations of the interface are generated in this case, one on `*ClientHTTP` and one on `*ClientJSON`.
The programs author is required to provide one rpc function template for each type, which `rpc-gen` will autocomplete.

Run the above command in the example directory and examine the output (`client_methods.go`). Alternatively, just run `go generate`.
The API generated is the same as that found in `examples/core`. Everything else is filler for the rpc mechanism,
but `rpc-gen` is relatively agnostic.
