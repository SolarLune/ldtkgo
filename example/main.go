package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/ebitenrenderer"
)

type Game struct {
	LDTKProject    *ldtkgo.Project
	EbitenRenderer *renderer.EbitenRenderer
	CurrentLevel   int
	ActiveLayers   []bool
}

func NewGame() *Game {

	g := &Game{
		ActiveLayers: []bool{true, true, true, true},
	}

	var err error

	// Load LDtk Project
	g.LDTKProject, err = ldtkgo.LoadFile("example.ldtk")

	if err != nil {
		panic(err)
	}

	// Choose a level...
	level := g.LDTKProject.Levels[0]

	// Create a new renderer...
	g.EbitenRenderer = renderer.NewEbitenRenderer(renderer.NewDiskLoader("")) // DiskLoader loads from disk using ebitenutil.NewImageFromFile().

	// ... And render the tiles out to *ebiten.Images. We'll draw them later in the Draw() loop.
	g.EbitenRenderer.Render(level)

	fmt.Println("Press the 1 - 4 keys to toggle the tileset layers. Press the Left or Right arrow keys to switch Levels.")

	return g

}

func (g *Game) Update() error {

	levelChanged := false

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.CurrentLevel++
		levelChanged = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.CurrentLevel--
		levelChanged = true
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.ActiveLayers[0] = !g.ActiveLayers[0]
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.ActiveLayers[1] = !g.ActiveLayers[1]
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.ActiveLayers[2] = !g.ActiveLayers[2]
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.ActiveLayers[3] = !g.ActiveLayers[3]
	}

	if g.CurrentLevel >= len(g.LDTKProject.Levels) {
		g.CurrentLevel = 0
	}

	if g.CurrentLevel < 0 {
		g.CurrentLevel = len(g.LDTKProject.Levels) - 1
	}

	if levelChanged {
		g.EbitenRenderer.Render(g.LDTKProject.Levels[g.CurrentLevel])
	}

	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(g.LDTKProject.Levels[0].BGColor) // We want to use the BG Color when possible

	for i, layer := range g.EbitenRenderer.RenderedLayers {
		if g.ActiveLayers[i] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}

}

func (g *Game) Layout(w, h int) (int, int) { return 320, 240 }

func main() {

	g := NewGame()

	ebiten.SetWindowResizable(true)

	ebiten.SetWindowTitle("LDtk-Go Example")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
