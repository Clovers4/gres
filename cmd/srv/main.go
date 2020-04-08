package main

import (
	"github.com/clovers4/gres"
)

func main() {
	// init and start server
	srv := gres.NewServer()
	srv.ListenAndServe()
	defer srv.Stop()
}
