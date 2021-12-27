package ac_algorithm

import "container/list"

const ROOT_STATE uint32 = 0

type Automaton struct {
	rootState   *stateNode
	maxStateNum uint32
}

func NewAutomaton(words []string) *Automaton {
	inst := &Automaton{
		rootState: newState(ROOT_STATE),
	}
	inst.rootState.failState = inst.rootState

	for _, word := range words {
		inst.addWord(word)
	}
	inst.initFailTransitions()
	return inst
}

func (m *Automaton) Search(text string) ([]string, []uint32) {
	matchTextCount := make(map[string]uint32, 0)

	curState := m.rootState
	for _, char := range text {

		nextState := curState.gotoState(char)
		for nextState == nil && curState.state != ROOT_STATE {
			curState = curState.failState
			nextState = curState.gotoState(char)
		}

		if nextState == nil{
			curState = m.rootState
		} else {
			curState = nextState
			for _, word := range curState.outputWords() {
				if _, ok := matchTextCount[word]; ok {
					matchTextCount[word] += 1
				} else {
					matchTextCount[word] = 1
				}
			}
		}
	}

	matches := make([]string, 0, len(matchTextCount))
	matchCounts := make([]uint32, 0, len(matchTextCount))

	for word, count := range matchTextCount {
		matches = append(matches, word)
		matchCounts = append(matchCounts, count)
	}

	return matches, matchCounts

}

func (m *Automaton) addWord(word string) {
	curState := m.rootState
	for _, char := range word {
		nextState := curState.gotoState(char)
		if nextState == nil {
			m.maxStateNum += 1
			nextState = newState(m.maxStateNum)
			curState.gotoStateMap[char] = nextState
		}
		curState = nextState
	}
	curState.addOutputWord(word)
}

func (m *Automaton) printFailTransitions() ([]uint32, []uint32) {
	states := make([]uint32, 0, m.maxStateNum)
	failStates := make([]uint32, 0, m.maxStateNum)

	queue := list.New()
	queue.PushBack(m.rootState)

	stateItem := queue.Front()

	for stateItem != nil {
		curState := stateItem.Value.(*stateNode)
		states = append(states, curState.state)
		failStates = append(failStates, curState.failState.state)
		for _, node := range curState.gotoStateMap {
			queue.PushBack(node)
		}

		stateItem = stateItem.Next()
	}

	return states, failStates
}

func (m *Automaton) initFailTransitions() {

	queue := list.New()
	for _, node := range m.rootState.gotoStateMap {
		node.failState = m.rootState
		queue.PushBack(node)
	}
	stateItem := queue.Front()

	for stateItem != nil {
		curState := stateItem.Value.(*stateNode)
		for char, gotoState := range curState.gotoStateMap {
			queue.PushBack(gotoState)
			// 递归地找出gotoNode的failedNode
			for curState.failState.state != ROOT_STATE && curState.failState.gotoState(char) == nil {
				curState = curState.failState
			}

			if curState.failState.state == ROOT_STATE && curState.failState.gotoState(char) == nil {
				gotoState.failState = m.rootState
			} else {
				gotoState.failState = curState.failState.gotoState(char)
			}
			// 更新output
			gotoState.addOutputWord(gotoState.failState.outputWords()...)
		}

		stateItem = stateItem.Next()
	}
}

type stateNode struct {
	state        uint32
	failState    *stateNode
	output       []string
	gotoStateMap map[int32]*stateNode
}

func newState(state uint32) *stateNode {
	return &stateNode{
		state:        state,
		failState:    nil,
		output:       nil,
		gotoStateMap: make(map[int32]*stateNode, 0),
	}
}

func (n *stateNode) gotoState(char int32) *stateNode {
	return n.gotoStateMap[char]
}

func (n *stateNode) addOutputWord(words ...string) {
	if len(words) <= 0 {
		return
	}

	for _, word := range words {
		var exists bool
		for _, ouputWord := range n.outputWords() {
			if ouputWord != word {
				continue
			}
			exists = true
			break
		}
		if !exists {
			n.output = append(n.output, word)
		}
	}
	return
}

func (n *stateNode) outputWords()[]string{
	return n.output
}
