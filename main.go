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
		window.Option(app.Size(unit.Dp(1200), unit.Dp(800)))
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

	var aiInputEditor widget.Editor
	aiInputEditor.SingleLine = false
	aiInputEditor.Submit = true

	var aiOutputEditor widget.Editor
	aiOutputEditor.ReadOnly = true
	aiOutputEditor.SingleLine = false

	var maskButton widget.Clickable
	var unmaskButton widget.Clickable
	var copyMaskedButton widget.Clickable
	var copyUnmaskedButton widget.Clickable
	var settingsButton widget.Clickable

	// Store mapping for unmasking
	var currentMapping map[string]string

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
				// Vertical layout: Top row | Bottom row
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Top row - Original & Masked
					layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							// Top-Left: Original text
							layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.H6(th, "Original")
										return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, title.Layout)
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
							// Center buttons for top row
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
										layout.Flexed(1, layout.Spacer{}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if maskButton.Clicked(gtx) {
												input := inputEditor.Text()
												result := sp.MaskTextWithMapping(input)
												outputEditor.SetText(result.MaskedText)
												currentMapping = result.Mapping
												log.Println("Masked. Mapping size:", len(currentMapping))
											}
											btn := material.Button(th, &maskButton, "Mask →")
											return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, btn.Layout)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if copyMaskedButton.Clicked(gtx) {
												text := outputEditor.Text()
												copyToClipboard(text)
												log.Println("Copied masked text")
											}
											btn := material.Button(th, &copyMaskedButton, "Copy")
											return btn.Layout(gtx)
										}),
										layout.Flexed(1, layout.Spacer{}.Layout),
									)
								})
							}),
							// Top-Right: Masked text
							layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.H6(th, "Masked")
										return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, title.Layout)
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
					}),
					// Middle spacing
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
					// Bottom row - AI Input & Unmasked
					layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							// Bottom-Left: AI response (to unmask)
							layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.H6(th, "AI Response")
										return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, title.Layout)
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										border := widget.Border{
											Color:        material.NewTheme().Palette.Fg,
											CornerRadius: unit.Dp(8),
											Width:        unit.Dp(1),
										}
										return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											ed := material.Editor(th, &aiInputEditor, "Paste AI response here...")
											ed.TextSize = unit.Sp(14)
											return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
										})
									}),
								)
							}),
							// Center buttons for bottom row
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
										layout.Flexed(1, layout.Spacer{}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if unmaskButton.Clicked(gtx) {
												aiInput := aiInputEditor.Text()
												if currentMapping != nil {
													unmasked := sp.UnmaskText(aiInput, currentMapping)
													aiOutputEditor.SetText(unmasked)
													log.Println("Unmasked with", len(currentMapping), "mappings")
												} else {
													log.Println("No mapping available. Mask text first!")
												}
											}
											btn := material.Button(th, &unmaskButton, "Unmask →")
											return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, btn.Layout)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if copyUnmaskedButton.Clicked(gtx) {
												text := aiOutputEditor.Text()
												copyToClipboard(text)
												log.Println("Copied unmasked text")
											}
											btn := material.Button(th, &copyUnmaskedButton, "Copy")
											return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, btn.Layout)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if settingsButton.Clicked(gtx) {
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
										layout.Flexed(1, layout.Spacer{}.Layout),
									)
								})
							}),
							// Bottom-Right: Unmasked result
							layout.Flexed(0.48, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.H6(th, "Unmasked")
										return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, title.Layout)
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										border := widget.Border{
											Color:        material.NewTheme().Palette.Fg,
											CornerRadius: unit.Dp(8),
											Width:        unit.Dp(1),
										}
										return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											ed := material.Editor(th, &aiOutputEditor, "")
											ed.TextSize = unit.Sp(14)
											return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
										})
									}),
								)
							}),
						)
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}

// Helper function for clipboard
func copyToClipboard(text string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "echo "+text+" | clip")
		cmd.Run()
	} else {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		cmd.Run()
	}
}
