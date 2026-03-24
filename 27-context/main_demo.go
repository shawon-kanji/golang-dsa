package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunAllDemos provides an interactive menu to run different context demos
func RunAllDemos() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("Go Context Package - Interactive Demo")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1.  context.WithCancel - Manual cancellation")
		fmt.Println("2.  context.WithDeadline - Absolute deadline")
		fmt.Println("3.  context.WithTimeout - Relative timeout (original example)")
		fmt.Println("4.  context.WithValue - Passing values")
		fmt.Println("5.  context.WithCancelCause - Cancel with reason (Go 1.20+)")
		fmt.Println("6.  context.AfterFunc - Cleanup callbacks (Go 1.21+)")
		fmt.Println("7.  context.WithoutCancel - Detached context (Go 1.21+)")
		fmt.Println("8.  HTTP-Style Context Usage")
		fmt.Println("9.  Parallel Operations with Context")
		fmt.Println("10. Context Propagation")
		fmt.Println("11. Graceful Shutdown Pattern")
		fmt.Println("0.  Exit")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			DemoCancel()
		case "2":
			DemoDeadline()
			DemoDeadlineSuccess()
		case "3":
			// Original example - runs the slowOperation from main.go
			fmt.Println("\n=== context.WithTimeout Demo (Original) ===")
			fmt.Println("See main.go for the original timeout example")
			fmt.Println("Run: go run main.go")
		case "4":
			DemoValue()
			DemoValueBestPractices()
		case "5":
			DemoCancelCause()
			DemoTimeoutWithCause()
		case "6":
			DemoAfterFunc()
			DemoAfterFuncWithStop()
			DemoAfterFuncMultiple()
		case "7":
			DemoWithoutCancel()
			DemoWithoutCancelUseCase()
		case "8":
			DemoHTTPStyleContext()
		case "9":
			DemoParallelWithContext()
		case "10":
			DemoPropagation()
		case "11":
			DemoGracefulShutdown()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}

		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}

// Uncomment the following to make this the main entry point:
// func main() {
// 	RunAllDemos()
// }
