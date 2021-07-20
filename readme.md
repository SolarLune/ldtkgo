# LDtk-Go

LDtk-Go is a loader for [LDtk v0.9.3](https://ldtk.io/) projects written in pure Go.

![Screenshot](https://i.imgur.com/fFDmCCw.png)

Generally, you first load a project using `ldtkgo.Open()` or `ldtkgo.Read()`; afterward, you render the layers out to images, and then draw them onscreen with a rendering or game development framework, like [ebiten](https://hajimehoshi.github.io/ebiten/), [Pixel](https://github.com/faiface/pixel), [raylib-go](https://github.com/gen2brain/raylib-go), or [raylib-go-plus](https://github.com/Lachee/raylib-goplus). An example for both Ebiten and Raylib are included.

All of the major elements of LDtk should be supported, including Levels, Layers, Tiles, AutoLayers, IntGrids, Entities, and Properties.

[Pkg.go.dev docs](https://pkg.go.dev/github.com/solarlune/ldtkgo)

## Loading Projects

Using LDtk-Go is simple. 

First, you load the LDTK project, either with `ldtkgo.Open()` or `ldtkgo.Read()`. After that, you have access to most of the useful parts of the entire LDTK project. Here's an excerpt from the provided example, showing loading an LDtk Project and then 

```go

ldtkProject *ldtkgo.Project
ebitenRenderer *renderer.EbitenRenderer

func main() {

	// Load the LDtk Project
    ldtkProject, err := ldtkgo.LoadFile("example.ldtk")

    if err != nil {
        panic(err)
    }
    
	// Choose a level...
	level := ldtkProject.Levels[0]

	// Create a new renderer...
	ebitenRenderer = renderer.NewEbitenRenderer()

	// ... And render the tiles for the level out to layers, which will be *ebiten.Images. We'll retrieve them to draw in a Draw() loop later.
	ebitenRenderer.Render(level)

}

```

As shown above, rendering is done by creating a renderer and using it to render with your framework of choice. An example is already provided for the ebiten and raylib frameworks, which can be used as a general guide for creating other renderers.

## Running the Example

`cd` to the example directory, and then call `go run ./` to run the example (as it's self-contained within its directory). The console will give instructions for interacting with the example.

## Anything Else?

The core LDtk loader requires the `encoding/json` and `image` packages, as well as [tidwall's gjson](https://github.com/tidwall/gjson) package. The ebiten renderer requires ebiten package as well, of course.
