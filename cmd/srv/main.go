package main

import (
	"github.com/clovers4/gres"
)

func main() {
	// init and start server
	srv := gres.NewServer()
	defer srv.Stop()
	srv.ListenAndServe()
}
