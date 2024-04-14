# LDtk-Go

LDtk-Go is a loader for [LDtk v0.9.3](https://ldtk.io/) projects written in pure Go.

![Screenshot](https://i.imgur.com/fFDmCCw.png)

Generally, you first load a project using `ldtkgo.Open()` or `ldtkgo.Read()`; afterward, you render the layers out to images, and then draw them onscreen with a rendering or game development framework, like [ebiten](https://hajimehoshi.github.io/ebiten/), [Pixel](https://github.com/faiface/pixel), [raylib-go](https://github.com/gen2brain/raylib-go), or [raylib-go-plus](https://github.com/Lachee/raylib-goplus). An example that uses ebiten for rendering is included.

All of the major elements of LDtk should be supported, including Levels, Layers, Tiles, AutoLayers, IntGrids, Entities, and Properties.

[Pkg.go.dev docs](https://pkg.go.dev/github.com/solarlune/ldtkgo)

## Loading Projects

Using LDtk-Go is not that simple. You should heck out the example to use.

## Running the Example

`cd` to the example directory, and then call `go run ./` to run the example (as it's self-contained within its directory). The console will give instructions for interacting with the example.

## Anything Else?

The core LDtk loader requires the `encoding/json` and `image` packages, as well as [tidwall's gjson](https://github.com/tidwall/gjson) package. The ebiten renderer requires ebiten package as well, of course.
