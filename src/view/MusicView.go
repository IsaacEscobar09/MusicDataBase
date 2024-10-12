package view

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
    "github.com/IsaacEscobar09/MusicDataBase/src/controller"
)

const maxCharLength = 55 // Máximo número de caracteres para mostrar en cada celda de la tabla

// truncateText es una función que trunca un texto si excede la longitud máxima permitida.
func truncateText(text string, maxLength int) string {
    if len(text) > maxLength {
        return text[:maxLength] + "..."
    }
    return text
}

// NewMusicView crea e inicializa la vista principal de la aplicación con una tabla de canciones, barra de búsqueda,
// botones de control y funcionalidad de minería de archivos MP3.
func NewMusicView() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Music Data Base")

    // Instanciar el controlador principal de la aplicación.
    mc := controller.NewMusicController()

    // Verificar la inicialización de archivos necesarios, como la configuración y la base de datos.
    if err := mc.CheckConfigAndDB(); err != nil {
        dialog.ShowError(err, myWindow)
        return
    }

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Se ha recuperado de un error inesperado:", r)
            // Puedes mostrar un diálogo de error aquí si lo deseas
        }
    }()

    // Crear barra de progreso para mostrar el avance de la minería.
    progressBar := widget.NewProgressBar()

    // Variables que almacenan los datos de las canciones que se mostrarán en la tabla.
    var songData [][]string
    var songDataWithHeader [][]string

    // Crear una tabla para mostrar los resultados de las canciones.
    songTable := widget.NewTable(
        func() (int, int) {
            return len(songDataWithHeader), 6 // Número de filas y columnas (Canción, Performer, Álbum).
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("") // Celda vacía inicial para la tabla.
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            if id.Row == 0 {
                label.SetText(songDataWithHeader[0][id.Col]) // Establecer encabezados.
                label.TextStyle = fyne.TextStyle{Bold: true}
            } else {
                label.SetText(truncateText(songDataWithHeader[id.Row][id.Col], maxCharLength)) // Mostrar datos truncados.
            }
        },
    )
    
    // Configurar el ancho de las columnas para ajustarse a los datos.
    songTable.SetColumnWidth(0, 500) // Ancho de la columna Canción.
    songTable.SetColumnWidth(1, 400) // Ancho de la columna Performer.
    songTable.SetColumnWidth(2, 400) // Ancho de la columna Álbum.
    songTable.SetColumnWidth(3, 100) // Ancho de la columna Año.
    songTable.SetColumnWidth(4, 200) // Ancho de la columna Genero.
    songTable.SetColumnWidth(5, 100) // Ancho de la columna No. de pista.
    

    // Función para cargar y actualizar los datos de la tabla de canciones desde el controlador.
    loadTableData := func() {
        data, err := mc.CreateSongTableData()
        if err != nil {
            dialog.ShowError(err, myWindow)
            return
        }
        songData = data
        songDataWithHeader = [][]string{{"Canción", "Performer", "Álbum", "Año", "Genero", "No. de pista"}}
        songDataWithHeader = append(songDataWithHeader, songData...)
        songTable.Refresh()
    }

    // Cargar los datos de la tabla al inicio de la aplicación.
    loadTableData()

    // Crear un campo de entrada para buscar canciones usando filtros por performer, álbum, etc.
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Buscar canción: 'p: <performer>, a: <album>, c: <cancion>, g: <genero>, y: <año>'")

    // Función para realizar la búsqueda y actualizar la tabla con los resultados.
    performSearch := func() {
        searchString := searchEntry.Text
        data, err := mc.SearchSongsTableData(searchString)
        if err != nil {
            dialog.ShowError(err, myWindow)
            return
        }
        songData = data
        songDataWithHeader = [][]string{{"Canción", "Performer", "Álbum", "Año", "Genero", "No. de pista"}}
        songDataWithHeader = append(songDataWithHeader, songData...)
        songTable.Refresh()
    }

    // Ejecutar la búsqueda cuando el usuario presiona Enter en el campo de búsqueda.
    searchEntry.OnSubmitted = func(text string) {
        performSearch()
    }

    

    // Crear botones de control para minimizar, pantalla completa y cerrar la aplicación.
    minimizeButton := widget.NewButton("−", func() {
        myWindow.Hide()
    })
    fullscreenButton := widget.NewButton("☐", func() {
        myWindow.SetFullScreen(!myWindow.FullScreen())
    })
    closeButton := widget.NewButton("X", func() {
        myApp.Quit()
    })

    // Crear botones adicionales como Help, Settings y un botón "Inicio" para restablecer la vista de todas las canciones.
    helpButton := widget.NewButton("Help", func() {
        mc.OpenHelp()
    })
    settingsButton := widget.NewButton("Settings", func() {
        mc.ShowSettingsDialog(myWindow)
    })
    homeButton := widget.NewButton("Inicio", func() {
        loadTableData()
    })

    // Agrupar los botones en un contenedor horizontal.
    buttonsContainer := container.NewHBox(
        helpButton,
        settingsButton,
        homeButton,
        layout.NewSpacer(),
        minimizeButton,
        fullscreenButton,
        closeButton,
    )

    // Botón "Miner" para iniciar el proceso de minería de archivos MP3, verificando archivos de configuración antes.
minerButton := widget.NewButton("Minero", func() {
    // Verificar los archivos de configuración y base de datos antes de comenzar la minería.
    if err := mc.CheckConfigAndDB(); err != nil {
        dialog.ShowError(err, myWindow)
        return
    }

    // Iniciar la minería con la barra de progreso y actualizar la tabla una vez finalizado.
    mc.StartMiningWithProgress(myWindow, progressBar, func() {
        loadTableData()
    })
})

    // Definir el contenido principal de la ventana, con la tabla de canciones, la barra de búsqueda y los botones de control.
    content := container.NewBorder(
        container.NewVBox(buttonsContainer, searchEntry, progressBar, minerButton), // Parte superior.
        nil, // Parte inferior.
        nil, // Parte izquierda.
        nil, // Parte derecha.
        songTable, // Área principal con la tabla de canciones.
    )

    // Configurar el tamaño inicial de la ventana y mostrar la aplicación.
    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}
