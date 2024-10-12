package model

// Song representa una canción dentro de la base de datos de música.
// Contiene información como el título, el artista, el álbum, el año, el género y el número de pista.
type Song struct {
    IDRola  int    // ID único de la canción (correspondiente al campo id_rola en la base de datos)
    Title   string // Título de la canción
    Artist  string // Artista o intérprete de la canción
    Album   string // Nombre del álbum en el que aparece la canción
    Year    int    // Año de lanzamiento de la canción
    Genre   string // Género musical de la canción
    Track   int    // Número de pista en el álbum
}

