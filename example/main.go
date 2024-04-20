package main

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/renderer/ebitengine"
)

type Game struct {
	LDTKProject  *ldtkgo.Project
	Renderer     *renderer.Renderer
	BGImage      *ebiten.Image
	CurrentLevel int
	ActiveLayers []bool
}

//go:embed assets
var assets embed.FS

func NewGame() *Game {

	g := &Game{
		ActiveLayers: []bool{true, true, true, true, true},
	}

	proj, err := ldtkgo.Open("assets/example.ldtk", assets)

	if err != nil {
		panic(err)
	}

	g.LDTKProject = proj

	// Next, we create a new Renderer to render our level.

	// We pass the filesystem to use for the tileset - in this case, these are embedded into the code as an fs.FS / embed.FS, but you could
	// also use, say, os.DirFS() to create a filesystem from your HDD, as we did above.

	// Here, we use fs.Sub to create a sub directory that is relative to the .ldtk file.
	subDir, err := fs.Sub(assets, "assets")

	if err != nil {
		panic(err)
	}

	// The alternative would be to create another working fs.FS that is just fs.Sub(assets, "assets") - the root of it would
	// then be where the .ldtk file is (as all resources used in the .ldtk file are, naturally, relative to that location).
	// You could also use a custom function or something like debme to easily customize virtual file system roots to suit
	// your needs.

	// Here, we create the Renderer.
	g.Renderer, err = renderer.New(subDir, g.LDTKProject)

	if err != nil {
		panic(err)
	}

	fmt.Println("Press the 1 - 4 keys to toggle the tileset layers. Press the Left or Right arrow keys to switch Levels.")

	return g

}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.CurrentLevel++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.CurrentLevel--
	}

	if g.CurrentLevel >= len(g.LDTKProject.Levels) {
		g.CurrentLevel = 0
	}

	if g.CurrentLevel < 0 {
		g.CurrentLevel = len(g.LDTKProject.Levels) - 1
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.ActiveLayers[0] = !g.ActiveLayers[0]
	}

	// g.ActiveLayers[1] is the entities layer

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.ActiveLayers[2] = !g.ActiveLayers[2]
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.ActiveLayers[3] = !g.ActiveLayers[3]
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.ActiveLayers[4] = !g.ActiveLayers[4]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	level := g.LDTKProject.Levels[g.CurrentLevel]

	opt := renderer.NewDefaultDrawOptions()

	// Now, something that we can do that's a bit cool is that we can render things in the LayerDrawCallback - if we render on a specific
	// layer index or layer type, then we can render in-between the other layers, allowing us to place objects behind tiles or vice-versa
	opt.LayerDrawCallback = func(layer *ldtkgo.Layer, layerIndex int) bool {

		for _, entity := range layer.Entities {

			if entity.TileRect != nil {

				tileset := g.Renderer.Tilesets[entity.TileRect.Tileset.Path]

				tileRect := entity.TileRect
				tile := tileset.SubImage(image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)).(*ebiten.Image)

				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(float64(entity.Position[0]), float64(entity.Position[1]))

				screen.DrawImage(tile, opt)

			}

		}

		return g.ActiveLayers[layerIndex]

	}

	g.Renderer.Render(level, screen, opt)

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
