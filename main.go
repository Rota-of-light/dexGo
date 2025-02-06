package main

import (
    "fmt" 
    "strings"
    "bufio"
    "os"
    "net/http"
    "encoding/json"
    "io"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*Config) error
}

type Config struct {
    Next     *string
    Previous *string
}

type Respond struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

//var response Respond

func commandExit(cfg *Config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

var commands map[string]cliCommand

func commandHelp(cfg *Config) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Printf("Usage:\n\n")
    for _, command := range commands {
        fmt.Printf("%v: %v\n", command.name, command.description)
    }
    return nil
}

func commandMap(cfg *Config) error {
    var httpString string
    var response Respond
    if cfg.Next == nil {
        httpString = "https://pokeapi.co/api/v2/location-area"
    } else {
        httpString = *cfg.Next
    }
    res, err := http.Get(httpString)
    if err != nil {
		return err
	}
    defer res.Body.Close()
    
    body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

    if err := json.Unmarshal(body, &response); err != nil {
        return err
    }
    for _, area := range response.Results {
        fmt.Printf("%v\n", area.Name)
    }

    cfg.Previous = response.Previous
    cfg.Next = &response.Next
    return nil
}

func commandMapb(cfg *Config) error {
    var httpString string
    var response Respond
    if cfg.Previous == nil {
        fmt.Println("you're on the first page")
        return nil
    }
    httpString = *cfg.Previous
    res, err := http.Get(httpString)
    if err != nil {
		return err
	}
    defer res.Body.Close()
    
    body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

    if err := json.Unmarshal(body, &response); err != nil {
        return err
    }
    
    for _, area := range response.Results {
        fmt.Printf("%v\n", area.Name)
    }

    cfg.Previous = response.Previous
    cfg.Next = &response.Next
    return nil
}

func cleanInput(test string) []string {
    return strings.Fields(strings.ToLower(test))
} 

func main() {
    config := Config{
        Next: nil,
        Previous: nil,
    }
    commands = map[string]cliCommand{
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "map": {
        name:        "map",
        description: "Shows the next 20 location areas",
        callback:    commandMap,
        },
        "mapb": {
        name:        "mapb",
        description: "Shows the previous 20 location areas",
        callback:    commandMapb,
        },
    }
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        userInput := scanner.Text()
        cleanedInput := cleanInput(userInput)
        if len(cleanedInput) == 0 {
            continue
        }
        cmd, ok := commands[cleanedInput[0]]
        if ok {
            err := cmd.callback(&config)
            if err != nil {
                fmt.Println("Error is:", err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}
