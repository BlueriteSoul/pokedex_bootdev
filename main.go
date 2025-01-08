package main

import (
	"bufio"
	//"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var Pokedex map[string]mapJSONstructPokemon

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *Cache, string) error
}
type config struct {
	next         string
	prev         string
	pokeEndPoint string
	limit        int
	offset       int
}
type mapJSONstruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
type mapJSONstructExplore struct {
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
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
type mapJSONstructPokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	cached   map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

func newCache(dur time.Duration) *Cache {
	result := Cache{
		cached:   make(map[string]cacheEntry),
		mu:       sync.RWMutex{},
		interval: dur,
	}
	result.reapLoop()
	return &result
}
func (c *Cache) add(key string, val []byte) {
	c.mu.Lock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.cached[key] = entry
	c.mu.Unlock()
}
func (c *Cache) get(key string) ([]byte, bool) {
	c.mu.Lock()
	result, ok := c.cached[key]
	c.mu.Unlock()
	return result.val, ok
}
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			<-ticker.C
			c.mu.Lock()
			for k, entry := range c.cached { // get both key and entry
				if time.Since(entry.createdAt) > c.interval { // time.Since is cleaner than Now().Sub
					delete(c.cached, k) // delete using the key
				}
			}
			c.mu.Unlock()
		}
	}()
}

