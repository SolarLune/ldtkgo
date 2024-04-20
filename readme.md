# LDtk-Go

LDtk-Go is a loader for [LDtk v0.9.3](https://ldtk.io/) projects written in pure Go.

![Screenshot](https://i.imgur.com/fFDmCCw.png)

Generally, you first load a project using `ldtkgo.Open()` or `ldtkgo.Read()`; afterward, you render the layers out to images, and then draw them onscreen with a rendering or game development framework, like [ebiten](https://hajimehoshi.github.io/ebiten/), [Pixel](https://github.com/faiface/pixel), [raylib-go](https://github.com/gen2brain/raylib-go), or [raylib-go-plus](https://github.com/Lachee/raylib-goplus). An example that uses ebiten for rendering is included.

All of the major elements of LDtk should be supported, including Levels, Layers, Tiles, AutoLayers, IntGrids, Entities, and Properties.

[Pkg.go.dev docs](https://pkg.go.dev/github.com/solarlune/ldtkgo)

## Loading Projects

Using LDtk-Go is simple. 

First, you load the LDTK project, either with `ldtkgo.Open()` or `ldtkgo.Read()`. After that, you have access to most of the useful parts of the entire LDTK project. You can then render the map using the included renderer if you're using [Ebitengine](https://ebitengine.org/).

```go
package main

import (
	"embed"

	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/renderer/ebitengine"
)

var ldtkProject *ldtkgo.Project
var ebitenRenderer *renderer.Ebitengine

//go:embed *.ldtk *.png
var fileSystem embed.FS

func main() {
	// Load the LDtk Project
	ldtkProject, err := ldtkgo.Open("example.ldtk", fileSystem)

	if err != nil {
		panic(err)
	}

	// ldtkProject now contains all data from the file.
	// If you'd like to render it, you could use the included renderer that uses Ebitengine:

	// Choose a level...
	level := ldtkProject.Levels[0]

	// Create a new renderer...
	ebitenRenderer = renderer.New(ldtkProject)

	// ... And render the tiles for the level out to layers, which will be *ebiten.Images. We can then retrieve them to draw in a Draw() loop later.
	ebitenRenderer.Render(level)
}


```

## Running the Example

`cd` to the example directory, and then call `go run .` to run the example (as it's self-contained within its directory). The console will give instructions for interacting with the example.

## Anything Else?

The core LDtk loader requires the `encoding/json` and `image` package. The Ebiten renderer requires Ebitengine as well, of course.

## To-do

- [ ] Add map clipping / viewports to Ebitengine renderer