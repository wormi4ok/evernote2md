module github.com/wormi4ok/evernote2md

go 1.24

require (
	github.com/bmatcuk/doublestar/v4 v4.8.1
	github.com/briandowns/spinner v1.23.2
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/hashicorp/logutils v1.0.0
	github.com/integrii/flaggy v1.5.2
	github.com/mattn/godown v0.0.2-0.20210508133137-72c48840c3e3
	github.com/sergi/go-diff v1.4.0
	golang.org/x/net v0.41.0
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
)

replace github.com/mattn/godown => github.com/wormi4ok/godown v0.5.0
