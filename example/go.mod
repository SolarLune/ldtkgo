module github.com/solarlune/ldtkgo-example

go 1.14

require (
	github.com/hajimehoshi/ebiten v1.12.4
	github.com/solarlune/ldtkgo v0.0.0-00010101000000-000000000000
	github.com/solarlune/ldtkgo/renderer v0.0.0-00010101000000-000000000000
)

replace github.com/solarlune/ldtkgo => ../
replace github.com/solarlune/ldtkgo/renderer => ../renderer
