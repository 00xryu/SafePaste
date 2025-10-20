package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	sp "safe-paste/safe_paste"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("SafePaste"))
		window.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	// Widgets
	var inputEditor widget.Editor
	inputEditor.SingleLine = false
	inputEditor.Submit = true

	var outputEditor widget.Editor
	outputEditor.ReadOnly = true
	outputEditor.SingleLine = false

	var maskButton widget.Clickable
	var copyButton widget.Clickable
	var settingsButton widget.Clickable

	for {
		e := window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Main layout with padding
			layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(20),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Horizontal layout: Left panel | Buttons | Right panel
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
					// Left panel - Original text
					layout.Flexed(0.45, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							// Title
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, "Original")
								return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
							}),
							// Input editor with border
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								// Border using background
								border := widget.Border{
									Color:        material.NewTheme().Palette.Fg,
									CornerRadius: unit.Dp(8),
									Width:        unit.Dp(1),
								}
								return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									ed := material.Editor(th, &inputEditor, "Paste your text here...")
									ed.TextSize = unit.Sp(14)
									return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
								})
							}),
						)
					}),
					// Center - Buttons (vertical)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								// Spacer to center buttons vertically
								layout.Flexed(1, layout.Spacer{}.Layout),
								// Mask button
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if maskButton.Clicked(gtx) {
										input := inputEditor.Text()
										log.Println("Input text:", input)
										masked := sp.MaskText(input)
										log.Println("Output text:", masked)
										outputEditor.SetText(masked)
									}
									btn := material.Button(th, &maskButton, "Mask")
									return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, btn.Layout)
								}),
								// Copy button
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if copyButton.Clicked(gtx) {
										text := outputEditor.Text()
										log.Println("Copying text:", text)
										if runtime.GOOS == "windows" {
											cmd := exec.Command("cmd", "/c", "echo "+text+" | clip")
											cmd.Run()
										} else {
											cmd := exec.Command("xclip", "-selection", "clipboard")
											cmd.Stdin = strings.NewReader(text)
											cmd.Run()
										}
									}
									btn := material.Button(th, &copyButton, "Copy")
									return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, btn.Layout)
								}),
								// Settings button
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if settingsButton.Clicked(gtx) {
										// Find config file path (in same directory as exe)
										exePath, _ := os.Executable()
										exeDir := filepath.Dir(exePath)
										configPath := filepath.Join(exeDir, "config.json")

										if runtime.GOOS == "windows" {
											exec.Command("notepad.exe", configPath).Start()
										} else {
											exec.Command("xdg-open", configPath).Start()
										}
									}
									btn := material.Button(th, &settingsButton, "Settings")
									return btn.Layout(gtx)
								}),
								// Spacer to center buttons vertically
								layout.Flexed(1, layout.Spacer{}.Layout),
							)
						})
					}),
					// Right panel - Masked text
					layout.Flexed(0.45, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							// Title
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th, "Output")
								return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
							}),
							// Output editor with border
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								// Border using background
								border := widget.Border{
									Color:        material.NewTheme().Palette.Fg,
									CornerRadius: unit.Dp(8),
									Width:        unit.Dp(1),
								}
								return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									ed := material.Editor(th, &outputEditor, "")
									ed.TextSize = unit.Sp(14)
									return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
								})
							}),
						)
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
