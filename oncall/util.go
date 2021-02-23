package oncall

func getSetVennDiagram(left, right []string) (leftOnly []string, rightOnly []string, intersection []string, sum []string) {
	// Create two "sets", one of current users and one of target users
	setLeft := map[string]bool{}
	setRight := map[string]bool{}
	setSum := map[string]bool{}
	for _, u := range right {
		setRight[u] = true
		setSum[u] = true
	}
	for _, u := range left {
		setLeft[u] = true
		setSum[u] = true
	}

	for u, _ := range setSum {
		sum = append(sum, u)
	}

	for u, _ := range setLeft {
		_, targetUser := setRight[u]
		if !targetUser {
			leftOnly = append(leftOnly, u)
		} else {
			intersection = append(intersection, u)
		}
	}

	for u, _ := range setRight {
		_, targetUser := setLeft[u]
		if !targetUser {
			rightOnly = append(rightOnly, u)
		}
	}

	return
}
