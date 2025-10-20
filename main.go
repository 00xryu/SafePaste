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

	// Widget'lar
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

			// Ana layout - tüm içeriği padding ile çevrele
			layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(20),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
					// Başlık
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H5(th, "Orijinal Text")
						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
					}),
					// Input editor
					layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, &inputEditor, "Metni buraya yapıştır...")
						ed.TextSize = unit.Sp(14)
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
					}),
					// Boşluk
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
					// Butonlar
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Flexed(0.33, func(gtx layout.Context) layout.Dimensions {
								if maskButton.Clicked(gtx) {
									input := inputEditor.Text()
									log.Println("Input text:", input)
									masked := sp.MaskText(input)
									log.Println("Output text:", masked)
									outputEditor.SetText(masked)
								}
								btn := material.Button(th, &maskButton, "Maskele")
								return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, btn.Layout)
							}),
							layout.Flexed(0.33, func(gtx layout.Context) layout.Dimensions {
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
								btn := material.Button(th, &copyButton, "Kopyala")
								return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, btn.Layout)
							}),
							layout.Flexed(0.33, func(gtx layout.Context) layout.Dimensions {
								if settingsButton.Clicked(gtx) {
									// Config dosyasının yolunu bul (exe ile aynı dizinde)
									exePath, _ := os.Executable()
									exeDir := filepath.Dir(exePath)
									configPath := filepath.Join(exeDir, "config.json")

									if runtime.GOOS == "windows" {
										exec.Command("notepad.exe", configPath).Start()
									} else {
										exec.Command("xdg-open", configPath).Start()
									}
								}
								btn := material.Button(th, &settingsButton, "Ayarlar")
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, btn.Layout)
							}),
						)
					}),
					// Boşluk
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
					// Başlık
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H5(th, "Maskelenmiş Text")
						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
					}),
					// Output editor
					layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, &outputEditor, "")
						ed.TextSize = unit.Sp(14)
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
