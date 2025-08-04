package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 1600
	screenHeight    = 900
	circleCount     = 200
	speed           = 0.3
	connectDistance = 80
)

type Vec2 struct {
	X, Y float32
}

func (v Vec2) len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vec2) normal() Vec2 {
	return Vec2{X: v.X / v.len(), Y: v.Y / v.len()}
}

func (v Vec2) dist(other Vec2) float32 {
	deltaX := other.X - v.X
	deltaY := other.Y - v.Y
	len := math.Sqrt(float64(deltaX*deltaX) + float64(deltaY*deltaY))
	return float32(math.Abs(len))
}

func addVec(a, b Vec2) Vec2 {
	return Vec2{X: a.X + b.X, Y: a.Y + b.Y}
}

type Dot struct {
	position  Vec2
	direction Vec2
}

type Game struct {
	circles []Dot
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	for i := range g.circles {
		// control speed
		velocity := Vec2{
			X: g.circles[i].direction.X * speed,
			Y: g.circles[i].direction.Y * speed,
		}
		g.circles[i].position = addVec(g.circles[i].position, velocity)

		// control direction
		if g.circles[i].position.X > screenWidth || g.circles[i].position.X < 0 {
			g.circles[i].direction.X = -g.circles[i].direction.X
		}
		if g.circles[i].position.Y > screenHeight || g.circles[i].position.Y < 0 {
			g.circles[i].direction.Y = -g.circles[i].direction.Y
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, c := range g.circles {
		vector.DrawFilledCircle(screen, c.position.X, c.position.Y, 1.5, color.White, true)
	}

	for i := 0; i < len(g.circles); i++ {
		for j := i + 1; j < len(g.circles); j++ {
			if g.circles[i].position.dist(g.circles[j].position) < connectDistance {
				strokeWidth := (connectDistance / g.circles[i].position.dist(g.circles[j].position)) * 0.2
				if strokeWidth > 1 {
					strokeWidth = 1
				}
				vector.StrokeLine(screen, g.circles[i].position.X, g.circles[i].position.Y, g.circles[j].position.X, g.circles[j].position.Y, strokeWidth, color.White, true)
			}

			// consider cursor position to connect dots
			cx, cy := ebiten.CursorPosition()
			cursorPos := Vec2{float32(cx), float32(cy)}
			if g.circles[i].position.dist(cursorPos) < connectDistance + 30 {
				strokeWidth := (connectDistance / g.circles[i].position.dist(cursorPos)) * 0.2
				if strokeWidth > 1 {
					strokeWidth = 1
				}
				vector.StrokeLine(screen, g.circles[i].position.X, g.circles[i].position.Y, cursorPos.X, cursorPos.Y, strokeWidth, color.White, true)
			}
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %v, TPS: %v", int(ebiten.ActualFPS()), int(ebiten.ActualTPS())), 10, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dots")
	g := &Game{
		circles: make([]Dot, circleCount),
	}
	for range circleCount {
		g.circles = append(g.circles, Dot{
			position: Vec2{
				X: float32(rand.Intn(screenWidth)), Y: float32(rand.Intn(screenHeight)),
			},
			direction: randDir(),
		})
	}
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

// returns a normalized random direction
func randDir() Vec2 {
	dir := Vec2{
		X: float32(rand.Float64()*2 - 1),
		Y: float32(rand.Float64()*2 - 1),
	}
	return dir.normal()
}
