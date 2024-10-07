package model

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/dhowden/tag"
)

type MP3Miner struct {
}

func findDatabaseFile(dbDir string) (string, error) {
    var dbFile string
    err := filepath.Walk(dbDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".db" {
            dbFile = path
            return filepath.SkipDir         }
        return nil
    })

    if dbFile == "" {
        return "", fmt.Errorf("no se encontró un archivo de base de datos en %s", dbDir)
    }

    return dbFile, err
}

func (m *MP3Miner) MineDirectory(path string, dbDir string) {
    dbPath, err := findDatabaseFile(dbDir)
    if err != nil {
        log.Fatalf("Error al encontrar el archivo de base de datos: %v\n", err)
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatalf("Error al abrir la base de datos: %v\n", err)
    }
    defer db.Close()

    err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(filePath) == ".mp3" {
            fmt.Printf("Analizando archivo: %s\n", filePath)
            ExtractMetadata(filePath, db)
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error al recorrer el directorio: %v\n", err)
    }
}

func ExtractMetadata(filePath string, db *sql.DB) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error al abrir el archivo: %s\n", err)
        return
    }
    defer file.Close()

    metadata, err := tag.ReadFrom(file)
    if err != nil {
        log.Printf("Error al leer los metadatos: %s\n", err)
        return
    }

    currentYear := time.Now().Year()

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
        album = filepath.Base(filepath.Dir(filePath)) 
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

    if songExists(db, filePath) {
        fmt.Printf("La canción ya existe en la base de datos, omitiendo: %s\n", filePath)
        return
    }

    insertAlbum(db, album, year, filepath.Dir(filePath))
    insertPerformer(db, artist)
    insertRola(db, artist, album, filePath, title, trackNum, year, genre)
}

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

func insertAlbum(db *sql.DB, name string, year int, path string) {
    _, err := db.Exec("INSERT OR IGNORE INTO albums (name, year, path) VALUES (?, ?, ?)", name, year, path)
    if err != nil {
        log.Printf("Error al insertar el álbum: %v\n", err)
    }
}

func insertPerformer(db *sql.DB, name string) {
    _, err := db.Exec("INSERT OR IGNORE INTO performers (name, id_type) VALUES (?, ?)", name, 2)
    if err != nil {
        log.Printf("Error al insertar el intérprete: %v\n", err)
    }
}

func insertRola(db *sql.DB, artist, album, filePath, title string, trackNum, year int, genre string) {
    var id_performer int
    var id_album int

    err := db.QueryRow("SELECT id_performer FROM performers WHERE name = ?", artist).Scan(&id_performer)
    if err != nil {
        log.Printf("Error al obtener el ID del intérprete: %v\n", err)
        return
    }

    err = db.QueryRow("SELECT id_album FROM albums WHERE name = ?", album).Scan(&id_album)
    if err != nil {
        log.Printf("Error al obtener el ID del álbum: %v\n", err)
        return
    }

    _, err = db.Exec("INSERT INTO rolas (id_performer, id_album, path, title, track, year, genre) VALUES (?, ?, ?, ?, ?, ?, ?)",
        id_performer, id_album, filePath, title, trackNum, year, genre)
    if err != nil {
        log.Printf("Error al insertar la rola: %v\n", err)
    }
}
