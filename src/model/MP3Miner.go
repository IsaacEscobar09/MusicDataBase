package model

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/dhowden/tag"
)

// Método principal para minar un directorio
func MineDirectory(path string) {
    // Recorre el directorio y busca archivos .mp3
    err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // Si el archivo tiene la extensión .mp3
        if filepath.Ext(filePath) == ".mp3" {
            fmt.Printf("Analizando archivo: %s\n", filePath)
            ExtractMetadata(filePath)
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error al recorrer el directorio: %v\n", err)
    }
}

// Método para extraer las etiquetas de un archivo MP3
func ExtractMetadata(filePath string) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error al abrir el archivo: %s\n", err)
        return
    }
    defer file.Close()

    // Usando la biblioteca tag para obtener los metadatos
    metadata, err := tag.ReadFrom(file)
    if err != nil {
        log.Printf("Error al leer los metadatos: %s\n", err)
        return
    }

    // Obtener el año actual o de creación del archivo
    currentYear := time.Now().Year()

    // Imprimir etiquetas con valores por omisión
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
        album = "Unknown"
    }

    year := metadata.Year()
    if year == 0 {
        year = currentYear
    }

    genre := metadata.Genre()
    if genre == "" {
        genre = "Unknown"
    }

    trackNum, totalTracks := metadata.Track()
    if trackNum == 0 {
        trackNum = 1 // Valor por omisión si no hay número de pista
    }

    // Imprimir las etiquetas obtenidas
    fmt.Printf("Título: %s\n", title)
    fmt.Printf("Artista: %s\n", artist)
    fmt.Printf("Álbum: %s\n", album)
    fmt.Printf("Año: %d\n", year)
    fmt.Printf("Género: %s\n", genre)
    fmt.Printf("Pista: %d/%d\n", trackNum, totalTracks)
}
