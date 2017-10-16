package graph

// lockNodes and unlockNodes is always same order to avoid dead-lock.

func lockNodes(n1, n2 Node) {
	s1, s2 := sortNodes(n1, n2)
	s1.lock()
	s2.lock()
}

func unlockNodes(n1, n2 Node) {
	s1, s2 := sortNodes(n1, n2)
	s2.unlock()
	s1.unlock()
}

func sortNodes(n1, n2 Node) (Node, Node) {
	if n1.lockNum() < n2.lockNum() {
		return n1, n2
	}
	return n2, n1
}
