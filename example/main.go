package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
)

type Game struct {
	LDTKProject    *ldtkgo.Project
	EbitenRenderer *EbitenRenderer
	BGImage        *ebiten.Image
	CurrentLevel   int
	ActiveLayers   []bool
}

func NewGame() *Game {

	g := &Game{
		ActiveLayers: []bool{true, true, true, true},
	}

	var err error

	// First, we load the LDtk Project. An error would indicate that ldtk-go was unable to find the project file or deserialize the JSON.
	g.LDTKProject, err = ldtkgo.Open("example.ldtk")

	if err != nil {
		panic(err)
	}

	// Seconds, we create a new Renderer.

	// EbitenRenderer.DiskLoader loads images from disk using ebitenutil.NewImageFromFile() and takes an argument of the base path to use when loading.
	// We pass a blank string to NewDiskLoader() because for the example, the assets are in the same directory.
	g.EbitenRenderer = NewEbitenRenderer(NewDiskLoader(""))

	// Then, we render the tiles out to *ebiten.Images contained in the EbitenRenderer. We'll grab them to draw later in the Draw() loop.
	g.RenderLevel()

	fmt.Println("Press the 1 - 4 keys to toggle the tileset layers. Press the Left or Right arrow keys to switch Levels.")

	return g

}

func (g *Game) RenderLevel() {

	if g.CurrentLevel >= len(g.LDTKProject.Levels) {
		g.CurrentLevel = 0
	}

	if g.CurrentLevel < 0 {
		g.CurrentLevel = len(g.LDTKProject.Levels) - 1
	}

	level := g.LDTKProject.Levels[g.CurrentLevel]

	if level.BGImage != nil {
		g.BGImage, _, _ = ebitenutil.NewImageFromFile(level.BGImage.Path)
	} else {
		g.BGImage = nil
	}

	g.EbitenRenderer.Render(level)
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.CurrentLevel++
		g.RenderLevel()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.CurrentLevel--
		g.RenderLevel()
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

	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(g.LDTKProject.Levels[g.CurrentLevel].BGColor) // We want to use the BG Color when possible

	if g.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := g.LDTKProject.Levels[g.CurrentLevel].BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		screen.DrawImage(g.BGImage, opt)
	}

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
