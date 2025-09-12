package parser

import (
	"os"
	"parttwo/processor/lexer"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONToMap(t *testing.T) {
	cases := []struct {
		Name          string
		InputPath     string
		Expected      any
		Out           any
		ExpectedError bool
	}{
		{
			Name:      "map[string]any with number (positive and negative), string, boolean, array, and nested object",
			InputPath: "./tests/test_map_complete_success.json",
			Expected: map[string]any{
				"key":   "value",
				"key-n": int64(101),
				"key-o": map[string]any{
					"test":  int64(-202),
					"float": float64(3.14),
					"aaa":   float64(3.1423),
				},
				"key-l": []any{int64(1), int64(2), float64(3.5), int64(40), float64(2.333333333)},
				"test":  int64(-202),
				"float": float64(3.14),
				"aaa":   float64(3.1423),
				"array-object": []any{
					map[string]any{
						"id":   int64(1),
						"desc": "bejirrr",
					},
					map[string]any{
						"id":   int64(2),
						"desc": "walaweehhh",
					},
					map[string]any{
						"id":       int64(3),
						"desc":     "awokdawodkawodawodk",
						"apatuhhh": "apanihhhh",
					},
					map[string]any{
						"id":         int64(4),
						"ddwadawdaw": "aokdawodkoa",
						"test":       "aaaaaaaa",
					},
				},
			},
			ExpectedError: false,
			Out:           &map[string]any{},
		},
		{
			Name:      "map[string]float64 with all json value being key number pair: success",
			InputPath: "./tests/test_map_string_number.json",
			Expected: map[string]float64{
				"x1": 1.00,
				"x2": 66.0,
				"y1": 2.40,
				"y2": 22.00,
			},
			ExpectedError: false,
			Out:           &map[string]float64{},
		},
		{
			Name:          "map[string]float64 error due to one of the json value not a number",
			InputPath:     "./tests/test_map_string_not_all_number.json",
			Expected:      map[string]float64{},
			ExpectedError: true,
			Out:           &map[string]float64{},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			input, err := os.ReadFile(c.InputPath)
			assert.Nil(t, err)
			l := lexer.New(string(input))
			l.Process()
			p := New(&l)
			p.Process()
			err = p.Decode(c.Out)
			if c.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, c.Expected, reflect.ValueOf(c.Out).Elem().Interface())
		})
	}

}

func TestJSONToArray(t *testing.T) {
	cases := []struct {
		Name        string
		InputPath   string
		Expected    any
		ExpectedErr error
		Out         any
	}{
		{
			Name:        "success",
			InputPath:   "./tests/test_array_same_type.json",
			Expected:    []string{"test", "test", "test", "test"},
			ExpectedErr: nil,
			Out:         &[]string{},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			input, err := os.ReadFile(c.InputPath)
			assert.Nil(t, err)
			l := lexer.New(string(input))
			l.Process()
			p := New(&l)
			p.Process()
			err = p.Decode(c.Out)
			assert.Equal(t, c.ExpectedErr, err)
		})
	}
}

func TestJSONToStruct(t *testing.T) {
	cases := []struct {
		name        string
		inputPath   string
		expected    any
		expectedErr error
		out         any
	}{
		{
			name:      "parse to struct success",
			inputPath: "./tests/test_struct_success.json",
			expected: Person{
				Name:    "arief",
				Age:     12,
				Balance: 200,
				Education: struct {
					InstitutionName string "json:\"institution_name\""
					Degree          string "json:\"degree\""
				}{
					InstitutionName: "apa cobaaa",
					Degree:          "apa hayooo",
				},
				Scores: []int{2, 2, 4},
			},
			expectedErr: nil,
			out:         &Person{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input, err := os.ReadFile(c.inputPath)
			assert.Nil(t, err)
			l := lexer.New(string(input))
			l.Process()
			p := New(&l)
			p.Process()
			err = p.Decode(c.out)
			assert.Equal(t, c.expectedErr, err)
			assert.Equal(t, c.expected, reflect.ValueOf(c.out).Elem().Interface())
		})
	}
}

type Person struct {
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Balance   float64 `json:"balance"`
	Education struct {
		InstitutionName string `json:"institution_name"`
		Degree          string `json:"degree"`
	} `json:"current_education"`
	Scores []int `json:"scores"`
}
