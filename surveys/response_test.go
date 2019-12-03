package surveys

import (
	"fmt"
	"regexp"
	"testing"
)

var reRangeError = regexp.MustCompile("Must be between [0-9]+ and [0-9]+")

func TestResponseValidate(t *testing.T) {

	for testIndex, tc := range []struct {
		input  Response
		expect map[string]*regexp.Regexp
	}{
		{input: Response{}, expect: nil},
		{
			input: Response{Age: -1},
			expect: map[string]*regexp.Regexp{
				"age": reRangeError,
			},
		},
	} {

		t.Run(fmt.Sprintf("case %d", testIndex), func(t *testing.T) {
			got := tc.input.Validate()
			for key, reMatcher := range tc.expect {
				if gotVal, gotOk := got[key]; !gotOk {
					t.Errorf("Missing key %s in response", key)
				} else if !reMatcher.MatchString(gotVal) {
					t.Errorf("At key %s, got %s", key, gotVal)
				}
			}
			for key := range got {
				if _, ok := tc.expect[key]; !ok {
					t.Errorf("Extra key %s in response", key)
				}
			}
		})
	}

}
