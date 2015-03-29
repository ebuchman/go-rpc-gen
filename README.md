# RPC-Gen

`rpc-gen` auto generates an rpc client using the same API as the core functions running on the server.

eg. `rpc-gen -interface Client -pkg core -type *ClientHTTP,*ClientJSON -exclude pipe.go -out-pkg rpc`

will make a new interface `Client`, with all the exported methods from the package `core` but excluding the files `pipe.go`. 
Two implementations of the interface are generated in this case, one on `*ClientHTTP` and one on `*ClientJSON`.
The programs author is required to provide one rpc function template for each type, which `rpc-gen` will autocomplete.

Run the above command in the example directory and examine the output.
