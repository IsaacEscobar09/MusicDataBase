package model

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // Importa el driver SQLite
    "github.com/dhowden/tag"        // Para leer metadatos ID3 de archivos MP3
    "fyne.io/fyne/v2/widget"        // Para manejar la barra de progreso
)

// MP3Miner es responsable de extraer metadatos de archivos MP3 y almacenarlos en la base de datos.
type MP3Miner struct {
    FileCount int // Contador de archivos MP3 procesados.
}

// findDatabaseFile busca el archivo de base de datos en un directorio especificado.
func findDatabaseFile(dbDir string) (string, error) {
    var dbFile string
    // Recorre el directorio y busca un archivo con extensión .db
    err := filepath.Walk(dbDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".db" {
            dbFile = path
            return filepath.SkipDir // Detiene la búsqueda una vez que se encuentra el archivo.
        }
        return nil
    })

    if dbFile == "" {
        return "", fmt.Errorf("no se encontró un archivo de base de datos en %s", dbDir)
    }

    return dbFile, err
}

// GetTotalFiles cuenta el número de archivos MP3 en un directorio y sus subdirectorios.
func (m *MP3Miner) GetTotalFiles(path string) int {
    fileCount := 0
    // Recorre el directorio y cuenta archivos con extensión .mp3
    filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err == nil && filepath.Ext(filePath) == ".mp3" {
            fileCount++
        }
        return nil
    })
    return fileCount
}

// MineDirectoryWithProgress procesa los archivos MP3 en un directorio, actualizando una barra de progreso.
func (m *MP3Miner) MineDirectoryWithProgress(path string, dbDir string, progressBar *widget.ProgressBar, totalFiles int) {
    // Encuentra el archivo de base de datos en el directorio especificado.
    dbPath, err := findDatabaseFile(dbDir)
    if err != nil {
        log.Fatalf("Error al encontrar el archivo de base de datos: %v\n", err)
    }

    // Abre la conexión a la base de datos.
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatalf("Error al abrir la base de datos: %v\n", err)
    }
    defer db.Close()

    currentFile := 0
    // Recorre el directorio y procesa cada archivo MP3.
    err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(filePath) == ".mp3" {
            fmt.Printf("Analizando archivo: %s\n", filePath)
            ExtractMetadata(filePath, db)

            // Actualiza la barra de progreso.
            currentFile++
            progressBar.SetValue(float64(currentFile) / float64(totalFiles)) // Actualiza progreso
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error al recorrer el directorio: %v\n", err)
    }
}

// ExtractMetadata extrae metadatos de un archivo MP3 y los inserta en la base de datos.
func ExtractMetadata(filePath string, db *sql.DB) {
    // Abre el archivo MP3.
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error al abrir el archivo: %s\n", err)
        return
    }
    defer file.Close()

    // Lee los metadatos ID3 del archivo.
    metadata, err := tag.ReadFrom(file)
    if err != nil {
        log.Printf("Error al leer los metadatos: %s\n", err)
        return
    }

    currentYear := time.Now().Year()

    // Obtiene los metadatos o asigna valores por defecto si faltan.
    title := metadata.Title()
    if title == "" {
        title = "Unknown"
    }

    artist := metadata.Artist()
    if artist == "" {
        artist = "Unknown"
    }

    album := metadata.Album()
    if album == "" {
        album = filepath.Base(filepath.Dir(filePath)) // Usa el nombre del directorio como álbum.
    }

    year := metadata.Year()
    if year == 0 {
        year = currentYear
    }

    genre := metadata.Genre()
    if genre == "" {
        genre = "Unknown"
    }

    trackNum, _ := metadata.Track()
    if trackNum == 0 {
        trackNum = 1
    }

    // Verifica si la canción ya existe en la base de datos.
    if songExists(db, filePath) {
        fmt.Printf("La canción ya existe en la base de datos, omitiendo: %s\n", filePath)
        return
    }

    // Inserta los datos en la base de datos.
    insertAlbum(db, album, year, filepath.Dir(filePath))
    insertPerformer(db, artist)
    insertRola(db, artist, album, filePath, title, trackNum, year, genre)
}

// songExists verifica si una canción ya existe en la base de datos.
func songExists(db *sql.DB, filePath string) bool {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM rolas WHERE path = ? LIMIT 1);`
    err := db.QueryRow(query, filePath).Scan(&exists)
    if err != nil {
        log.Printf("Error al verificar si la canción existe: %v\n", err)
        return false
    }
    return exists
}

// insertAlbum inserta un álbum en la base de datos si no existe.
func insertAlbum(db *sql.DB, name string, year int, path string) {
    _, err := db.Exec("INSERT OR IGNORE INTO albums (name, year, path) VALUES (?, ?, ?)", name, year, path)
    if err != nil {
        log.Printf("Error al insertar el álbum: %v\n", err)
    }
}

// insertPerformer inserta un intérprete en la base de datos si no existe.
func insertPerformer(db *sql.DB, name string) {
    _, err := db.Exec("INSERT OR IGNORE INTO performers (name, id_type) VALUES (?, ?)", name, 2)
    if err != nil {
        log.Printf("Error al insertar el intérprete: %v\n", err)
    }
}

// insertRola inserta una canción en la base de datos, asociándola con su intérprete y álbum.
func insertRola(db *sql.DB, artist, album, filePath, title string, trackNum, year int, genre string) {
    var id_performer int
    var id_album int

    // Obtiene el ID del intérprete.
    err := db.QueryRow("SELECT id_performer FROM performers WHERE name = ?", artist).Scan(&id_performer)
    if err != nil {
        log.Printf("Error al obtener el ID del intérprete: %v\n", err)
        return
    }

    // Obtiene el ID del álbum.
    err = db.QueryRow("SELECT id_album FROM albums WHERE name = ?", album).Scan(&id_album)
    if err != nil {
        log.Printf("Error al obtener el ID del álbum: %v\n", err)
        return
    }

    // Inserta la canción (rola) en la base de datos.
    _, err = db.Exec("INSERT INTO rolas (id_performer, id_album, path, title, track, year, genre) VALUES (?, ?, ?, ?, ?, ?, ?)",
        id_performer, id_album, filePath, title, trackNum, year, genre)
    if err != nil {
        log.Printf("Error al insertar la rola: %v\n", err)
    }
}

