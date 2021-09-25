package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type AllValidatorsStruct struct {
	One    string   `validate:"len:6"`
	Two    string   `validate:"regexp:^\\d+$"`
	Three  UserRole `validate:"in:a,b"`
	Four   string   `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
	Five   []string `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
	Six    string
	Seven  int   `validate:"min:1"`
	Eight  int   `validate:"max:50"`
	Nine   int   `validate:"in:0,1,2"`
	Ten    int   `validate:"min:1|in:0,25,51|max:50"`
	Eleven []int `validate:"min:1|in:0,1,25,50,51|max:50"`
	Twelve int
}

func TestValidate(t *testing.T) {
	notStruct := 10

	validStruct := AllValidatorsStruct{
		One:    "123123",
		Two:    "111",
		Three:  "a",
		Four:   "111",
		Five:   []string{"111", "222"},
		Six:    "20",
		Seven:  1,
		Eight:  50,
		Nine:   1,
		Ten:    25,
		Eleven: []int{1, 50},
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{notStruct, ErrNotStruct},
		{validStruct, nil},
		{struct {
			Field string `validate:"len:6"`
		}{"12"}, ErrStructValidation},
		{struct {
			Field string `validate:"regexp:^\\d+$"`
		}{"qw"}, ErrStructValidation},
		{struct {
			Field UserRole `validate:"in:a,b"`
		}{"c"}, ErrStructValidation},
		{struct {
			Field string `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
		}{"a"}, ErrStructValidation},
		{struct {
			Field string `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
		}{"33"}, ErrStructValidation},
		{struct {
			Field []string `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
		}{[]string{"111", "a"}}, ErrStructValidation},
		{struct {
			Field []string `validate:"len:3|regexp:^\\d+$|in:a,111,222,33"`
		}{[]string{"555", "33"}}, ErrStructValidation},
		{struct {
			Field int `validate:"min:1"`
		}{0}, ErrStructValidation},
		{struct {
			Field int `validate:"max:50"`
		}{51}, ErrStructValidation},
		{struct {
			Field int `validate:"in:0,1,2"`
		}{-3}, ErrStructValidation},
		{struct {
			Field int `validate:"min:1|in:0,25,51|max:50"`
		}{0}, ErrStructValidation},
		{struct {
			Field int `validate:"min:1|in:0,25,51|max:50"`
		}{51}, ErrStructValidation},
		{struct {
			Field []int `validate:"min:1|in:0,1,25,50,51|max:50"`
		}{[]int{1, 25, 51}}, ErrStructValidation},
		{struct {
			Field []int `validate:"min:1|in:0,1,25,50,51|max:50"`
		}{[]int{2, 49}}, ErrStructValidation},
		{struct {
			Field []int `validate:"min:1|in:0,1,25,50,51|something:50"`
		}{[]int{2, 49}}, ErrProgram},
		{struct {
			Field []int `validate:"min:1|in0,1,25,50,51|max:50"`
		}{[]int{2, 49}}, ErrProgram},
		{struct {
			Field []int `validate:"min:a"`
		}{[]int{2, 49}}, ErrProgram},
		{struct {
			Field []int `validate:"max:a"`
		}{[]int{2, 49}}, ErrProgram},
		{struct {
			Field []int `validate:"in:a,b,c"`
		}{[]int{2, 49}}, ErrProgram},
		{struct {
			Field string `validate:"regexp:[["`
		}{"a"}, ErrProgram},
		{struct {
			Field bool `validate:"something:true"`
		}{true}, ErrProgram},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Equal(t, nil, result)
			} else {
				require.True(t, errors.Is(result, tt.expectedErr))
			}
		})
	}
}
