package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 1600
	screenHeight    = 900
	pointsCount     = 200
	speed           = 0.2
	connectDistance = 100
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
	dx := other.X - v.X
	dy := other.Y - v.Y
	len := math.Sqrt(float64(dx*dx) + float64(dy*dy))
	return float32(math.Abs(len))
}

func addVec(a, b Vec2) Vec2 {
	return Vec2{X: a.X + b.X, Y: a.Y + b.Y}
}

type Point struct {
	position  Vec2
	direction Vec2
}

type Game struct {
	points []Point
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	for i := range g.points {
		// control speed
		velocity := Vec2{
			X: g.points[i].direction.X * speed,
			Y: g.points[i].direction.Y * speed,
		}
		g.points[i].position = addVec(g.points[i].position, velocity)

		// control direction
		if g.points[i].position.X > screenWidth || g.points[i].position.X < 0 {
			g.points[i].direction.X = -g.points[i].direction.X
		}
		if g.points[i].position.Y > screenHeight || g.points[i].position.Y < 0 {
			g.points[i].direction.Y = -g.points[i].direction.Y
		}
	}

	return nil
}

// returns a normalized random direction
func randDir() Vec2 {
	dir := Vec2{
		X: float32(rand.Float64()*2 - 1),
		Y: float32(rand.Float64()*2 - 1),
	}
	return dir.normal()
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, c := range g.points {
		vector.DrawFilledCircle(screen, c.position.X, c.position.Y, 1.5, color.White, true)
	}

	sort.Slice(g.points, func(i, j int) bool {
		di := g.points[i].position.dist(Vec2{0, 0})
		dj := g.points[j].position.dist(Vec2{0, 0})
		return di < dj
	})

	for i := 0; i < len(g.points); i++ {
		countingLimit := i + 50
		for j := i + 1; j < countingLimit; j++ {
			if j >= len(g.points) {
				break
			}
			if g.points[i].position.dist(g.points[j].position) < connectDistance {
				strokeWidth := float32(math.Abs(-1 + float64((g.points[i].position.dist(g.points[j].position) / 100))))
				if strokeWidth > 0.5 {
					strokeWidth = 0.5
				}
				vector.StrokeLine(screen, g.points[i].position.X, g.points[i].position.Y, g.points[j].position.X, g.points[j].position.Y, strokeWidth, color.White, true)
			}

			// consider cursor position to connect dots
			cx, cy := ebiten.CursorPosition()
			cursorPos := Vec2{float32(cx), float32(cy)}
			if g.points[i].position.dist(cursorPos) < connectDistance+30 {
				strokeWidth := float32(math.Abs(-1 + float64((g.points[i].position.dist(cursorPos) / 100))))
				if strokeWidth > 0.5 {
					strokeWidth = 0.5
				}
				vector.StrokeLine(screen, g.points[i].position.X, g.points[i].position.Y, cursorPos.X, cursorPos.Y, strokeWidth, color.White, true)
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
	ebiten.SetVsyncEnabled(false)
	g := &Game{
		points: make([]Point, 0),
	}
	for range pointsCount {
		g.points = append(g.points, Point{
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
