package email_test

import (
	"testing"
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/lib/email"
)

func Test1(t *testing.T) {
	log.DebugOn()
	tests := []struct{addr string; valid bool}{
		{"1", false},
		{"2", false},
		{"a", false},
		{"ab", false},
		{"A", false},
		{"AB", false},
		{"A-B", false},
		{"A--B", false},
		{"A-", false},
		//invalid
		{"-A", false},
		{"A.B", false},
		{"A..B", false},
		{".A", false},
		{"A.", false},
		{"@A", false},
		{"A@", false},
		{"-", false},
		{".", false},
		{"@", false},
		{"~A", false},
		{"a@b", false},
		{"a@b.c", true},
		{"one.two@three.four", true},
		{"one-two.three-four@five-size.seven-eight.nine-ten", true},
	}

	for _,test := range tests {
		//log.Debugf("Test[%d]: %s", i, test.addr)
		err := email.CheckAddress(test.addr);
		result := "valid"
		if err != nil {
			result="invalid"
		}
		log.Debugf("%10s %s", result, test.addr)
		if test.valid && err != nil {
			t.Fatalf("%s is valid but not accepted", test.addr)
		}
		if !test.valid && err == nil {
			t.Fatalf("%s is invalid but accepted", test.addr)
		}
	}
}