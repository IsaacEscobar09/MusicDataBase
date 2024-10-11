package model

import (
    "database/sql"
    "fmt"
    "strings"

    _ "github.com/mattn/go-sqlite3"
)

type Compiler struct{}

// CompileSearch analiza el string de búsqueda y devuelve un arreglo de canciones que coinciden con los filtros
func (c *Compiler) CompileSearch(searchString string, db *sql.DB) ([]Song, error) { // Usa la estructura Song de Song.go
    // Dividimos el string en pares clave-valor separados por coma
    filters := strings.Split(searchString, ",")
    queryConditions := []string{}
    args := []interface{}{}

    // Variable para verificar si hay filtros específicos
    hasSpecificFilters := false

    // Recorremos los filtros para construir la consulta SQL
    for _, filter := range filters {
        filter = strings.TrimSpace(filter)
        if filter == "" {
            continue
        }

        parts := strings.SplitN(filter, ":", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(strings.ToLower(parts[0]))
        value := strings.TrimSpace(parts[1])

        if value == "" {
            continue // Si no hay valor después de ":", lo ignoramos
        }

        // Dependiendo del tipo de filtro, construimos diferentes condiciones
        switch key {
        case "performer":
            queryConditions = append(queryConditions, "performers.name LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "album":
            queryConditions = append(queryConditions, "albums.name LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "cancion":
            queryConditions = append(queryConditions, "rolas.title LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "genero":
            queryConditions = append(queryConditions, "rolas.genre LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "año":
            queryConditions = append(queryConditions, "rolas.year = ?")
            args = append(args, value)
            hasSpecificFilters = true
        }
    }

    // Si no se pasaron filtros específicos, realizar una búsqueda general
    if !hasSpecificFilters {
        // Buscamos en los campos principales de la base de datos (canción, álbum, artista)
        generalCondition := "(rolas.title LIKE ? OR performers.name LIKE ? OR albums.name LIKE ? OR rolas.genre LIKE ?)"
        queryConditions = append(queryConditions, generalCondition)
        searchTerm := "%" + searchString + "%"
        args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
    }

    // Construimos la consulta SQL
    query := `
    SELECT rolas.id_rola, rolas.title, performers.name AS artist, albums.name AS album, rolas.year, rolas.genre, rolas.track
    FROM rolas
    JOIN performers ON rolas.id_performer = performers.id_performer
    JOIN albums ON rolas.id_album = albums.id_album
    WHERE ` + strings.Join(queryConditions, " AND ")

    // Ejecutamos la consulta
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("Error ejecutando la consulta: %v", err)
    }
    defer rows.Close()

    // Recorremos los resultados y los agregamos al arreglo de canciones
    var songs []Song // Usa la estructura Song definida en Song.go
    for rows.Next() {
        var song Song // Usa la estructura Song
        err := rows.Scan(&song.IDRola, &song.Title, &song.Artist, &song.Album, &song.Year, &song.Genre, &song.Track)
        if err != nil {
            return nil, fmt.Errorf("Error leyendo los resultados: %v", err)
        }
        songs = append(songs, song)
    }

    // Si no se encontraron canciones, devolver un mensaje
    if len(songs) == 0 {
        return nil, fmt.Errorf("No se encontraron canciones que coincidan con los filtros")
    }

    return songs, nil
}

