package main

import (
	"github.com/clovers4/gres"
)

func main() {
	// init and start server
	srv := gres.NewServer(gres.DbnumOption(8))
	srv.ListenAndServe()
	defer srv.Stop()
}
