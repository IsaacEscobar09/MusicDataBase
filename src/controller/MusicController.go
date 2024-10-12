package controller

import (
    "database/sql"
    "fmt"
    "os/exec"
    "strconv" // Importado para convertir int a string

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    _ "github.com/mattn/go-sqlite3"
    "github.com/IsaacEscobar09/MusicDataBase/src/model"
)

type MusicController struct {
    ConfigFile    *model.ConfigurationFile
    MP3Miner      *model.MP3Miner
    MusicDatabase *model.MusicDataBase
    DB            *sql.DB
    Compiler      *model.Compiler
}

// NewMusicController crea una nueva instancia de MusicController.
// Inicializa los modelos de archivo de configuración, MP3Miner, base de datos y conexión a la base de datos SQLite.
func NewMusicController() *MusicController {
    configFile := model.NewConfigurationFile()
    mp3Miner := &model.MP3Miner{}
    musicDatabase := model.NewMusicDataBase(configFile.DefaultDBPath)

    db, err := sql.Open("sqlite3", configFile.DefaultDBPath)
    if err != nil {
    dialog.ShowError(fmt.Errorf("Error al abrir la base de datos: %v", err), nil)
    return nil
    }

    return &MusicController{
        ConfigFile:    configFile,
        MP3Miner:      mp3Miner,
        MusicDatabase: musicDatabase,
        DB:            db,
    }
}

// GetAllSongs devuelve todas las canciones almacenadas en la base de datos.
// Hace una consulta SQL que une las tablas de canciones, intérpretes y álbumes, devolviendo un arreglo de canciones.
func (mc *MusicController) GetAllSongs() ([]model.Song, error) {
    rows, err := mc.DB.Query(
        "SELECT r.id_rola, r.title, p.name, a.name, r.year, r.genre, r.track " +
            "FROM rolas r " +
            "JOIN performers p ON r.id_performer = p.id_performer " +
            "JOIN albums a ON r.id_album = a.id_album")
    if err != nil {
        return nil, fmt.Errorf("error al obtener las canciones: %v", err)
    }
    defer rows.Close()

    var songs []model.Song
    for rows.Next() {
        var song model.Song
        err := rows.Scan(&song.IDRola, &song.Title, &song.Artist, &song.Album, &song.Year, &song.Genre, &song.Track)
        if err != nil {
            return nil, fmt.Errorf("error al escanear canción: %v", err)
        }
        songs = append(songs, song)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error en la iteración de canciones: %v", err)
    }

    return songs, nil
}

// CreateSongTableData crea los datos de la tabla a partir de las canciones en la base de datos.
// Devuelve una matriz bidimensional de strings que representan las filas y columnas de la tabla.
func (mc *MusicController) CreateSongTableData() ([][]string, error) {
    songs, err := mc.GetAllSongs()
    if err != nil {
        return nil, fmt.Errorf("error al obtener canciones: %v", err)
    }

    songData := make([][]string, len(songs))
    for i, song := range songs {
        songData[i] = []string{
            song.Title,
            song.Artist,
            song.Album,
            strconv.Itoa(song.Year),   // Convertir int a string
            song.Genre,
            strconv.Itoa(song.Track),  // Convertir int a string
        }
    }

    return songData, nil
}

// SearchSongsTableData busca canciones en la base de datos con base en el string de búsqueda proporcionado.
// Devuelve los datos de la tabla filtrados de acuerdo con los resultados de la búsqueda.
func (mc *MusicController) SearchSongsTableData(searchString string) ([][]string, error) {
    songs, err := mc.SearchSongs(searchString)
    if err != nil {
        return nil, fmt.Errorf("error al buscar canciones: %v", err)
    }

    if len(songs) == 0 {
        return [][]string{}, nil
    }

    songData := make([][]string, len(songs))
    for i, song := range songs {
        songData[i] = []string{
            song.Title,
            song.Artist,
            song.Album,
            strconv.Itoa(song.Year),   // Convertir int a string
            song.Genre,
            strconv.Itoa(song.Track),  // Convertir int a string
        }
    }

    return songData, nil
}

