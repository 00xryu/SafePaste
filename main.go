package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	sp "safe-paste/safe_paste"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const Version = "v1.0.0"

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
	// Load initial config
	cfg := sp.LoadConfig()
	isDark := cfg.Theme == "dark"

	th := material.NewTheme()
	updateTheme(th, isDark)

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
	var clearMaskButton widget.Clickable
	var clearUnmaskButton widget.Clickable
	var themeSwitchButton widget.Clickable

	// Animation state
	var animProgress float32
	if isDark {
		animProgress = 1.0
	}
	var lastTime time.Time

	// Store mapping for unmasking
	var currentMapping map[string]string

	for {
		e := window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Fill background
			paint.Fill(gtx.Ops, th.Palette.Bg)

			// Animation logic
			target := float32(0.0)
			if isDark {
				target = 1.0
			}

			now := e.Now
			if !lastTime.IsZero() {
				dt := float32(now.Sub(lastTime).Seconds())
				// Smooth transition
				diff := target - animProgress
				if math.Abs(float64(diff)) > 0.01 {
					animProgress += diff * dt * 10 // Speed factor
					window.Invalidate()
				} else {
					animProgress = target
				}
			}
			lastTime = now

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
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
											Color:        th.Fg,
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
											return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, btn.Layout)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if clearMaskButton.Clicked(gtx) {
												inputEditor.SetText("")
												outputEditor.SetText("")
												currentMapping = nil
												log.Println("Cleared masked section")
											}
											btn := material.Button(th, &clearMaskButton, "Clear")
											btn.Background = color.NRGBA{R: 0xFF, G: 0x88, B: 0x88, A: 0xFF}
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
											Color:        th.Fg,
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
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
											Color:        th.Fg,
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
											if clearUnmaskButton.Clicked(gtx) {
												aiInputEditor.SetText("")
												aiOutputEditor.SetText("")
												log.Println("Cleared unmasked section")
											}
											btn := material.Button(th, &clearUnmaskButton, "Clear")
											btn.Background = color.NRGBA{R: 0xFF, G: 0x88, B: 0x88, A: 0xFF}
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
											Color:        th.Fg,
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
					// Footer
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							// Theme Switcher
							if themeSwitchButton.Clicked(gtx) {
								isDark = !isDark
								updateTheme(th, isDark)
								cfg.Theme = "light"
								if isDark {
									cfg.Theme = "dark"
								}
								sp.SaveConfig(cfg)
								window.Invalidate()
							}

							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Clickable(gtx, &themeSwitchButton, func(gtx layout.Context) layout.Dimensions {
										return drawThemeIcon(gtx, th.Fg, animProgress)
									})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Body2(th, "SafePaste "+Version)
									lbl.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
									return layout.Center.Layout(gtx, lbl.Layout)
								}),
								// Spacer to balance the left button
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{Size: image.Point{X: 40, Y: 40}}
								}),
							)
						})
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}

func updateTheme(th *material.Theme, isDark bool) {
	if isDark {
		th.Palette.Fg = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		th.Palette.Bg = color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff}
		th.Palette.ContrastBg = color.NRGBA{R: 0x3e, G: 0x3e, B: 0x3e, A: 0xff}
		th.Palette.ContrastFg = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	} else {
		th.Palette.Fg = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		th.Palette.Bg = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		th.Palette.ContrastBg = color.NRGBA{R: 0x3f, G: 0x51, B: 0xb5, A: 0xff}
		th.Palette.ContrastFg = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	}
}

func drawThemeIcon(gtx layout.Context, col color.NRGBA, progress float32) layout.Dimensions {
	size := unit.Dp(30)
	sz := gtx.Dp(size)
	center := f32.Pt(float32(sz)/2, float32(sz)/2)

	// Helper to draw a line
	drawLine := func(start, end f32.Point) {
		path := clip.Path{}
		path.Begin(gtx.Ops)
		path.MoveTo(start)
		path.LineTo(end)
		paint.FillShape(gtx.Ops, col, clip.Stroke{Path: path.End(), Width: float32(gtx.Dp(2))}.Op())
	}

	dims := layout.Dimensions{Size: image.Pt(sz, sz)}

	// Rotate based on progress
	defer op.Affine(f32.Affine2D{}.Rotate(center, progress*3.14159)).Push(gtx.Ops).Pop()

	// Draw Sun (when progress is close to 0)
	if progress < 0.5 {
		// Sun body
		circle := clip.Ellipse{
			Min: image.Pt(sz/4, sz/4),
			Max: image.Pt(sz*3/4, sz*3/4),
		}.Op(gtx.Ops)
		paint.FillShape(gtx.Ops, col, circle)

		// Rays
		for i := 0; i < 8; i++ {
			angle := float64(i) * (2 * math.Pi / 8)
			start := f32.Pt(
				center.X+float32(sz)*0.3*float32(math.Cos(angle)),
				center.Y+float32(sz)*0.3*float32(math.Sin(angle)),
			)
			end := f32.Pt(
				center.X+float32(sz)*0.45*float32(math.Cos(angle)),
				center.Y+float32(sz)*0.45*float32(math.Sin(angle)),
			)
			drawLine(start, end)
		}
	} else {
		// Draw Moon (when progress is close to 1)
		// Simple crescent using two arcs
		path := clip.Path{}
		path.Begin(gtx.Ops)
		// Outer arc
		path.Arc(f32.Pt(center.X+float32(sz)*0.4, center.Y), f32.Pt(center.X-float32(sz)*0.4, center.Y), 3.14)
		// This is hard to get right without trial and error.
		// Let's use the "draw circle, then draw background circle" trick.
		// But we need the background color.
		// Since we don't have it passed, let's assume we can just draw the crescent shape.

		// Let's try a simpler shape: A 'C'
		path.MoveTo(f32.Pt(center.X+float32(sz)*0.2, center.Y-float32(sz)*0.3))
		path.QuadTo(
			f32.Pt(center.X-float32(sz)*0.4, center.Y),
			f32.Pt(center.X+float32(sz)*0.2, center.Y+float32(sz)*0.3),
		)
		path.QuadTo(
			f32.Pt(center.X-float32(sz)*0.1, center.Y),
			f32.Pt(center.X+float32(sz)*0.2, center.Y-float32(sz)*0.3),
		)
		path.Close()
		paint.FillShape(gtx.Ops, col, clip.Outline{Path: path.End()}.Op())
	}
	return dims
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
