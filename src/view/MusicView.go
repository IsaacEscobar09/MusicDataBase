package view

import (
    "os/exec" 
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
)

type MusicView struct {
}

func NewMusicView() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Music Data Base")

    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Buscar canción") // Placeholder para la barra de búsqueda

    minimizeButton := widget.NewButton("−", func() {
        myWindow.Hide() 
    })

    fullscreenButton := widget.NewButton("☐", func() {
        myWindow.SetFullScreen(!myWindow.FullScreen())     })

    closeButton := widget.NewButton("X", func() {
        myApp.Quit() 
    })

    helpButton := widget.NewButton("Help", func() {
        exec.Command("xdg-open", "https://github.com/IsaacEscobar09/MusicDataBase").Start()
    })

    settingsButton := widget.NewButton("Settings", func() {
    })

    buttonsContainer := container.NewHBox(
        helpButton,
        settingsButton,
        layout.NewSpacer(), 
        minimizeButton,
        fullscreenButton,
        closeButton,
    )

    myWindow.SetContent(container.NewVBox(
        buttonsContainer, 
        container.NewMax(searchEntry), 
    ))

    myWindow.Resize(fyne.NewSize(800, 600)) 
    myWindow.ShowAndRun()
}

