package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("\nStarting ZDE shell...")

	// handle SIGINT (Ctrl+C) to gracefully exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		for range sigChan {
			fmt.Println("\nReceived interrupt. Exiting ZDE shell...")
			os.Exit(0)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		userInfo, err := getUserInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		fmt.Print(userInfo)

		// read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		// handle the input execution
		if err = execInput(input); err != nil {
			fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		}
	}
}

func getUserInfo() (string, error) {
	dir, err := os.Getwd() // get current working directory
	if err != nil {
		return "", err
	}

	hostName, err := os.Hostname() // get machine's host name
	if err != nil {
		return "", err
	}

	currentUser, err := user.Current() // get the current user
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\n%s@%s %s> ", currentUser.Username, hostName, dir), nil
}

func execInput(input string) error {
	input = strings.TrimSpace(input) // remove the newline character and any trailing spaces

	args := strings.Fields(input) // split the input to separate the command and the arguments

	// check if no command was entered
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "cd":
		return changeDirectory(args)
	case "exit":
		fmt.Println("Exiting ZDE shell...")
		os.Exit(0)
	}

	// Prepare the command to execute
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func changeDirectory(args []string) error {
	if len(args) < 2 { // navigate to home directory if no  path is provided
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(homeDir)
	}
	return os.Chdir(args[1]) // navigate to the provided path
}
