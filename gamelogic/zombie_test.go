package gamelogic

import (
	"testing"
)

func TestZombie_move(t *testing.T) {
	z := newZombie()
	initX, initY := z.x, z.y

	t.Run("move, not at any edge of the board", func(t *testing.T) {
		z.move()
		if z.x-initX < 0 || z.x-initX > 1 {
			t.Fatalf("move() should have added 0 or 1 to z.x")
		}
		if z.y-initY < -1 || z.y-initY > 1 {
			t.Fatalf("move() should have added -1, 0 or 1 to y")
		}
	})

	t.Run("move, at the top edge of the board", func(t *testing.T) {
		// init state
		for z.y != 30 {
			z.move()
		}
		initX, initY = z.x, z.y

		z.move()
		if z.x-initX < 0 || z.x-initX > 1 {
			t.Fatalf("move() should have added 0 or 1 to z.x")
		}
		if z.y-initY != -1 {
			t.Fatalf("move() at the top edge should have decremented y by 1")
		}
	})

	t.Run("move, at the bottom edge of the board", func(t *testing.T) {
		// init state
		for z.y != 0 {
			z.move()
		}
		initX, initY = z.x, z.y

		z.move()
		if z.x-initX < 0 || z.x-initX > 1 {
			t.Fatalf("move() should have added 0 or 1 to z.x")
		}
		if z.y-initY != 1 {
			t.Fatalf("move() at the bottom edge should have incremented y by 1")
		}
	})
}