// the return is unreachable, but the assingment says this is the
// function signiture we use
func commandExit(*config, *Cache, string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(*config, *Cache, string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}

	//fmt.Println("help: Displays a help message")
	//fmt.Println("exit: Exit the Pokedex")
	return nil
}
func commandMap(myconfig *config, c *Cache, notGonnaUse string) error {
	value, ok := c.cached[myconfig.next]
	var mapJSON mapJSONstruct
	if ok {
		if err := json.Unmarshal(value.val, &mapJSON); err != nil {
			return err
		}
	} else {
		res, err := http.Get(myconfig.next)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&mapJSON); err != nil {
			return err
		}

		// Marshal the `mapJSON` back to bytes for caching
		data, err := json.Marshal(mapJSON)
		if err != nil {
			return err
		}
		c.add(myconfig.next, data)

	}
	myconfig.prev = myconfig.next
	myconfig.next = mapJSON.Next
	// Process `mapJSON.Results` outside of if-else
	for _, area := range mapJSON.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(myconfig *config, c *Cache, notGonnaUse string) error {
	fmt.Println("Checking cache key:", myconfig.prev) //nbjfb jksfbjvkfbjk
	value, ok := c.cached[myconfig.prev]
	var mapJSON mapJSONstruct
	fmt.Println("Cache hit:", ok) //sfjkgbjdfbjkdf
	if ok {
		fmt.Println("Using cached data for:", myconfig.prev) //djkbvjkdbvjkd

		if err := json.Unmarshal(value.val, &mapJSON); err != nil {
			return err
		}
	} else {
		fmt.Println("Fetching new data for:", myconfig.prev) //dkjgbrjkvbjkd
		res, err := http.Get(myconfig.prev)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&mapJSON); err != nil {
			return err
		}

		// Marshal the `mapJSON` back to bytes for caching
		data, err := json.Marshal(mapJSON)
		if err != nil {
			return err
		}
		c.add(myconfig.prev, data)

	}
	myconfig.prev = mapJSON.Previous
	myconfig.next = mapJSON.Next
	// Process `mapJSON.Results` outside of if-else
	for _, area := range mapJSON.Results {
		fmt.Println(area.Name)
	}

	return nil
	/*
		res, err := http.Get(myconfig.prev)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		//fmt.Println(response)

		var mapJSON mapJSONstruct
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&mapJSON); err != nil {
			return err
		}
		myconfig.prev = mapJSON.Previous
		myconfig.next = mapJSON.Next
		for _, area := range mapJSON.Results {
			fmt.Println(area.Name)
		}
		return nil
	*/
}
func commandExplore(myconfig *config, c *Cache, area string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + area + "/"
	value, ok := c.cached[url]
	var mapJSON mapJSONstructExplore
	fmt.Println("Cache hit:", ok) //sfjkgbjdfbjkdf
	if ok {
		fmt.Println("Using cached data for:", url) //djkbvjkdbvjkd

		if err := json.Unmarshal(value.val, &mapJSON); err != nil {
			return err
		}
	} else {
		fmt.Println("Fetching new data for:", url) //dkjgbrjkvbjkd
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&mapJSON); err != nil {
			return err
		}

		// Marshal the `mapJSON` back to bytes for caching
		data, err := json.Marshal(mapJSON)
		if err != nil {
			return err
		}
		c.add(url, data)

	}

	// Process `mapJSON.Results` outside of if-else
	for _, pokemon := range mapJSON.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil

}
func normalizeXP(xp int) float64 {
	// Ensure that xp is within the range [1, 500]
	if xp < 1 {
		xp = 1
	} else if xp > 500 {
		xp = 500
	}

	// Scale the value between 0.3 and 1, compressing larger values more
	const minXP = 1.0
	const maxXP = 500.0
	const minNorm = 0.3
	const maxNorm = 1.0

	// Logarithmic scale compression
	logScale := (maxNorm-minNorm)*(math.Log(float64(xp))-math.Log(minXP))/(math.Log(maxXP)-math.Log(minXP)) + minNorm

	return logScale
}
func commandCatch(myconfig *config, c *Cache, pokeName string) error {
	url2 := myconfig.pokeEndPoint + pokeName + "/"
	value, ok := c.cached[url2]
	var mapJSON mapJSONstructPokemon
	fmt.Println("Cache hit:", ok, "I am also where I'm supposed to be") //sfjkgbjdfbjkdf
	if ok {
		fmt.Println("Using cached data for:", url2) //djkbvjkdbvjkd

		if err := json.Unmarshal(value.val, &mapJSON); err != nil {
			return err
		}
	} else {
		fmt.Println("Fetching new data for:", url2) //dkjgbrjkvbjkd
		res, err := http.Get(url2)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&mapJSON); err != nil {
			return err
		}

		// Marshal the `mapJSON` back to bytes for caching
		data, err := json.Marshal(mapJSON)
		if err != nil {
			return err
		}
		c.add(url2, data)

	}

	fmt.Printf("Throwing a Pokeball at %s...\n", mapJSON.Name)
	if normalizeXP(mapJSON.BaseExperience)*float64(rand.Intn(100)) > 70 {
		Pokedex[mapJSON.Name] = mapJSON
		fmt.Println(mapJSON.Name, " was caught!")
	} else {
		fmt.Println(mapJSON.Name, " escaped!")
	}

	return nil

}
func commandInspect(myconfig *config, c *Cache, pokeName string) error {
	pokemon, ok := Pokedex[pokeName]
	if !ok {
		fmt.Println("Pokemon not found in Pokedex")
		return nil
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Println(stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typee := range pokemon.Types {
		fmt.Println("-", typee.Type.Name)
	}
	return nil
}
func commandPokedex(myconfig *config, c *Cache, pokeName string) error {
	if len(Pokedex) == 0 {
		fmt.Println("You didn't catch any pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pok := range Pokedex {
		fmt.Println("-", pok.Name)
	}
	return nil
}
func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.ToLower(strings.TrimSpace(word))
	}
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}, "help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}, "map": {
		name:        "map",
		description: "Shows next 20 available locations",
		callback:    commandMap,
	}, "mapb": {
		name:        "mapb",
		description: "Shows previous 20available locations",
		callback:    commandMapb,
	}, "explore": {
		name:        "explore",
		description: "Shows pokemon available in the area e.g. 'explore pastoria-city-area",
		callback:    commandExplore,
	}, "catch": {
		name:        "catch",
		description: "Attempt to catch a pokemon e.g. catch squirtle",
		callback:    commandCatch,
	}, "inspect": {
		name:        "inspect",
		description: "Describes the pokemon e.g. inspect pikachu",
		callback:    commandInspect,
	}, "pokedex": {
		name:        "pokedex",
		description: "Displays pokemons you have caught",
		callback:    commandPokedex,
	}}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	Pokedex = make(map[string]mapJSONstructPokemon)
	fmt.Print("Pokedex >")
	arg := ""
	commands := getCommands()
	cache := newCache(5 * time.Second)
	//the API has 1054 results, I limit the query to 17, because 1054 % 17 == 0
	//this ensures I will always get 17 results, even if I go to the end
	//that would glich out, if travelling in multiples of 20
	myConfig := config{next: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=17", //to test the far end ?offset=1022&limit=17
		prev: "", pokeEndPoint: "https://pokeapi.co/api/v2/pokemon/", offset: 1000, limit: 20}
	for scanner.Scan() {

		input := scanner.Text()
		words := cleanInput(input)
		command := words[0]
		if len(words) > 1 {
			arg = words[1]
		}

		value, ok := commands[command]
		if !ok {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex >")
			continue
		}
		err := value.callback(&myConfig, cache, arg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("Pokedex >")

	}
}
