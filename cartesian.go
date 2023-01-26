package match

import "sync"

func cartesian(params ...[]string) [][]string {
	cp := [][]string{}
	generator := iterCartesian(params...)
	for product := range generator {
		cp = append(cp, product)
	}
	return cp
}

func iterCartesian(params ...[]string) chan []string {
	c := make(chan []string)
	var wg sync.WaitGroup
	wg.Add(1)
	iterateCartesian(&wg, c, []string{}, params...)
	go func() { wg.Wait(); close(c) }()
	return c
}

func iterateCartesian(wg *sync.WaitGroup, channel chan []string, result []string, params ...[]string) {
	defer wg.Done()
	if len(params) == 0 {
		channel <- result
		return
	}
	p, params := params[0], params[1:]
	for i := 0; i < len(p); i++ {
		wg.Add(1)
		resultCopy := append([]string{}, result...)
		go iterateCartesian(wg, channel, append(resultCopy, p[i]), params...)
	}
}
