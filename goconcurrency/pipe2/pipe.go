package main

import (
	"fmt"
	"bytes"
	"log"
	"os/exec"
)

func main() {
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "bash")

	var outputBuf1 bytes.Buffer
	cmd1.Stdout = &outputBuf1

	if err := cmd1.Start(); err != nil {
		log.Fatalf("Error: The first command can not be startup: %s\n", err)
	}

	if err := cmd1.Wait(); err != nil {
		log.Fatalf("Error: Couldn't wait for the first command: %s\n", err)
	}

	cmd2.Stdin = &outputBuf1
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		log.Fatalf("Error: The second command can not be startup: %s\n", err)
	}

	if err := cmd2.Wait(); err != nil {
		log.Fatalf("Error: Couldn't wait for the second command: %s\n", err)
	}

	fmt.Printf("%s\n", outputBuf2.Bytes())
}