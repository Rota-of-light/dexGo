package main

import (
    "fmt" 
    "strings"
    "bufio"
    "os"
    "net/http"
    "encoding/json"
    "io"
    "time"
    "github.com/Rota-of-light/dexGo/internal/pokecache"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*Config, string) error
}

type Config struct {
    Next     *string
    Previous *string
    cache *pokecache.Cache
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

type Encounters struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExit(cfg *Config, locate string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

var commands map[string]cliCommand

func commandHelp(cfg *Config, locate string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Printf("Usage:\n\n")
    for _, command := range commands {
        fmt.Printf("%v: %v\n", command.name, command.description)
    }
    return nil
}

func commandMap(cfg *Config, locate string) error {
    var httpString string
    var response Respond
    if cfg.Next == nil {
        httpString = "https://pokeapi.co/api/v2/location-area"
    } else {
        httpString = *cfg.Next
    }
    if cachedData, found := cfg.cache.Get(httpString); found {
        if err := json.Unmarshal(cachedData, &response); err != nil {
            return err
        }
    } else {
        res, err := http.Get(httpString)
        if err != nil {
            return err
        }
        defer res.Body.Close()
        
        body, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
        
        cfg.cache.Add(httpString, body)

        if err := json.Unmarshal(body, &response); err != nil {
            return err
        }
    }
    for _, area := range response.Results {
        fmt.Printf("%v\n", area.Name)
    }

    cfg.Previous = response.Previous
    cfg.Next = &response.Next
    return nil
}

func commandMapb(cfg *Config, locate string) error {
    var httpString string
    var response Respond
    if cfg.Previous == nil {
        fmt.Println("you're on the first page")
        return nil
    }
    httpString = *cfg.Previous
    if cachedData, found := cfg.cache.Get(httpString); found {
        if err := json.Unmarshal(cachedData, &response); err != nil {
            return err
        }
    } else {
        res, err := http.Get(httpString)
        if err != nil {
            return err
        }
        defer res.Body.Close()
        
        body, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
                
        cfg.cache.Add(httpString, body)

        if err := json.Unmarshal(body, &response); err != nil {
            return err
        }
    }
    for _, area := range response.Results {
        fmt.Printf("%v\n", area.Name)
    }

    cfg.Previous = response.Previous
    cfg.Next = &response.Next
    return nil
}

func commandExplore(cfg *Config, locate string) error {
    var response Encounters
    if locate == "" {
        fmt.Println("No location given.")
        return nil
    }
    httpString := "https://pokeapi.co/api/v2/location-area/" + locate
    if cachedData, found := cfg.cache.Get(httpString); found {
        if err := json.Unmarshal(cachedData, &response); err != nil {
            return err
        }
    } else {
        res, err := http.Get(httpString)
        if res.StatusCode != 200 {
            return fmt.Errorf("Location area '%s' not found", locate)
        }
        if err != nil {
            return err
        }
        defer res.Body.Close()
        
        body, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
                
        cfg.cache.Add(httpString, body)

        if err := json.Unmarshal(body, &response); err != nil {
            return err
        }
    }
    fmt.Println("Exploring " + locate + "...")
    fmt.Println("Found Pokemon:")
    for _, encounter := range response.PokemonEncounters {
        fmt.Println(" - " + encounter.Pokemon.Name)
    }
    return nil
}

func cleanInput(test string) []string {
    return strings.Fields(strings.ToLower(test))
} 

func main() {
    config := Config{
        Next: nil,
        Previous: nil,
        cache: pokecache.NewCache(5 * time.Minute),
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
        "explore": {
            name:        "explore",
            description: "Displays a location's pokemon",
            callback:    commandExplore,
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
            var locate string
            if len(cleanedInput) > 1 {
                locate = cleanedInput[1]
            }
            err := cmd.callback(&config, locate)
            if err != nil {
                fmt.Println("Error is:", err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}
