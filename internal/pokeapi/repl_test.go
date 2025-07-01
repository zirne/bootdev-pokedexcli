package pokeapi

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		ep      string
		success bool
	}{
		{
			ep:      "https://pokeapi.co/api/v2/location-area/",
			success: true,
		},
		{
			ep:      "https://pokeapi.co/api/v2/gurka/",
			success: false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			client := NewClient()
			val, err := client.Get(c.ep)
			if c.success {
				if err != nil {
					t.Errorf("Expected success. Got error %v with data: %v", err, val)
				}
			} else {
				if err == nil {
					t.Errorf("Expected failure. Got data '%v'", val)
				}
			}
		})
	}
}
