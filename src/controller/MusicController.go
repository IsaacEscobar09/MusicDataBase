package controller

import (
    "database/sql"
    "fmt"
    "os/exec"

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

// Constructor de MusicController
func NewMusicController() *MusicController {
    configFile := model.NewConfigurationFile()
    mp3Miner := &model.MP3Miner{}
    musicDatabase := model.NewMusicDataBase(configFile.DefaultDBPath)

    db, err := sql.Open("sqlite3", configFile.DefaultDBPath)
    if err != nil {
        panic(fmt.Sprintf("Error al abrir la base de datos: %v", err))
    }

    return &MusicController{
        ConfigFile:    configFile,
        MP3Miner:      mp3Miner,
        MusicDatabase: musicDatabase,
        DB:            db,
    }
}

// Método para obtener todas las canciones de la base de datos
func (mc *MusicController) GetAllSongs() ([]model.Song, error) {
    rows, err := mc.DB.Query("SELECT r.id_rola, r.title, p.name, a.name, r.year, r.genre, r.track FROM rolas r JOIN performers p ON r.id_performer = p.id_performer JOIN albums a ON r.id_album = a.id_album")
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

// Método para crear los datos de la tabla con las canciones de la base de datos
func (mc *MusicController) CreateSongTableData() ([][]string, error) {
    songs, err := mc.GetAllSongs()
    if err != nil {
        return nil, fmt.Errorf("error al obtener canciones: %v", err)
    }

    // Crear una matriz de strings para representar las filas de la tabla
    songData := make([][]string, len(songs))
    for i, song := range songs {
        songData[i] = []string{song.Title, song.Artist, song.Album}
    }

    return songData, nil
}

// Método para buscar canciones y obtener los datos de la tabla
func (mc *MusicController) SearchSongsTableData(searchString string) ([][]string, error) {
    songs, err := mc.SearchSongs(searchString)
    if err != nil {
        return nil, fmt.Errorf("error al buscar canciones: %v", err)
    }

    if len(songs) == 0 {
        return [][]string{}, nil
    }

    // Crear una matriz de strings para representar las filas de la tabla
    songData := make([][]string, len(songs))
    for i, song := range songs {
        songData[i] = []string{song.Title, song.Artist, song.Album}
    }

    return songData, nil
}

func (mc *MusicController) SearchSongs(searchString string) ([]model.Song, error) {
    // Si el string de búsqueda está vacío, devuelve todas las canciones
    if searchString == "" {
        return mc.GetAllSongs()
    }

    // Si tienes un compilador que interpreta el searchString, úsalo aquí
    if mc.Compiler == nil {
        mc.Compiler = &model.Compiler{}
    }

    songs, err := mc.Compiler.CompileSearch(searchString, mc.DB)
    if err != nil {
        return nil, fmt.Errorf("error al buscar canciones: %v", err)
    }

    return songs, nil
}

// Iniciar minería de música en una gorutina
func (mc *MusicController) StartMiningWithProgress(parent fyne.Window, progressBar *widget.ProgressBar, onMiningComplete func()) {
    go func() {
        totalFiles := mc.MP3Miner.GetTotalFiles(mc.ConfigFile.DefaultMusicDir)
        mc.MP3Miner.MineDirectoryWithProgress(mc.ConfigFile.DefaultMusicDir, mc.ConfigFile.DefaultDBPath, progressBar, totalFiles)
        onMiningComplete()
        dialog.ShowInformation("Miner", "Minería completada.", parent)
    }()
}

// Verificar y crear configuración y base de datos si es necesario
func (mc *MusicController) CheckConfigAndDB() error {
    if err := mc.MusicDatabase.InitializeDatabase(); err != nil {
        return fmt.Errorf("error al inicializar la base de datos: %v", err)
    }

    if err := mc.ConfigFile.CreateDefaultConfig(); err != nil {
        return fmt.Errorf("error al crear el archivo de configuración: %v", err)
    }

    return nil
}

// Mostrar diálogo de configuración
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

// Actualizar ruta de directorio de música
func (mc *MusicController) UpdateMusicDirectory(newDir string) {
    mc.ConfigFile.DefaultMusicDir = newDir
    mc.ConfigFile.CreateDefaultConfig()
}

// Actualizar ruta de base de datos
func (mc *MusicController) UpdateDatabasePath(newDBPath string) {
    mc.ConfigFile.DefaultDBPath = newDBPath
    mc.ConfigFile.CreateDefaultConfig()
}

// Abrir ayuda en el navegador
func (mc *MusicController) OpenHelp() {
    exec.Command("xdg-open", "https://github.com/IsaacEscobar09/MusicDataBase").Start()
}

