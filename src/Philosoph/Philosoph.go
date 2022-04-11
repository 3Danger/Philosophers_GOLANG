package Philosoph

import (
	"ProducerConsumer/src/Color"
	"fmt"
	"sync"
	"time"
)

type info_common struct {
	eat          time.Duration
	sleep        time.Duration
	time_to_dead time.Duration
	base         time.Time
	msg_mut      sync.Mutex
	toExit       bool
}

type Philosoph struct {
	id             int
	last_eat       time.Time
	info_general   *info_common
	lPhilo, rPlilo *Philosoph
	lFork, rFork   *sync.Mutex
}

func (p *Philosoph) inOrderLockFork() {
	if p.id%2 == 1 {
		p.rFork.Lock()
		p.lFork.Lock()
	} else {
		p.lFork.Lock()
		p.rFork.Lock()
	}
}

func (p *Philosoph) Run() {
	last := p.lPhilo
	for i := 0; i < last.id; i++ {
		p.last_eat = time.Now()
		go p.life()
		p = p.rPlilo
	}
}

func (p *Philosoph) life() {
	for {
		if p.think() || p.eat() || p.sleep() {
			break
		}
	}
}

func (p *Philosoph) isDead() bool {
	if p.info_general.toExit {
		return true
	} else if time.Now().After(p.last_eat.Add(p.info_general.time_to_dead)) {
		p.info_general.toExit = true
		p.message("dead", Color.Red)
		return true
	}
	return false
}

func (p *Philosoph) sleep() bool {
	if p.isDead() {
		return true
	} else {
		p.message("sleep", Color.Purple)
		time.Sleep(p.info_general.sleep)
		return false
	}
}

func (p *Philosoph) think() bool {
	if p.isDead() {
		return true
	} else {
		p.message("think", Color.Blue)
		return false
	}
}

func (p *Philosoph) message(do string, color string) {
	p.info_general.msg_mut.Lock()
	fmt.Print(color)
	fmt.Print(int(time.Now().Sub(p.info_general.base)) / 1000000)
	fmt.Printf("ms. id %2d is %sing%s\n", p.id, do, Color.Reset)
	if do != "dead" {
		p.info_general.msg_mut.Unlock()
	}
}

func (p *Philosoph) eat() bool {
	if p.isDead() {
		return true
	} else {
		p.inOrderLockFork()
		p.last_eat = time.Now()
		p.message("eat", Color.Green)
		time.Sleep(p.info_general.eat)
		p.lFork.Unlock()
		p.rFork.Unlock()
		return false
	}
}

func newPhilosoph(id int, info *info_common, lPhilo, rPhilo *Philosoph, lFork, rFork *sync.Mutex) *Philosoph {
	return &Philosoph{
		id:           id,
		last_eat:     time.Time{},
		info_general: info,
		lPhilo:       lPhilo,
		rPlilo:       rPhilo,
		lFork:        lFork,
		rFork:        rFork,
	}
}

func makeInfo(eat, sleep, time_to_dead int) info_common {
	return info_common{
		eat:          time.Millisecond * time.Duration(eat),
		sleep:        time.Millisecond * time.Duration(sleep),
		time_to_dead: time.Millisecond * time.Duration(time_to_dead),
		base:         time.Now(),
		toExit:       false,
	}
}

func NewPhilosopher(howMuchPhilo, eat, sleep, dead int) *Philosoph {
	info_main := &info_common{}
	result := make([]*Philosoph, howMuchPhilo)
	result[0] = newPhilosoph(1, info_main, nil, nil, nil, &sync.Mutex{})
	for i := 1; i < howMuchPhilo; i++ {
		result[i] = newPhilosoph(i+1, info_main, result[i-1], nil, result[i-1].rFork, &sync.Mutex{})
		result[i-1].rPlilo = result[i]
		result[i-1].lFork = result[i].rFork
	}
	result[0].lPhilo = result[howMuchPhilo-1]
	result[howMuchPhilo-1].rPlilo = result[0]
	*info_main = makeInfo(eat, sleep, dead)
	return result[0]
}

func (p *Philosoph) IsFinished() bool {
	return p.info_general.toExit
}
