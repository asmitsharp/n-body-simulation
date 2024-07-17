package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1000
	screenHeight = 800
	G            = 6.67430e-11 // gravitational constant
	timeStep     = 1.0 / 60    // simulation time step
	scaleFactor  = 1e-9        // scale factor to make the simulation visible
	orbitScale   = 1e-9        // scale down the orbit sizes to fit on screen
	speedScale   = 300000
)

type Vector2D struct {
	X, Y float64
}

type Body struct {
	Position Vector2D
	Velocity Vector2D
	Mass     float64
	Radius   float64
	Color    color.Color
}

type Simulation struct {
	Bodies []Body
}

func NewSimulation() *Simulation {
	return &Simulation{
		Bodies: make([]Body, 0),
	}
}

func (s *Simulation) AddBody(b Body) {
	s.Bodies = append(s.Bodies, b)
}

func (s *Simulation) Update() {
	for i := range s.Bodies {
		force := Vector2D{}
		for j := range s.Bodies {
			if i != j {
				force = addVectors(force, calculateGravitationalForce(&s.Bodies[i], &s.Bodies[j]))
			}
		}
		acceleration := scaleVector(force, 1/s.Bodies[i].Mass)
		s.Bodies[i].Velocity = addVectors(s.Bodies[i].Velocity, scaleVector(acceleration, timeStep))
		s.Bodies[i].Position = addVectors(s.Bodies[i].Position, scaleVector(s.Bodies[i].Velocity, timeStep))

		// Keep bodies within the screen
		s.Bodies[i].Position.X = math.Mod(s.Bodies[i].Position.X+screenWidth, screenWidth)
		s.Bodies[i].Position.Y = math.Mod(s.Bodies[i].Position.Y+screenHeight, screenHeight)
	}
}

func calculateGravitationalForce(b1, b2 *Body) Vector2D {
	dx := b2.Position.X - b1.Position.X
	dy := b2.Position.Y - b1.Position.Y
	distSq := dx*dx + dy*dy
	dist := math.Sqrt(distSq)

	// Softening factor to prevent extreme forces at small distances
	softening := 1e7
	force := G * b1.Mass * b2.Mass / (distSq + softening*softening)

	return Vector2D{
		X: force * dx / dist * scaleFactor,
		Y: force * dy / dist * scaleFactor,
	}
}

func addVectors(v1, v2 Vector2D) Vector2D {
	return Vector2D{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func scaleVector(v Vector2D, scalar float64) Vector2D {
	return Vector2D{X: v.X * scalar, Y: v.Y * scalar}
}

type Game struct {
	sim *Simulation
}

func (g *Game) Update() error {
	g.sim.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, body := range g.sim.Bodies {
		ebitenutil.DrawCircle(screen, body.Position.X, body.Position.Y, body.Radius, body.Color)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	sim := NewSimulation()

	sun := Body{
		Position: Vector2D{X: screenWidth / 2, Y: screenHeight / 2},
		Velocity: Vector2D{X: 0, Y: 0},
		Mass:     1.989e30, // Mass of the Sun in kg
		Radius:   20,
		Color:    color.RGBA{255, 255, 0, 255},
	}
	sim.AddBody(sun)

	// Venus
	venusOrbitRadius := 108.2e9 * orbitScale         // 108.2 million km
	venusSpeed := 35.02e3 * speedScale * scaleFactor // 35.02 km/s
	venus := Body{
		Position: Vector2D{X: screenWidth/2 + venusOrbitRadius, Y: screenHeight / 2},
		Velocity: Vector2D{X: 0, Y: -venusSpeed},
		Mass:     4.867e24, // Mass of Venus in kg
		Radius:   4,
		Color:    color.RGBA{255, 198, 73, 255}, // Light orange
	}
	sim.AddBody(venus)

	// Earth
	earthOrbitRadius := 149.6e9 * orbitScale         // 149.6 million km
	earthSpeed := 29.78e3 * speedScale * scaleFactor // 29.78 km/s
	earth := Body{
		Position: Vector2D{X: screenWidth/2 + earthOrbitRadius, Y: screenHeight / 2},
		Velocity: Vector2D{X: 0, Y: -earthSpeed},
		Mass:     5.972e24, // Mass of the Earth in kg
		Radius:   5,
		Color:    color.RGBA{0, 0, 255, 255},
	}
	sim.AddBody(earth)

	// Earth's Moon
	moonOrbitRadius := 384400e3 * orbitScale                                              // 384,400 km
	moonSpeed := (1.022e3 + earthSpeed/scaleFactor/speedScale) * speedScale * scaleFactor // 1.022 km/s + Earth's speed
	moon := Body{
		Position: Vector2D{X: earth.Position.X + moonOrbitRadius, Y: earth.Position.Y},
		Velocity: Vector2D{X: 0, Y: -moonSpeed},
		Mass:     7.34767309e22, // Mass of the Moon in kg
		Radius:   2,
		Color:    color.RGBA{200, 200, 200, 255}, // Light grey
	}
	sim.AddBody(moon)

	// Mars
	marsOrbitRadius := 227.9e9 * orbitScale          // 227.9 million km
	marsSpeed := 24.077e3 * speedScale * scaleFactor // 24.077 km/s
	mars := Body{
		Position: Vector2D{X: screenWidth/2 + marsOrbitRadius, Y: screenHeight / 2},
		Velocity: Vector2D{X: 0, Y: -marsSpeed},
		Mass:     6.39e23, // Mass of Mars in kg
		Radius:   4,
		Color:    color.RGBA{255, 0, 0, 255},
	}
	sim.AddBody(mars)

	// Jupiter
	jupiterOrbitRadius := 778.5e9 * orbitScale         // 778.5 million km
	jupiterSpeed := 13.07e3 * speedScale * scaleFactor // 13.07 km/s
	jupiter := Body{
		Position: Vector2D{X: screenWidth/2 + jupiterOrbitRadius, Y: screenHeight / 2},
		Velocity: Vector2D{X: 0, Y: -jupiterSpeed},
		Mass:     1.898e27, // Mass of Jupiter in kg
		Radius:   15,
		Color:    color.RGBA{255, 140, 0, 255}, // Dark orange
	}
	sim.AddBody(jupiter)

	game := &Game{
		sim: sim,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Solar System Simulation")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
