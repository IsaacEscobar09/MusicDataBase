package model

import (
    "database/sql"
    "fmt"
    "strings"

    _ "github.com/mattn/go-sqlite3"
)

// Compiler es una estructura vacía que contiene el método CompileSearch para compilar búsquedas en la base de datos.
type Compiler struct{}

// CompileSearch analiza el string de búsqueda y genera una consulta SQL para filtrar canciones en la base de datos.
// Recibe como parámetros el string de búsqueda, la conexión a la base de datos y devuelve un arreglo de canciones.
func (c *Compiler) CompileSearch(searchString string, db *sql.DB) ([]Song, error) { 
    // Divide la cadena de búsqueda en pares clave-valor, separados por comas.
    filters := strings.Split(searchString, ",")
    queryConditions := []string{}
    args := []interface{}{}

    // Variable para verificar si existen filtros específicos en la búsqueda.
    hasSpecificFilters := false

    // Recorremos los filtros para construir la condición de la consulta SQL.
    for _, filter := range filters {
        filter = strings.TrimSpace(filter) // Elimina espacios en blanco innecesarios.
        if filter == "" {
            continue
        }

        parts := strings.SplitN(filter, ":", 2) // Divide la clave y el valor en cada filtro.
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(strings.ToLower(parts[0])) // Clave del filtro (e.g., performer, album).
        value := strings.TrimSpace(parts[1]) // Valor del filtro (e.g., nombre del artista o álbum).

        if value == "" {
            continue // Ignora los filtros sin valores.
        }

        // Construcción de condiciones SQL según el tipo de filtro.
        switch key {
        case "p":
            queryConditions = append(queryConditions, "performers.name LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "a":
            queryConditions = append(queryConditions, "albums.name LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "c":
            queryConditions = append(queryConditions, "rolas.title LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "g":
            queryConditions = append(queryConditions, "rolas.genre LIKE ?")
            args = append(args, "%"+value+"%")
            hasSpecificFilters = true
        case "y":
            queryConditions = append(queryConditions, "rolas.year = ?")
            args = append(args, value)
            hasSpecificFilters = true
        }
    }

    // Si no se pasan filtros específicos, realizar una búsqueda general en los campos principales.
    if !hasSpecificFilters {
        generalCondition := "(rolas.title LIKE ? OR rolas.year LIKE ? OR rolas.genre LIKE ? OR rolas.track LIKE ? OR performers.name LIKE ? OR albums.name LIKE ?)"
        queryConditions = append(queryConditions, generalCondition)
        searchTerm := "%" + searchString + "%"
        args = append(args, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
    }

    // Construcción final de la consulta SQL.
    query := `
    SELECT rolas.id_rola, rolas.title, performers.name AS artist, albums.name AS album, rolas.year, rolas.genre, rolas.track
    FROM rolas
    JOIN performers ON rolas.id_performer = performers.id_performer
    JOIN albums ON rolas.id_album = albums.id_album
    WHERE ` + strings.Join(queryConditions, " AND ")

    // Ejecución de la consulta SQL.
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("Error ejecutando la consulta: %v", err)
    }
    defer rows.Close()

    // Recorremos los resultados y los agregamos al arreglo de canciones.
    var songs []Song
    for rows.Next() {
        var song Song
        err := rows.Scan(&song.IDRola, &song.Title, &song.Artist, &song.Album, &song.Year, &song.Genre, &song.Track)
        if err != nil {
            return nil, fmt.Errorf("Error leyendo los resultados: %v", err)
        }
        songs = append(songs, song)
    }

    // Si no se encontraron resultados, devolver un error.
    if len(songs) == 0 {
        return nil, fmt.Errorf("No se encontraron canciones que coincidan con los filtros")
    }

    return songs, nil
}

