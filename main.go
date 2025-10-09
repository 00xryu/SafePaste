package main

import (
	"image/color"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"safe-paste/safe_paste"
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
	th := material.NewTheme() // Basit tema
	var ops op.Ops

	// Widget'lar
	var inputEditor widget.Editor
	inputEditor.SingleLine = false
	inputEditor.Submit = true

	var outputText string

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

			// Layout: Dikey flex, spacing düzeltildi
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H3(th, "Orijinal Text:").Layout(gtx)
				}),
				layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, &inputEditor, "Metni buraya yapıştır...")
					return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, ed.Layout)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Butonlar yatay, spacing düzeltildi
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &maskButton, "Maskele")
							btn.Layout(gtx)
							if maskButton.Clicked(gtx) {
								outputText = safe_paste.maskText(inputEditor.Text())
							}
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &copyButton, "Kopyala")
							btn.Layout(gtx)
							if copyButton.Clicked(gtx) {
								if outputText == "" {
									return layout.Dimensions{}
								}
								if runtime.GOOS == "windows" {
									cmd := exec.Command("powershell", "-Command", "Set-Clipboard -Value '"+strings.ReplaceAll(outputText, "'", "''")+"'")
									cmd.Run()
								} else {
									cmd := exec.Command("xclip", "-selection", "clipboard")
									cmd.Stdin = strings.NewReader(outputText)
									cmd.Run()
								}
							}
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &settingsButton, "Ayarlar")
							btn.Layout(gtx)
							if settingsButton.Clicked(gtx) {
								if runtime.GOOS == "windows" {
									exec.Command("notepad.exe", "config.json").Start()
								} else {
									exec.Command("xdg-open", "config.json").Start()
								}
							}
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H3(th, "Maskelenmiş:").Layout(gtx)
				}),
				layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), outputText)
					lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
					lbl.MaxLines = -1
					return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, lbl.Layout)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