// SearchSongs busca canciones utilizando un compilador de consultas y devuelve los resultados como un arreglo de canciones.
// Si el string de búsqueda está vacío, se devuelve el conjunto completo de canciones.
func (mc *MusicController) SearchSongs(searchString string) ([]model.Song, error) {
    if searchString == "" {
        return mc.GetAllSongs()
    }

    if mc.Compiler == nil {
        mc.Compiler = &model.Compiler{}
    }

    songs, err := mc.Compiler.CompileSearch(searchString, mc.DB)
    if err != nil {
        return nil, fmt.Errorf("error al buscar canciones: %v", err)
    }

    return songs, nil
}

// StartMiningWithProgress inicia el proceso de minería de archivos MP3 en una gorutina.
// Muestra una barra de progreso mientras se realiza la minería y refresca los datos de la tabla al completarse.
func (mc *MusicController) StartMiningWithProgress(parent fyne.Window, progressBar *widget.ProgressBar, onMiningComplete func()) {
    go func() {
        defer func() {
        if r := recover(); r != nil {
            dialog.ShowError(fmt.Errorf("Error inesperado durante la minería: %v", r), parent)
        }
    }()
        totalFiles := mc.MP3Miner.GetTotalFiles(mc.ConfigFile.DefaultMusicDir)
        mc.MP3Miner.MineDirectoryWithProgress(mc.ConfigFile.DefaultMusicDir, mc.ConfigFile.DefaultDBPath, progressBar, totalFiles)
        onMiningComplete()
        dialog.ShowInformation("Miner", "Minería completada.", parent)
    }()
}

// CheckConfigAndDB verifica si existen la base de datos y el archivo de configuración.
// Si no existen, los crea utilizando los métodos apropiados de los modelos.
func (mc *MusicController) CheckConfigAndDB() error {
    if err := mc.MusicDatabase.InitializeDatabase(); err != nil {
        return fmt.Errorf("error al inicializar la base de datos: %v", err)
    }

    if err := mc.ConfigFile.CreateDefaultConfig(); err != nil {
        return fmt.Errorf("error al crear el archivo de configuración: %v", err)
    }

    return nil
}

// ShowSettingsDialog muestra un diálogo para actualizar las rutas de la música y la base de datos.
// Permite al usuario cambiar estas configuraciones y las guarda en el archivo de configuración.
func (mc *MusicController) ShowSettingsDialog(parent fyne.Window) {
    musicDirEntry := widget.NewEntry()
    musicDirEntry.SetText(mc.ConfigFile.DefaultMusicDir)

    dbPathEntry := widget.NewEntry()
    dbPathEntry.SetText(mc.ConfigFile.DefaultDBPath)

    dialog.ShowForm("Settings", "Guardar", "Cancelar", []*widget.FormItem{
        {Text: "Ruta de Música", Widget: musicDirEntry},
        {Text: "Ruta de Base de Datos", Widget: dbPathEntry},
    }, func(response bool) {
        if response {
            mc.UpdateMusicDirectory(musicDirEntry.Text)
            mc.UpdateDatabasePath(dbPathEntry.Text)
            dialog.ShowInformation("Configuración", "Rutas actualizadas con éxito.", parent)
        }
    }, parent)
}

// UpdateMusicDirectory actualiza la ruta del directorio de música y guarda el cambio en el archivo de configuración.
func (mc *MusicController) UpdateMusicDirectory(newDir string) {
    mc.ConfigFile.DefaultMusicDir = newDir
    mc.ConfigFile.CreateDefaultConfig()
}

// UpdateDatabasePath actualiza la ruta de la base de datos y guarda el cambio en el archivo de configuración.
func (mc *MusicController) UpdateDatabasePath(newDBPath string) {
    mc.ConfigFile.DefaultDBPath = newDBPath
    mc.ConfigFile.CreateDefaultConfig()
}

// OpenHelp abre el navegador del sistema en la URL de ayuda del proyecto.
func (mc *MusicController) OpenHelp() {
    err := exec.Command("xdg-open", "https://github.com/IsaacEscobar09/MusicDataBase").Start()
    if err != nil {
        dialog.ShowError(fmt.Errorf("No se pudo abrir el navegador: %v", err), nil)
    }
}
