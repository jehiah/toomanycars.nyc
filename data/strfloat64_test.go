package parkingdata

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStrFloat64_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		input    string
		expected float64
	}
	tests := []testCase{
		{
			input:    `"1.1"`,
			expected: 1.1,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			var s StrFloat64
			err := json.Unmarshal([]byte(tc.input), &s)
			if err != nil {
				t.Fatal(err)
			}
			if float64(s) != tc.expected {
				t.Errorf("got %f expected %f for %q", s, tc.expected, tc.input)
			}
		})
	}
}
