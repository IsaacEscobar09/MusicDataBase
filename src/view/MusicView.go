package view

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
    "github.com/IsaacEscobar09/MusicDataBase/src/controller"
)

type MusicView struct {
    controller *controller.MusicController
}

func NewMusicView() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Music Data Base")

    // Instanciar el controlador
    mc := controller.NewMusicController()

    // Verificar la inicialización de archivos necesarios
    if err := mc.CheckConfigAndDB(); err != nil {
        dialog.ShowError(err, myWindow)
        return
    }

    // Componentes de la UI
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Buscar canción")

    // Botones de la UI
    minimizeButton := widget.NewButton("−", func() {
        myWindow.Hide()
    })
    fullscreenButton := widget.NewButton("☐", func() {
        myWindow.SetFullScreen(!myWindow.FullScreen())
    })
    closeButton := widget.NewButton("X", func() {
        myApp.Quit()
    })
    helpButton := widget.NewButton("Help", func() {
        mc.OpenHelp()
    })
    settingsButton := widget.NewButton("Settings", func() {
        mc.ShowSettingsDialog(myWindow)
    })
    minerButton := widget.NewButton("Miner", func() {
        mc.StartMining(myWindow)
    })

    buttonsContainer := container.NewHBox(
        helpButton,
        settingsButton,
        layout.NewSpacer(),
        minimizeButton,
        fullscreenButton,
        closeButton,
    )

    // Obtener contenedores de canciones y agregarlos a la vista
    songContainers := mc.CreateSongContainers()

    // Crear un contenedor desplazable con las canciones
    scrollableSongs := container.NewVScroll(
        container.NewVBox(songContainers...),
    )
    scrollableSongs.SetMinSize(fyne.NewSize(800, 400)) // Define el tamaño mínimo del área desplazable

    // Definir el contenido principal, usando un spacer para empujar la lista de canciones
    myWindow.SetContent(container.NewBorder(
        container.NewVBox(buttonsContainer, searchEntry, minerButton), // Parte superior
        nil,    // Parte inferior vacía
        nil,    // Izquierda vacía
        nil,    // Derecha vacía
        scrollableSongs,  // Centro: lista de canciones que ocupa el resto de la pantalla
    ))

    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}

