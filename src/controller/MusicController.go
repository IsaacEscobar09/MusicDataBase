package controller

import (
    "database/sql"
    "fmt"
    "os/exec"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    _ "github.com/mattn/go-sqlite3" 
    "github.com/IsaacEscobar09/MusicDataBase/src/model"
)

type MusicController struct {
    ConfigFile    *model.ConfigurationFile
    MP3Miner      *model.MP3Miner
    MusicDatabase *model.MusicDataBase
    DB            *sql.DB 
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

// Definición de la estructura Song dentro de MusicController.go
type Song struct {
    IDRola int
    Title  string
    Artist string
    Album  string
    Year   int
    Genre  string
    Track  int
}

// Método para obtener todas las canciones de la base de datos
func (mc *MusicController) GetAllSongs() ([]Song, error) {
    rows, err := mc.DB.Query("SELECT r.id_rola, r.title, p.name, a.name, r.year, r.genre, r.track FROM rolas r JOIN performers p ON r.id_performer = p.id_performer JOIN albums a ON r.id_album = a.id_album")
    if err != nil {
        return nil, fmt.Errorf("error al obtener las canciones: %v", err)
    }
    defer rows.Close()

    var songs []Song
    for rows.Next() {
        var song Song
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

// Crear contenedores para cada canción en la base de datos
func (mc *MusicController) CreateSongContainers() []fyne.CanvasObject {
    songs, err := mc.GetAllSongs()
    if err != nil {
        return []fyne.CanvasObject{widget.NewLabel("Error al obtener canciones")}
    }

    songContainers := []fyne.CanvasObject{}
    for _, song := range songs {
        songContainer := mc.createSongContainer(song)
        songContainers = append(songContainers, songContainer)
    }
    return songContainers
}

// Función privada para crear un contenedor de canción
func (mc *MusicController) createSongContainer(song Song) *fyne.Container {
    songInfo := widget.NewLabel(fmt.Sprintf("Canción: %s | Artista: %s | Álbum: %s | Año: %d | Género: %s | No. de pista: %d",
        song.Title, song.Artist, song.Album, song.Year, song.Genre, song.Track))

    container := container.NewVBox(songInfo)
    return container
}

// Iniciar minería de música en una gorutina
func (mc *MusicController) StartMiningWithProgress(parent fyne.Window, progressBar *widget.ProgressBar) {
    go func() {
        totalFiles := mc.MP3Miner.GetTotalFiles(mc.ConfigFile.DefaultMusicDir)  
        mc.MP3Miner.MineDirectoryWithProgress(mc.ConfigFile.DefaultMusicDir, mc.ConfigFile.DefaultDBPath, progressBar, totalFiles)
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

