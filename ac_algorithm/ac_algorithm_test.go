package ac_algorithm

import "testing"


func TestAutomatonSearch(t *testing.T) {
	automaton := NewAutomaton([]string{"abc", "add", "aaabc", "sdfdd"})

	states, failStatues := automaton.printFailTransitions()
	t.Log(states)
	t.Log(failStatues)

	matches, counts := automaton.Search("safujgoaaabcsdabcasjfaddddabc")
	t.Log(matches, counts)

}
