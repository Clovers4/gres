package main


func main() {
	cli := NewClient()

	cli.Interact()
	defer cli.GracefulExit()
}
