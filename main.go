package main

import (
	"fmt"
	"time"
)

var elves = 0
var reindeer = 0
var mutex = NewSemaphore(1)
var santaSem = NewSemaphore(0)
var reindeerSem = NewSemaphore(0)
var elfTex = NewSemaphore(1)

func main() {
	for i := 1; i <= 9; i++ {
		go reindeerArrives()
	}

	for i := 1; i <= 100; i++ {
		go elfArrives()
	}

	for {
		santaSem.Wait()
		mutex.Wait()

		if reindeer == 9 {
			prepareSleigh()
			reindeerSem.Signal()
			fmt.Println("É Natal, entregando presentes!")
			return // Encerre o programa
		} else if elves == 3 {
			helpElves()
		}

		mutex.Signal()
	}
}

func reindeerArrives() {
	time.Sleep(time.Duration(1) * time.Second) // Reindeer's arrival
	fmt.Printf("Rena: Volta das férias\n")

	mutex.Wait()
	reindeer++
	if reindeer == 9 {
		santaSem.Signal()
	}
	mutex.Signal()

	reindeerSem.Wait()
	getHitched()
}

func elfArrives() {
	time.Sleep(time.Duration(1) * time.Second) // Elf arrives
	fmt.Println("Elfo: Pede Ajuda")

	elfTex.Wait()
	mutex.Wait()
	elves++
	if elves == 3 {
		santaSem.Signal()
	} else {
		elfTex.Signal()
	}
	mutex.Signal()

	getHelp()

	mutex.Wait()
	elves--
	if elves == 0 {
		elfTex.Signal()
	}
	mutex.Signal()
}

func prepareSleigh() {
	fmt.Println("Santa: Preparando o trenó")
}

func getHitched() {
	fmt.Println("Santa: Preparando as renas")
	reindeerSem.Signal()
}

func helpElves() {
	fmt.Println("Santa: Ajudando os elfos")
	time.Sleep(1 * time.Second)
}

func getHelp() {
	fmt.Println("Elfos: Recebendo ajuda")
	time.Sleep(1 * time.Second) // Simulate the elf getting help
}

// The rest of the code remains the same.



type Semaphore struct {
	v    int
	fila chan struct{}
	sc   chan struct{}
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v:    init,
		fila: make(chan struct{}),
		sc:   make(chan struct{}, 1),
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sc <- struct{}{}
	s.v--
	if s.v < 0 {
		<-s.sc
		s.fila <- struct{}{}
	} else {
		<-s.sc
	}
}

func (s *Semaphore) Signal() {
	s.sc <- struct{}{}
	s.v++
	if s.v <= 0 {
		<-s.fila
	}
	<-s.sc
}

