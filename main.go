package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/webp"
)

var supportedImageFormats = map[string]bool{
	"gif":  true,
	"webp": true,
}

func loadImages(path string, imageContainer *fyne.Container, scroller *container.Scroll) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		fyne.LogError(fmt.Sprintf("Failed to read directory: %s", path), err)
		return
	}

	var objects []fyne.CanvasObject
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		fullPath := filepath.Join(path, dirEntry.Name())
		if !isSupportedImageFormat(fullPath) {
			continue
		}

		img, err := decodeImage(fullPath)
		if err != nil {
			fmt.Println("error decoding image:", err)
			continue
		}
		raster := canvas.NewImageFromImage(img)
		raster.FillMode = canvas.ImageFillOriginal
		objects = append(objects, raster)
	}

	imageContainer.Objects = objects
	imageContainer.Refresh()

	scroller.ScrollToTop()
}

func decodeImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func isSupportedImageFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	if len(ext) < 2 {
		return false
	}
	return supportedImageFormats[ext[1:]]
}

func main() {
	image.RegisterFormat("webp", "RIFF????WEBPVP8 ", webp.Decode, webp.DecodeConfig)

	a := app.NewWithID("io.fyne.emojiviewer")
	w := a.NewWindow("Emoji Viewer")
	w.Resize(fyne.NewSize(420, 420))

	imageContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)))
	scroller := container.NewScroll(imageContainer)

	buttons := []string{"Anime", "Dance", "Emoji", "Pepe", "React", "Tech"}
	buttonContainer := container.NewVBox()
	var imageButtons []*widget.Button

	updateButtonState := func(activeBtn *widget.Button) {
		for _, b := range imageButtons {
			if b == activeBtn {
				b.Importance = widget.HighImportance
			} else {
				b.Importance = widget.LowImportance
			}
			b.Refresh()
		}
	}

	for _, label := range buttons {
		currentLabel := label
		button := widget.NewButton(currentLabel, nil)
		button.Importance = widget.LowImportance
		buttonContainer.Add(button)
		imageButtons = append(imageButtons, button)
		currentButton := button
		currentButton.OnTapped = func() {
			loadImages(fmt.Sprintf("./emojis/%s/", currentLabel), imageContainer, scroller)
			updateButtonState(currentButton)
		}
	}

	split := container.NewHSplit(buttonContainer, scroller)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}
