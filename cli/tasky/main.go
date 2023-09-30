package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github/tasky"
	"io"
	"os"
	"strings"
)

const taskFile = ".tasky.json"

func main() {
	// Define command-line flags
	add := flag.Bool("add", false, "add new task")
	complete := flag.Int("complete", 0, "task completed")
	rm := flag.Int("rm", 0, "task removed successfully")
	list := flag.Bool("list", false, "list all tasks")
	flag.Parse()

	tasks := &tasky.Todos{}

	// Load tasks from file
	if err := tasks.Load(taskFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			handleError(err)
		}

		tasks.Add(task)
		err = tasks.Store(taskFile)
		handleError(err)

	case *complete > 0:
		err := tasks.Complete(*complete)
		handleError(err)

		err = tasks.Store(taskFile)
		handleError(err)

	case *rm > 0:
		err := tasks.Delete(*rm)
		handleError(err)

		err = tasks.Store(taskFile)
		handleError(err)

	case *list:
		tasks.Print()

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("your task is empty :)")
	}

	return text, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
