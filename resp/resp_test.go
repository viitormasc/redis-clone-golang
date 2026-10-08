package resp_test

import (
	"redis-clone/resp"
	"testing"
)

type testCases struct {
	desc     string
	argument []string
	expect   any
}

func TestBulkStringParser(t *testing.T) {
	tests := []testCases{
		{
			desc:     "Normal array of strings",
			argument: []string{"This is an array"},
			expect:   "$16\r\nThis is an array\r\n",
		},
		{
			desc:     "Normal array of strings",
			argument: []string{""},
			expect:   "$0\r\n\r\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := resp.BulkString(tc.argument)
			if got != tc.expect {
				t.Errorf("got %v, want %v", got, tc.expect)
			}
		})
	}
}
