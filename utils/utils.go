package utils

type Config struct {
	PREV_URL string
	NEXT_URL string
}

type GetLocationsResponse struct {
	Count    int
	Next     string
	Previous any
	Results  []struct {
		Name string
		URL  string
	}
}

func Map(config *Config) error {

}

func Mapb(config *Config) {

}
