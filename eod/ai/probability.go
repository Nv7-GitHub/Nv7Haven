package ai

import (
	"math/rand"
	"sync"
)

type Probability struct {
	lock *sync.RWMutex
	Data map[int]int
	Sum  int
}

func NewProbability() *Probability {
	return &Probability{
		lock: &sync.RWMutex{},
		Data: make(map[int]int),
	}
}

func (p *Probability) Add(id int, nolock bool) {
	if !nolock {
		p.lock.Lock()
		defer p.lock.Unlock()
	}
	_, exists := p.Data[id]
	if exists {
		p.Data[id]++
	} else {
		p.Data[id] = 1
	}
	p.Sum++
}

func (p *Probability) Predict() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	num := rand.Intn(p.Sum)

	// Weighted random
	for k, v := range p.Data {
		if num < v {
			return k
		}
		num -= v
	}
	return 0
}
