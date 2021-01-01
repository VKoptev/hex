package hex

import (
	"encoding/json"
	"testing"
)

type TestStruct struct {
	Hex Hex `json:"hex"`
}

type a struct {
	A string `json:"a"`
}

type TestMapStruct struct {
	Map map[Hex]a `json:"map"`
}

type testMapSource struct {
	Map map[string]a `json:"map"`
}

func (t *TestMapStruct) UnmarshalJSON(data []byte) error {
	m := new(testMapSource)
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	t.Map = make(map[Hex]a, len(m.Map))
	for hh, _a := range m.Map {
		var h Hex
		if err := json.Unmarshal([]byte(hh), &h); err != nil {
			return err
		}
		t.Map[h] = _a
	}
	return nil
}

func (t TestMapStruct) MarshalJSON() ([]byte, error) {
	m := testMapSource{
		Map: make(map[string]a, len(t.Map)),
	}

	for h, _a := range t.Map {
		hh, err := h.MarshalJSON()
		if err != nil {
			return nil, err
		}
		m.Map[string(hh)] = _a
	}

	return json.Marshal(m)
}

func TestHex_MarshalJSON(t *testing.T) {
	t.Parallel()
	bb, err := json.Marshal(TestStruct{Hex: New(0, 10)})
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}
	if string(bb) != `{"hex":[0,10]}` {
		t.Errorf("wrong marshaling: %s", bb)
	}
}

func TestHex_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	ts := new(TestStruct)
	if err := json.Unmarshal([]byte(`{"hex":[0,10]}`), ts); err != nil {
		t.Fatalf("error in test: %v", err)
	}
	if ts.Hex != New(0, 10) {
		t.Errorf("wrong unmarshaling: %+v", ts)
	}
}

func TestMap_MarshalJSON(t *testing.T) {
	t.Parallel()
	bb, err := json.Marshal(TestMapStruct{Map: map[Hex]a{
		New(0, 10): {A: "test"},
	}})
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}
	if string(bb) != `{"map":{"[0,10]":{"a":"test"}}}` {
		t.Errorf("wrong marshaling: %s", bb)
	}
}

func TestMap_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tms := new(TestMapStruct)
	if err := json.Unmarshal([]byte(`{"map":{"[0,10]":{"a":"test"}}}`), tms); err != nil {
		t.Fatalf("error in test: %v", err)
	}
	h := New(0, 10)
	if _, ok := tms.Map[h]; !ok {
		t.Errorf("wrong unmarshaling: %+v", tms)
	}
	if tms.Map[h].A != "test" {
		t.Errorf("wrong unmarshaling: %+v", tms)
	}
}
