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

    // Barra de progreso
    progressBar := widget.NewProgressBar()

    // Crear el contenedor desplazable para las canciones
    scrollableSongs := container.NewVScroll(container.NewVBox()) // Inicializar el contenedor de canciones vacío
    scrollableSongs.SetMinSize(fyne.NewSize(800, 400)) // Define el tamaño mínimo del área desplazable

    // Componentes de la UI
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Buscar canción")

    // Función para realizar la búsqueda y actualizar la lista
    performSearch := func() {
        searchString := searchEntry.Text // Capturamos el texto ingresado por el usuario
        songContainers := mc.CreateSearchResultsContainers(searchString) // Llamamos al método de búsqueda
        scrollableSongs.Content = container.NewVBox(songContainers...) // Actualizamos el contenido
        scrollableSongs.Refresh() // Refrescamos la vista para mostrar los nuevos resultados
    }

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

    // Botón "Inicio" para restablecer la vista de todas las canciones
    homeButton := widget.NewButton("Inicio", func() {
        songContainers := mc.CreateSongContainers()
        scrollableSongs.Content = container.NewVBox(songContainers...)
        scrollableSongs.Refresh()
    })

    // Inicialmente cargar todas las canciones de la base de datos
    songContainers := mc.CreateSongContainers()
    scrollableSongs.Content = container.NewVBox(songContainers...)

    buttonsContainer := container.NewHBox(
        helpButton,
        settingsButton,
        homeButton, // Agregamos el botón "Inicio" al contenedor de botones
        layout.NewSpacer(),
        minimizeButton,
        fullscreenButton,
        closeButton,
    )

    // Botón para iniciar minería con barra de progreso
    minerButton := widget.NewButton("Miner", func() {
        // Llamar a la minería con barra de progreso y callback para refrescar las canciones
        mc.StartMiningWithProgress(myWindow, progressBar, func() {
            songContainers := mc.CreateSongContainers()
            scrollableSongs.Content = container.NewVBox(songContainers...)
            scrollableSongs.Refresh()
        })
    })

    // Acción al presionar la tecla "Enter" en la barra de búsqueda
    searchEntry.OnSubmitted = func(text string) {
        performSearch() // Ejecuta la búsqueda cuando el usuario presiona Enter
    }

    // Definir el contenido principal, usando un spacer para empujar la lista de canciones
    myWindow.SetContent(container.NewBorder(
        container.NewVBox(buttonsContainer, searchEntry, progressBar, minerButton), // Parte superior
        nil,    // Parte inferior vacía
        nil,    // Izquierda vacía
        nil,    // Derecha vacía
        scrollableSongs,  // Centro: lista de canciones que ocupa el resto de la pantalla
    ))

    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}

