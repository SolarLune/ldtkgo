package main

import (
	"image"
	"path/filepath"
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
	BasePath string
	Filter   ebiten.Filter
}

// NewDiskLoader creates a new DiskLoader struct.
func NewDiskLoader(basePath string) *DiskLoader {
	return &DiskLoader{
		BasePath: basePath,
		Filter:   ebiten.FilterNearest,
	}
}

// LoadTileset simply loads a Tileset image from disk using ebitenutil's NewImageFromFile() function.
func (d *DiskLoader) LoadTileset(tilesetPath string) *ebiten.Image {
	if img, _, err := ebitenutil.NewImageFromFile(filepath.Join(d.BasePath, tilesetPath)); err == nil {
		return img
	}
	return nil
}

// RenderedLayer represents an LDtk.Layer that was rendered out to an *ebiten.Image.
type RenderedLayer struct {
	Image *ebiten.Image // The image that was rendered out
	Layer *ldtkgo.Layer // The layer used to render the image
}

// EbitenRenderer is a struct that renders LDtk levels to *ebiten.Images.
type EbitenRenderer struct {
	Tilesets       map[string]*ebiten.Image
	CurrentTileset string
	RenderedLayers []*RenderedLayer
	Loader         TilesetLoader // Loader for the renderer; defaults to a DiskLoader instance, though this can be switched out with something else as necessary.
}

// NewEbitenRenderer creates a new Renderer instance. TilesetLoader should be an instance of a struct designed to return *ebiten.Images for each Tileset requested (by path relative to the LDtk project file).
func NewEbitenRenderer(loader TilesetLoader) *EbitenRenderer {

	return &EbitenRenderer{
		Tilesets:       map[string]*ebiten.Image{},
		RenderedLayers: []*RenderedLayer{},
		Loader:         loader,
	}

}

// Clear clears the renderer's Result.
func (er *EbitenRenderer) Clear() {
	for _, layer := range er.RenderedLayers {
		layer.Image.Dispose()
	}
	er.RenderedLayers = []*RenderedLayer{}
}

// beginLayer gets called when necessary between rendering indidvidual Layers of a Level.
func (er *EbitenRenderer) beginLayer(layer *ldtkgo.Layer, w, h int) {

	_, exists := er.Tilesets[layer.Tileset.Path]

	if !exists {
		er.Tilesets[layer.Tileset.Path] = er.Loader.LoadTileset(layer.Tileset.Path)
	}

	er.CurrentTileset = layer.Tileset.Path

	renderedImage := ebiten.NewImage(w, h)

	er.RenderedLayers = append(er.RenderedLayers, &RenderedLayer{Image: renderedImage, Layer: layer})

}

// renderTile gets called by LDtkgo.Layer.RenderTiles(), and is currently provided the following arguments to handle rendering each tile in a Layer:
// x, y = position of the drawn tile
// srcX, srcY = position on the source tilesheet of the specified tile
// srcW, srcH = width and height of the tile
// flipBit = the flip bit of the tile; if the first bit is set, it should flip horizontally. If the second is set, it should flip vertically.
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
	er.RenderedLayers[len(er.RenderedLayers)-1].Image.DrawImage(tile, opt)

}

// Render clears, and then renders out each visible Layer in an ldtgo.Level instance.
func (er *EbitenRenderer) Render(level *ldtkgo.Level) {

	er.Clear()

	for _, layer := range level.Layers {

		switch layer.Type {

		case ldtkgo.LayerTypeIntGrid: // IntGrids get autotiles automatically
			fallthrough
		case ldtkgo.LayerTypeAutoTile:
			fallthrough
		case ldtkgo.LayerTypeTile:

			if tiles := layer.AllTiles(); len(tiles) > 0 {

				er.beginLayer(layer, level.Width, level.Height)

				for _, tileData := range tiles {
					// er.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)

					// Subimage the Tile from the Tileset
					tile := er.Tilesets[er.CurrentTileset].SubImage(image.Rect(tileData.Src[0], tileData.Src[1], tileData.Src[0]+layer.GridSize, tileData.Src[1]+layer.GridSize)).(*ebiten.Image)

					opt := &ebiten.DrawImageOptions{}

					// We have to offset the tile to be centered before flipping
					opt.GeoM.Translate(float64(-layer.GridSize/2), float64(-layer.GridSize/2))

					// Handle flipping; first bit in byte is horizontal flipping, second is vertical flipping.

					if tileData.FlipX() {
						opt.GeoM.Scale(-1, 1)
					}
					if tileData.FlipY() {
						opt.GeoM.Scale(1, -1)
					}

					// Undo offsetting
					opt.GeoM.Translate(float64(layer.GridSize/2), float64(layer.GridSize/2))

					// Move tile to final position; note that slightly unlike LDtk, layer offsets in LDtk-Go are added directly into the final tiles' X and Y positions. This means that with this renderer,
					// if a layer's offset pushes tiles outside of the layer's render Result image, they will be cut off. On LDtk, the tiles are still rendered, of course.
					opt.GeoM.Translate(float64(tileData.Position[0]+layer.OffsetX), float64(tileData.Position[1]+layer.OffsetY))

					// Finally, draw the tile to the Result image.
					er.RenderedLayers[len(er.RenderedLayers)-1].Image.DrawImage(tile, opt)

				}

			}

		}

	}

	// Reverse sort the layers when drawing because in LDtk, the numbering order is from top-to-bottom, but the drawing order is from bottom-to-top.
	sort.Slice(er.RenderedLayers, func(i, j int) bool {
		return i > j
	})

}
