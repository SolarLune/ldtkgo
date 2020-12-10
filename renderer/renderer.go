//Package renderer holds an ebiten Renderer for LDtk Projects.
package renderer

import (
	"image"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"

	_ "image/png" // Importing for loading PNGs
)

// TilesetLoader represents an interface that can be implemented to load a tileset from string, returning an *ebiten.Image.
type TilesetLoader interface {
	LoadTileset(string) *ebiten.Image
}

// DiskLoader is an implementation of TilesetLoader that simply loads a Tileset image from disk using ebitenutil's NewImageFromFile() function.
type DiskLoader struct {
	Filter ebiten.Filter
}

// NewDiskLoader creates a new DiskLoader struct.
func NewDiskLoader() *DiskLoader {
	return &DiskLoader{Filter: ebiten.FilterNearest}
}

// LoadTileset simply loads a Tileset image from disk using ebitenutil's NewImageFromFile() function.
func (d *DiskLoader) LoadTileset(tilesetPath string) *ebiten.Image {
	if img, _, err := ebitenutil.NewImageFromFile(tilesetPath); err == nil {
		return img
	}
	return nil
}

// EbitenRenderer is a struct that renders LDtk levels to *ebiten.Images.
type EbitenRenderer struct {
	Tilesets       map[string]*ebiten.Image
	CurrentTileset string
	RenderedLayers []*ebiten.Image
	Loader         TilesetLoader // Loader for the renderer; defaults to a DiskLoader instance, though this can be switched out with something else as necessary.
}

// NewEbitenRenderer creates a new Renderer instance. Width and Height should be the width and height of the level in pixels.
func NewEbitenRenderer() *EbitenRenderer {

	return &EbitenRenderer{
		Tilesets:       map[string]*ebiten.Image{},
		RenderedLayers: []*ebiten.Image{},
		Loader:         NewDiskLoader(),
	}

}

// Clear clears the renderer's Result.
func (er *EbitenRenderer) Clear() {
	for _, layer := range er.RenderedLayers {
		layer.Dispose()
	}
	er.RenderedLayers = []*ebiten.Image{}
}

// setTileset gets called when necessary between rendering indidvidual Layers of a Level.
func (er *EbitenRenderer) setTileset(tilesetPath string, w, h int) {

	_, exists := er.Tilesets[tilesetPath]

	if !exists {
		er.Tilesets[tilesetPath] = er.Loader.LoadTileset(tilesetPath)
	}

	er.CurrentTileset = tilesetPath

	newLayer := ebiten.NewImage(w, h)
	er.RenderedLayers = append(er.RenderedLayers, newLayer)

}

// renderTile gets called by LDtkgo.Layer.RenderTiles(), and is currently provided the following arguments to handle rendering each tile in a Layer:
// x, y = position of the drawn tile
// srcX, srcY = position on the source tilesheet of the specified tile
// srcW, srcH = width and height of the tile
// flipX, flipY = If the tile should be flipped horizontally or vertically
func (er *EbitenRenderer) renderTile(x, y, srcX, srcY, srcW, srcH int, flipBit byte) {

	// Subimage the Tile from the Tileset
	tile := er.Tilesets[er.CurrentTileset].SubImage(image.Rect(srcX, srcY, srcX+srcW, srcY+srcH)).(*ebiten.Image)

	opt := &ebiten.DrawImageOptions{}

	// We have to offset the tile to be centered before flipping
	opt.GeoM.Translate(float64(-srcW/2), float64(-srcH/2))

	// Handle flipping; first bit in byte is horizontal flipping, second is vertical flipping.

	if flipBit&1 > 0 {
		opt.GeoM.Scale(-1, 1)
	}
	if flipBit&2 > 0 {
		opt.GeoM.Scale(1, -1)
	}

	// Undo offsetting
	opt.GeoM.Translate(float64(srcW/2), float64(srcH/2))

	// Move tile to final position; note that slightly unlike LDtk, layer offsets in LDtk-Go are added directly into the final tiles' X and Y positions. This means that with this renderer,
	// if a layer's offset pushes tiles outside of the layer's render Result image, they will be cut off. On LDtk, the tiles are still rendered, of course.
	opt.GeoM.Translate(float64(x), float64(y))

	// Finally, draw the tile to the Result image.
	er.RenderedLayers[len(er.RenderedLayers)-1].DrawImage(tile, opt)

}

// Render clears, and then renders out each visible Layer in an ldtgo.Level instance.
func (er *EbitenRenderer) Render(level *ldtkgo.Level) {

	er.Clear()

	for _, layer := range level.Layers {

		switch layer.Type {

		case ldtkgo.LayerTypeTile:

			if len(layer.Tiles) > 0 {

				er.setTileset(layer.TilesetPath, level.Width, level.Height)

				for _, tile := range layer.Tiles {
					er.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)
				}

			}

		case ldtkgo.LayerTypeIntGrid: // IntGrids get autotiles automatically
			fallthrough
		case ldtkgo.LayerTypeAutoTile:

			if len(layer.AutoTiles) > 0 {

				er.setTileset(layer.TilesetPath, level.Width, level.Height)

				for _, tile := range layer.AutoTiles {
					er.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)
				}

			}

		}

	}

	// Reverse sort the layers when drawing because in LDtk, the numbering order is from top-to-bottom, but the drawing order is from bottom-to-top.
	sort.Slice(er.RenderedLayers, func(i, j int) bool {
		return i > j
	})

}
