// Package gamelogic contains all game-related logic
package gamelogic

// Game represents one game of winter-is-coming
type Game struct {
	PlayerName            string
	Zombie                Zombie
	ZombieGotHitChan      chan struct{}
	ZombieReachedWallChan chan struct{}
}

// NewGame creates a new game
func NewGame(playerName string) *Game {
	game := &Game{
		PlayerName:            playerName,
		Zombie:                *newZombie(),
		ZombieGotHitChan:      make(chan struct{}),
		ZombieReachedWallChan: make(chan struct{}),
	}

	// the zombie starts walking every two seconds
	game.Zombie.startMoving()

	// when zombie gets hit and dies, notify the tcp handler that created the game
	go func() {
		for {
			<-game.Zombie.diedChan
			close(game.ZombieGotHitChan) // close channel to trigger read from awaiting goroutine
			return
		}
	}()

	cancelSubscriptionChan := make(chan struct{})
	// when zombie reaches wall, stop the zombie and notify the tcp handler that created the game
	cancelSubscription := game.Zombie.SubscribeToMovements(func() {
		if game.Zombie.X() >= 10 {
			defer close(cancelSubscriptionChan)
			defer close(game.ZombieReachedWallChan) // close to trigger a read from awaiting goroutine
			game.Zombie.stopMoving()
		}
	})
	go func() {
		<-cancelSubscriptionChan
		cancelSubscription()
	}()

	return game
}

// Shoot lets the player shoot and try to hit the zombie
func (g *Game) Shoot(x, y int) {
	if g.Zombie.x == x && g.Zombie.y == y {
		g.Zombie.die()
	}
}
