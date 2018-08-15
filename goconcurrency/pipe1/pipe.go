package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang")

	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 cannot be startup: %s\n", err)
		return
	}

	output0 := make([]byte, 40)
	n, err := stdout0.Read(output0)
	if err != nil {
		fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
		return
	}

	fmt.Printf("%s\n", output0[:n])
}