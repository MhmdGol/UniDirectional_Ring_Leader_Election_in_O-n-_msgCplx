package main

import (
	"fmt"
	"math/rand"
	"time"
)

// the main part, my election O(n) messaging algorithms
// *intentionally written in long way
func node(take <-chan int, give chan<- int, UID int, leaderUID chan int) {
	// ignnore: fmt.Println(take, give, UID, leaderUID)

	// roles: Leader, Relay, TempLeader
	role := ""

	// start state
	// first message sent by each process
	give <- UID
	recievedUID := <-take

	// CMP state
	if recievedUID < UID {
		role = "TempLeader"
	} else if recievedUID > UID {
		role = "Relay"
		give <- recievedUID
	}

	// ignore: fmt.Println(UID, role)

	// Temp Leader state
	// leader is a temp leader who never recieves a bigger UID
	for role == "TempLeader" {
		recievedUID = <-take
		if recievedUID == UID {
			// Leader state
			role = "Leader"
			// using leaderUID chan causes termination
			leaderUID <- UID
			break
		} else if recievedUID > UID {
			role = "Relay"
			give <- recievedUID
			break
		}
	}

	//Relay state
	for role == "Relay" {
		msg := <-take
		give <- msg
	}
}

func main() {
	//a channel to block main goroutin till leader election
	leaderUID := make(chan int, 1)

	// ----------< making a list of UIDs >----------
	// number of nodes in ring
	n := 100
	// buffered channel
	listOfUIDs := make([]int, n)
	UIDmap := make(map[int]bool)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < n; i++ {
		exist := true
		var generatedUID int
		for exist {
			generatedUID = r1.Intn(10 * n)
			exist = UIDmap[generatedUID]
		}
		UIDmap[generatedUID] = true
		listOfUIDs[i] = generatedUID
	}
	// --------------\/\/\/\/\/---------------

	// ----------< configuring the ring >----------
	// endch links node n to 1
	// buffered channel
	endch := make(chan int, n)
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			go node(endch, ch, listOfUIDs[i], leaderUID)
		} else if i == n-1 {
			go node(ch, endch, listOfUIDs[i], leaderUID)
		} else {
			newch := make(chan int, n)
			go node(ch, newch, listOfUIDs[i], leaderUID)
			ch = newch
		}
	}
	// --------------\/\/\/\/\/---------------

	// printing phase
	msg := <-leaderUID
	fmt.Println(listOfUIDs, "MAX:", max(listOfUIDs))
	fmt.Println(msg)
}

func max(T []int) int {
	m := -1
	for _, v := range T {
		if v > m {
			m = v
		}
	}
	return m
}
