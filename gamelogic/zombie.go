package gamelogic

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nmrshll/winter-is-coming-backend/utils/random"
)

// Zombie represents a zombie on the board on the game
type Zombie struct {
	x    int
	y    int
	name string

	// movementSubscriptions can be used to register callbacks that will be run after every zombie movement
	movementSubscriptions subscriptions
	// a channel to tell it to stop
	stopMovingChan chan struct{}

	// a channel to signal that it got hit and died
	diedChan chan struct{}
}

func newZombie() *Zombie {
	return &Zombie{
		name:     "white-walker",
		y:        rand.Intn(999) % 30, // between 0 and 30
		diedChan: make(chan struct{}),
	}
}

// X returns the x (horizontal) position of the zombie
func (z *Zombie) X() int {
	return z.x
}

// Y returns the y (vertical) position of the zombie
func (z *Zombie) Y() int {
	return z.y
}

// Name returns the name of the zombie
func (z *Zombie) Name() string {
	return z.name
}

// move makes the zombie move 0 or 1 forward on the x axis and -1, 0, or 1 on the y axis, within game board bounds
func (z *Zombie) move() {
	z.x += rand.Intn(999) % 2 // 0 or 1, approximately random

	switch z.y {
	case 0:
		z.y++
	case 30: // upper bound of game board
		z.y--
	default:
		z.y += rand.Intn(999)%3 - 1 // -1, 0, or 1, approximately random
	}

	for _, callback := range z.movementSubscriptions.callbacks {
		callback()
	}
}

// startMoving makes the zombie start moving every 2 seconds
func (z *Zombie) startMoving() (stoppedMovingChan <-chan struct{}) {
	z.stopMovingChan = make(chan struct{})
	stoppedMovingChanBidir := make(chan struct{})

	// move every two seconds until received signal to stop moving
	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				z.move()
			case <-z.stopMovingChan:
				close(stoppedMovingChanBidir) // close channel to trigger a read from awaiting goroutine
				return
			}
		}
	}()
	return stoppedMovingChanBidir
}

func (z *Zombie) stopMoving() {
	close(z.stopMovingChan) // close to trigger read
}

// die kills the zombie and signals it on a channel
func (z *Zombie) die() {
	z.stopMoving()
	close(z.diedChan)
}

type subscriptions struct {
	lock      *sync.RWMutex // the callbacks map might be accessed concurrently, we need to protect against race conditions
	callbacks map[string]func()
}

// SubscribeToMovements lets you register a function that gets run every time the zombie moves
func (z *Zombie) SubscribeToMovements(callback func()) (cancelSubscription func()) {
	// init fields
	{
		if z.movementSubscriptions.lock == nil {
			z.movementSubscriptions.lock = &sync.RWMutex{}
		}
		z.movementSubscriptions.lock.Lock()
		defer z.movementSubscriptions.lock.Unlock()

		if z.movementSubscriptions.callbacks == nil {
			z.movementSubscriptions.callbacks = make(map[string]func())
		}
	}

	// add the passed callback to the list of callbacks
	key := random.String(8)
	z.movementSubscriptions.callbacks[key] = callback

	// return cancel function
	return func() {
		z.movementSubscriptions.lock.Lock()
		defer z.movementSubscriptions.lock.Unlock()
		delete(z.movementSubscriptions.callbacks, key)
	}
}

// CancelAllMovementSubscriptions removes all zombie movement callbacks at once
func (z *Zombie) CancelAllMovementSubscriptions() {
	z.movementSubscriptions.lock.Lock()
	defer z.movementSubscriptions.lock.Unlock()
	for key := range z.movementSubscriptions.callbacks {
		delete(z.movementSubscriptions.callbacks, key)
	}
}
