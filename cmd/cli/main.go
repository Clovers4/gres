package main

func main() {
	cli := NewClient()
	defer cli.GracefulExit()
	cli.Interact()
}
