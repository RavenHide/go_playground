package ac_algorithm

import "container/list"

const ROOT_STATE uint32 = 0

type Automaton struct {
	rootNode    *node
	maxStateNum uint32
}

func NewAutomaton(words []string) *Automaton {
	inst := &Automaton{
		rootNode: newNode(ROOT_STATE),
	}
	inst.rootNode.failNode = inst.rootNode

	for _, word := range words {
		inst.addWord(word)
	}
	inst.initFailTransitions()
	return inst
}

func (m *Automaton) Search(text string) ([]string, []uint32) {
	matchTextCount := make(map[string]uint32, 0)

	node := m.rootNode
	for _, char := range text {

		nextNode := node.gotoNode(char)
		for nextNode == nil && node.state != ROOT_STATE {
			node = node.failNode
			nextNode = node.gotoNode(char)
		}

		if nextNode == nil{
			node = m.rootNode
		} else {
			node = nextNode
			for _, word := range node.outputWords() {
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
	node := m.rootNode
	for _, char := range word {
		nextNode := node.gotoNode(char)
		if nextNode == nil {
			m.maxStateNum += 1
			nextNode = newNode(m.maxStateNum)
			node.gotoNodeMap[char] = nextNode
		}
		node = nextNode
	}
	node.addOutputWord(word)
}

func (m *Automaton) printFailTransitions() ([]uint32, []uint32) {
	states := make([]uint32, 0, m.maxStateNum)
	failStates := make([]uint32, 0, m.maxStateNum)

	queue := list.New()
	queue.PushBack(m.rootNode)

	item := queue.Front()

	for item != nil {
		curNode := item.Value.(*node)
		states = append(states, curNode.state)
		failStates = append(failStates, curNode.failNode.state)
		for _, node := range curNode.gotoNodeMap {
			queue.PushBack(node)
		}

		item = item.Next()
	}

	return states, failStates
}

func (m *Automaton) initFailTransitions() {

	queue := list.New()
	for _, node := range m.rootNode.gotoNodeMap {
		node.failNode = m.rootNode
		queue.PushBack(node)
	}
	curNode := queue.Front()

	for curNode != nil {
		mNode := curNode.Value.(*node)
		for char, gotoNode := range mNode.gotoNodeMap {
			queue.PushBack(gotoNode)
			// 递归地找出gotoNode的failedNode
			for mNode.failNode.state != ROOT_STATE && mNode.failNode.gotoNode(char) == nil {
				mNode = mNode.failNode
			}

			if mNode.failNode.state == ROOT_STATE && mNode.failNode.gotoNode(char) == nil {
				gotoNode.failNode = m.rootNode
			} else {
				gotoNode.failNode = mNode.failNode.gotoNode(char)
			}
			// 更新output
			gotoNode.addOutputWord(gotoNode.failNode.outputWords()...)
		}

		curNode = curNode.Next()
	}
}

type node struct {
	state       uint32
	failNode    *node
	output      []string
	gotoNodeMap map[int32]*node
}

func newNode(state uint32) *node {
	return &node{
		state:       state,
		failNode:    nil,
		output:      nil,
		gotoNodeMap: make(map[int32]*node, 0),
	}
}

func (n *node) gotoNode(char int32) *node {
	return n.gotoNodeMap[char]
}

func (n *node) addOutputWord(words ...string) {
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

func (n *node) outputWords()[]string{
	return n.output
}
