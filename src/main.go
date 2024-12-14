package main

import (
	"clli_ld/src/myutils"
	"fmt"
	"time"

	// "flag"
	"os"
	"strconv"
	// import the config package
)
var backendStatus chan bool

func main() {
	backendStatus = make(chan bool,1)
	// load the config
	config , err := myutils.LoadConfig()

	// test the backend
	go myutils.TestBackend(config,backendStatus)

	if(err != nil) {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if(len(os.Args) < 2) {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	// take in the first argument as the command
	cmd := os.Args[1]
	fmt.Println("Command: ", cmd)

	switch cmd {
	case "add":
		a, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ", err)
			// exitt the program
			os.Exit(1)
		}
		b, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Error: ", err)
			// exitt the program
			os.Exit(1)
		}
		// parse the a b to int and pass to add function
		fmt.Println("a + b = ", add(a, b))

	case "factorial":
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ", err)
			// exitt the program
			os.Exit(1)
		}
		// parse the a b to int and pass to add function
		fmt.Println("factorial of ", n, " is ", factorial(n))


	case "listencpp" : 
		 delay := 1000
		myutils.CompileAndListenCpp(os.Args[2],&delay)
	

	default:
		fmt.Println("Command not found")
	}
	time.Sleep(time.Second * 6)
	fmt.Println("Backend status: ", <-backendStatus)

}

// define a function to add two numbers
func add(a, b int) int {
	return a + b
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return (n * factorial(n-1))%10000007
}
