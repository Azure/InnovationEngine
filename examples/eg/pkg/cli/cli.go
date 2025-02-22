package cli

import (
    "fmt"
    "os"
)

func ExecuteCommand(command string) {
    switch command {
    case "hello":
        HelloCommand()
    default:
        fmt.Println("Unknown command:", command)
        os.Exit(1)
    }
}

func HelloCommand() {
    fmt.Println("Hello, World!")
}