package bytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestIs(t *testing.T) {
	testFunc := func(input []byte, key []rune) {
		for i := range input {
			correct := true
			ok, _ := Is(key[i])(input[i:])
			if ok != (key[i] == rune(input[i])) {
				correct = false
				t.Fail()
			}

			t.Log(
				fmt.Sprintf(
					"is(%v)(%v) -> %v\tcorrect=%v",
					string(input[i]), string(key[i]), ok, correct))
		}
	}

	testFunc([]byte("fooBaz"), []rune("foobar"))
	testFunc([]byte("andy"), []rune("foobar"))
	testFunc(nil, []rune("foobar"))
	testFunc([]byte(""), []rune("foobar"))
	testFunc(nil, nil)
	testFunc([]byte(""), []rune(""))
}

func TestNot(t *testing.T) {
	testFunc := func(input []byte, key []rune) {
		for i := range input {
			correct := true
			ok, _ := Not(key[i])(input[i:])
			if ok != (key[i] != rune(input[i])) {
				correct = false
				t.Fail()
			}
			t.Log(
				fmt.Sprintf(
					"not(%v)(%v) ->  %v\tcorrect=%v",
					string(input[i]), string(key[i]), ok, correct))
		}
	}

	testFunc([]byte("fooBaz"), []rune("foobar"))
	testFunc([]byte("andy"), []rune("foobar"))
	testFunc(nil, []rune("foobar"))        //FIXME
	testFunc([]byte(""), []rune("foobar")) //FIXME
}
func TestRepeat_Is(t *testing.T) {
	for i := 0; i < 10; i++ {
		input := bytes.Repeat([]byte("*"), 10)
		input[i] = byte('?')

		ok, n := Repeats(Is('*'))(input)
		if n != i {
			t.Fail()
		}
		t.Log(
			fmt.Sprintf(
				"repeats(is('*'))(%v) -> %v\tcorrect=%v",
				string(input), ok, i == n))
	}
}

func TestRepeat_Not(t *testing.T) {
	for i := 0; i < 10; i++ {
		input := bytes.Join(
			[][]byte{bytes.Repeat([]byte("*"), i),
				bytes.Repeat([]byte("?"), 10-i)}, []byte(""))

		ok, n := Repeats(Not('?'))(input)
		if n != i {
			t.Fail()
		}

		t.Log(
			fmt.Sprintf(
				"repeats(not('?'))(%v) -> %v\tcorrect=%v",
				string(input), ok, i == n))
	}
}

func TestThen(t *testing.T) {
	input := []byte("****Hi!")

	// Test 1
	ok1, n1 := Repeats(Is('*')).Then(Is('H')).Then(Is('i')).Then(Is('!'))(input)
	correct := ok1 && (n1 == len(input))
	if !correct {
		t.Fail()
	}
	t.Log(fmt.Sprintf("Test 1\t(%v,%q)\tcorrect=%v", ok1, string(input[n1:]), correct))

	// Test 2
	ok2, n2 := Repeats(Is('*')).Then(Is('H')).Then(Is('i'))(input)
	correct = ok2 && (n2 == (len(input) - 1))
	if !correct {
		t.Fail()
	}
	t.Log(fmt.Sprintf("Test 2\t(%v,%q)\tcorrect=%v", ok2, string(input[n2:]), correct))

	// Test 3
	ok3, n3 := Repeats(Is('*')).Then(Is('H')).Then(Is('i')).Then(Is('?'))(input)
	correct = !ok3 && (n3 == (len(input) - 1))
	if !correct {
		t.Fail()
	}
	t.Log(fmt.Sprintf("Test 3\t(%v,%q)\tcorrect=%v", ok3, string(input[n3:]), correct))

	// Test 4
	ok4, n4 := Is('*').Then(Is('*')).Then(Is('*')).Then(Not('?'))(input)
	correct = ok4 && (n4 == (len(input) - 3))
	if !correct {
		t.Fail()
	}
	t.Log(fmt.Sprintf("Test 4\t(%v,%q)\tcorrect=%v", ok4, string(input[n4:]), correct))
}
