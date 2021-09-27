package json

import (
	"encoding/json"

	"github.com/VKoptev/hex"
)

// Hex embeds hex.Hex to implement json.Marshaler and json.Unmarshaler interfaces.
type Hex struct {
	hex.Hex
}

// MarshalJSON returns JSON representation of Hex.
func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Array())
}

// UnmarshalJSON fills Hex from JSON representation.
func (h *Hex) UnmarshalJSON(data []byte) error {
	var a [2]int
	err := json.Unmarshal(data, &a)

	if err != nil {
		return err
	}

	h.Hex = hex.NewWithArray(a)

	return nil
}
