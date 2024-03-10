package main

import (
	"fmt"
	"image"
	"log"
	"os"

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

	// First, we load the LDtk Project. An error would indicate that ldtk-go was unable to find the project file or deserialize the JSON.
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	g.LDTKProject, err = ldtkgo.Open("example.ldtk", os.DirFS(dir))

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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	level := g.LDTKProject.Levels[g.CurrentLevel]

	screen.Fill(level.BGColor) // We want to use the BG Color when possible

	if g.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := level.BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		screen.DrawImage(g.BGImage, opt)
	}

	for i, layer := range g.EbitenRenderer.RenderedLayers {
		if g.ActiveLayers[i] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}

	// We'll additionally render the entities onscreen.
	for _, layer := range level.Layers {
		// In truth, we don't have to check to see if it's an entity layer before looping through,
		// because only Entity layers have entities in the Entities slice.
		for _, entity := range layer.Entities {

			if entity.TileRect != nil {

				tileset := g.EbitenRenderer.Tilesets[entity.TileRect.Tileset.Path]
				tileRect := entity.TileRect
				tile := tileset.SubImage(image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)).(*ebiten.Image)

				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(float64(entity.Position[0]), float64(entity.Position[1]))

				screen.DrawImage(tile, opt)

			}

		}

	}

}

func (g *Game) Layout(w, h int) (int, int) { return 320, 240 }

func main() {

	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetWindowTitle("LDtk-Go Example")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
