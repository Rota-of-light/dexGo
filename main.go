package main

import (
    "fmt" 
    "strings"
)

func cleanInput(test string) []string {
    return strings.Fields(strings.ToLower(test))
} 

func main() {
    fmt.Printf("Hello, World!")
}
