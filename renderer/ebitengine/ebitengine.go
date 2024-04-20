package ebitengine

// eb is a render system that uses ebiten to draw LDTK levels to the screen.

import (
	"errors"
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"

	_ "image/png" // Importing for loading PNGs
)

var ErrorBackgroundNotFound = "background image not found at given filepath"
var ErrorTilesetNotFound = "tileset image not found at given filepath"
var ErrorNoLevelGiven = "level pointer is nil"

// Ebitengine is a struct that draws LDtk levels to an *ebiten.screen.
type Ebitengine struct {
	Tilesets          map[string]*ebiten.Image
	Backgrounds       map[string]*ebiten.Image
	CurrentTileset    *ebiten.Image
	CurrentBackground *ebiten.Image
	FileSystem        fs.FS
}

// New creates a new Ebitengine renderer. This is used to render a level to one or more *ebiten.Images.
// The file system passed is the file system to use to load tileset images for the Renderer to use.
func New(fs fs.FS, project *ldtkgo.Project) (*Ebitengine, error) {

	renderer := &Ebitengine{
		Backgrounds: map[string]*ebiten.Image{},
		Tilesets:    map[string]*ebiten.Image{},
		FileSystem:  fs,
	}

	for _, level := range project.Levels {

		if level.BGImage == nil {
			continue
		}

		_, exists := renderer.Backgrounds[level.BGImage.Path]

		if !exists {
			img, _, err := ebitenutil.NewImageFromFileSystem(renderer.FileSystem, level.BGImage.Path)
			if err != nil {
				return nil, errors.New(ErrorBackgroundNotFound + ": [" + level.BGImage.Path + "]")
			}
			renderer.Backgrounds[level.BGImage.Path] = img
		}

	}

	for _, tileset := range project.Tilesets {

		_, exists := renderer.Tilesets[tileset.Path]

		if !exists {
			img, _, err := ebitenutil.NewImageFromFileSystem(renderer.FileSystem, tileset.Path)
			if err != nil {
				return nil, errors.New(ErrorTilesetNotFound + ": [" + tileset.Path + "]")
			}
			renderer.Tilesets[tileset.Path] = img
		}

	}

	return renderer, nil

}

type DrawOptions struct {
	BackgroundColorFill   bool                                      // Whether to fill the screen with the background color or not
	BackgroundDraw        bool                                      // Whether to render the background image when drawing the ldtkgo.Level
	BackgroundDrawOptions *ebiten.DrawImageOptions                  // The options to use when drawing the background
	LayerDrawOptions      *ebiten.DrawImageOptions                  // The options to use when drawing the tile layers
	LayerDrawCallback     func(layer *ldtkgo.Layer, index int) bool // A callback that is called in the order in which layers are rendering. Returning false stops the current layer from rendering and moves on.
}

// NewDefaultDrawOptions creates a RenderOptions struct with the default set of render options.
func NewDefaultDrawOptions() *DrawOptions {
	return &DrawOptions{
		BackgroundColorFill:   true,
		BackgroundDraw:        true,
		BackgroundDrawOptions: &ebiten.DrawImageOptions{},
		LayerDrawOptions:      &ebiten.DrawImageOptions{},
	}
}

// Render draws an *ldtkgo.Level to the destination screen specified using render options to control the process.
func (e *Ebitengine) Render(level *ldtkgo.Level, screen *ebiten.Image, drawOptions *DrawOptions) error {

	if level == nil {
		return errors.New(ErrorNoLevelGiven)
	}

	if drawOptions == nil {
		drawOptions = NewDefaultDrawOptions()
	}

	if drawOptions.BackgroundColorFill {
		screen.Fill(level.BGColor) // We want to use the BG Color when possible
	}

	if drawOptions.BackgroundDraw && level.BGImage != nil && level.BGImage.Path != "" {
		e.CurrentBackground = e.Backgrounds[level.BGImage.Path]
		opt := *drawOptions.BackgroundDrawOptions
		opt.GeoM.Translate(-level.BGImage.CropRect[0], -level.BGImage.CropRect[1])
		opt.GeoM.Scale(level.BGImage.ScaleX, level.BGImage.ScaleY)
		screen.DrawImage(e.CurrentBackground, &opt)
	}

	// Reverse sort the layers when drawing because in LDtk, the numbering order is from top-to-bottom, but the drawing order is from bottom-to-top.
	for layerIndex := len(level.Layers) - 1; layerIndex >= 0; layerIndex-- {

		layer := level.Layers[layerIndex]

		if drawOptions.LayerDrawCallback != nil {
			if !drawOptions.LayerDrawCallback(layer, layerIndex) {
				continue
			}
		}

		if tiles := layer.AllTiles(); len(tiles) > 0 {

			e.CurrentTileset = e.Tilesets[layer.Tileset.Path]

			for _, tileData := range tiles {
				// er.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)

				// Subimage the Tile from the Tileset
				tile := e.CurrentTileset.SubImage(image.Rect(tileData.Src[0], tileData.Src[1], tileData.Src[0]+layer.GridSize, tileData.Src[1]+layer.GridSize)).(*ebiten.Image)

				opt := *drawOptions.LayerDrawOptions // Clone the draw options used to render the tiles, because we'll be transforming them

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
				screen.DrawImage(tile, &opt)

			}

		}

	}

	return nil

}
