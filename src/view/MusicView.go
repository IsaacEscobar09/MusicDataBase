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

const maxCharLength = 30 // Cambia esto al número máximo de caracteres que desees mostrar

// Función para truncar texto
func truncateText(text string, maxLength int) string {
    if len(text) > maxLength {
        return text[:maxLength] + "..." // Agregar "..." si el texto se trunca
    }
    return text
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

    // Barra de progreso
    progressBar := widget.NewProgressBar()

    // Variables para los datos de la tabla
    var songData [][]string
    var songDataWithHeader [][]string

    // Crear una tabla para mostrar los resultados de las canciones
    songTable := widget.NewTable(
        func() (int, int) {
            return len(songDataWithHeader), 3 // Número de filas y 3 columnas
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            if id.Row == 0 {
                // Encabezados de columna
                label.SetText(songDataWithHeader[0][id.Col])
                label.TextStyle = fyne.TextStyle{Bold: true}
            } else {
                // Datos de las canciones, aplicando truncamiento
                label.SetText(truncateText(songDataWithHeader[id.Row][id.Col], maxCharLength))
                label.TextStyle = fyne.TextStyle{}
            }
        },
    )
    songTable.SetColumnWidth(0, 600) // Ancho de la columna Canción
    songTable.SetColumnWidth(1, 600) // Ancho de la columna Performer
    songTable.SetColumnWidth(2, 600) // Ancho de la columna Álbum

    // Función para cargar los datos de la tabla
    loadTableData := func() {
        data, err := mc.CreateSongTableData()
        if err != nil {
            dialog.ShowError(err, myWindow)
            return
        }
        songData = data

        // Añadir encabezados
        songDataWithHeader = [][]string{{"Canción", "Performer", "Álbum"}}
        songDataWithHeader = append(songDataWithHeader, songData...)

        songTable.Refresh()
    }

    // Cargar los datos al inicio
    loadTableData()

    // Componentes de la UI
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Buscar canción:     performer: <nombre>, album: <album>, cancion: <cancion>, genero: <genero>, año: <año>")

    // Función para realizar la búsqueda y actualizar la tabla
    performSearch := func() {
        searchString := searchEntry.Text
        data, err := mc.SearchSongsTableData(searchString)
        if err != nil {
            dialog.ShowError(err, myWindow)
            return
        }
        songData = data

        // Añadir encabezados
        songDataWithHeader = [][]string{{"Canción", "Performer", "Álbum"}}
        songDataWithHeader = append(songDataWithHeader, songData...)

        songTable.Refresh()
    }

    searchEntry.OnSubmitted = func(text string) {
        performSearch() // Ejecuta la búsqueda cuando el usuario presiona Enter
    }

    // Manejar clics en las filas de la tabla
    songTable.OnSelected = func(id widget.TableCellID) {
        if id.Row > 0 { // Ignorar el encabezado
            infoWindow := myApp.NewWindow("Información de la Canción")
            infoLabel := widget.NewLabel("Información de la canción en proceso")
            infoContainer := container.NewVBox(infoLabel)
            infoWindow.SetContent(infoContainer)
            infoWindow.Resize(fyne.NewSize(300, 200))
            infoWindow.Show()
        }
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
        loadTableData()
    })

    // Crear los botones en un contenedor horizontal
    buttonsContainer := container.NewHBox(
        helpButton,
        settingsButton,
        homeButton,
        layout.NewSpacer(),
        minimizeButton,
        fullscreenButton,
        closeButton,
    )

    // Botón para iniciar minería con barra de progreso
    minerButton := widget.NewButton("Miner", func() {
        // Llamar a la minería con barra de progreso y callback para refrescar las canciones
        mc.StartMiningWithProgress(myWindow, progressBar, func() {
            loadTableData()
        })
    })

    // Definir el contenido principal
    content := container.NewBorder(
        container.NewVBox(buttonsContainer, searchEntry, progressBar, minerButton), // Parte superior
        nil, // Parte inferior
        nil, // Parte izquierda
        nil, // Parte derecha
        songTable, // Área principal con la tabla de canciones
    )

    myWindow.SetContent(content)

    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}

