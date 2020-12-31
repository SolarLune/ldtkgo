package raylibrenderer

import (
	"path/filepath"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/ldtkgo"
)

// TilesetLoader represents an interface that can be implemented to load a tileset from string, returning an *ebiten.Image.
type TilesetLoader interface {
	LoadTileset(string) rl.Texture2D
}

// DiskLoader is an implementation of TilesetLoader that simply loads a Tileset image from disk using ebitenutil's NewImageFromFile() function.
type DiskLoader struct {
	BasePath string
}

// NewDiskLoader creates a new DiskLoader struct.
func NewDiskLoader(basePath string) *DiskLoader {
	return &DiskLoader{
		BasePath: basePath,
	}
}

// LoadTileset simply loads a Tileset image from disk using ebitenutil's NewImageFromFile() function.
func (d *DiskLoader) LoadTileset(tilesetPath string) rl.Texture2D {
	return rl.LoadTexture(filepath.Join(d.BasePath, tilesetPath))
}

// RenderedLayer represents an LDtk.Layer that was rendered out to a raylib.RenderTexture2D.
type RenderedLayer struct {
	Image rl.RenderTexture2D // The image that was rendered out
	Layer *ldtkgo.Layer      // The layer used to render the image
}

// RaylibRenderer is a struct that renders LDtk levels to raylib.RenderTexture2D images.
type RaylibRenderer struct {
	Tilesets       map[string]rl.Texture2D
	CurrentTileset string
	RenderedLayers []*RenderedLayer
	Loader         TilesetLoader // Loader for the renderer; defaults to a DiskLoader instance, though this can be switched out with something else as necessary.
}

// NewRaylibRenderer creates a new Renderer instance. TilesetLoader should be an instance of a struct designed to return *ebiten.Images for each Tileset requested (by path relative to the LDtk project file).
func NewRaylibRenderer(loader TilesetLoader) *RaylibRenderer {

	return &RaylibRenderer{
		Tilesets:       map[string]rl.Texture2D{},
		RenderedLayers: []*RenderedLayer{},
		Loader:         loader,
	}

}

// Clear clears the renderer's Result.
func (raylibRenderer *RaylibRenderer) Clear() {
	for _, layer := range raylibRenderer.RenderedLayers {
		rl.UnloadRenderTexture(layer.Image)
	}
	raylibRenderer.RenderedLayers = []*RenderedLayer{}
}

// beginLayer gets called when necessary between rendering indidvidual Layers of a Level.
func (raylibRenderer *RaylibRenderer) beginLayer(layer *ldtkgo.Layer, w, h int) {

	tilesetPath := layer.TilesetPath

	_, exists := raylibRenderer.Tilesets[tilesetPath]

	if !exists {
		raylibRenderer.Tilesets[tilesetPath] = raylibRenderer.Loader.LoadTileset(tilesetPath)
	}

	raylibRenderer.CurrentTileset = tilesetPath

	renderedImage := rl.LoadRenderTexture(int32(w), int32(h))

	raylibRenderer.RenderedLayers = append(raylibRenderer.RenderedLayers, &RenderedLayer{Image: renderedImage, Layer: layer})

}

// renderTile gets called by LDtkgo.Layer.RenderTiles(), and is currently provided the following arguments to handle rendering each tile in a Layer:
// x, y = position of the drawn tile
// srcX, srcY = position on the source tilesheet of the specified tile
// srcW, srcH = width and height of the tile
// flipBit = the flip bit of the tile; if the first bit is set, it should flip horizontally. If the second is set, it should flip vertically.
func (raylibRenderer *RaylibRenderer) renderTile(x, y, srcX, srcY, srcW, srcH int, flipBit byte) {

	tileset := raylibRenderer.Tilesets[raylibRenderer.CurrentTileset]

	src := rl.NewRectangle(float32(srcX), float32(srcY), float32(srcW), float32(srcH))
	dst := src
	dst.X = float32(x)
	dst.Y = float32(y)

	if flipBit&1 > 0 {
		src.Width *= -1
	}
	if flipBit&2 > 0 {
		src.Height *= -1
	}

	rl.DrawTexturePro(tileset, src, dst, rl.Vector2{}, 0, rl.White)

}

// Render clears, and then renders out each visible Layer in an ldtgo.Level instance.
func (raylibRenderer *RaylibRenderer) Render(level *ldtkgo.Level) {

	raylibRenderer.Clear()

	rl.BeginDrawing()

	for _, layer := range level.Layers {

		switch layer.Type {

		case ldtkgo.LayerTypeTile:

			if len(layer.Tiles) > 0 {

				raylibRenderer.beginLayer(layer, level.Width, level.Height)

				target := raylibRenderer.RenderedLayers[len(raylibRenderer.RenderedLayers)-1]
				rl.BeginTextureMode(target.Image)

				for _, tile := range layer.Tiles {
					raylibRenderer.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)
				}

				rl.EndTextureMode()

			}

		case ldtkgo.LayerTypeIntGrid: // IntGrids get autotiles automatically
			fallthrough
		case ldtkgo.LayerTypeAutoTile:

			if len(layer.AutoTiles) > 0 {

				raylibRenderer.beginLayer(layer, level.Width, level.Height)

				target := raylibRenderer.RenderedLayers[len(raylibRenderer.RenderedLayers)-1]
				rl.BeginTextureMode(target.Image)

				for _, tile := range layer.AutoTiles {
					raylibRenderer.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)
				}

				rl.EndTextureMode()

			}

		}

	}

	rl.EndDrawing()

	// Reverse sort the layers when drawing because in LDtk, the numbering order is from top-to-bottom, but the drawing order is from bottom-to-top.
	sort.Slice(raylibRenderer.RenderedLayers, func(i, j int) bool {
		return i > j
	})

}
