package model

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "os/user"
    "path/filepath"
)

type MusicDataBase struct {
    dbPath string
}

// NewMusicDataBase es el constructor para MusicDataBase
func NewMusicDataBase(dbPath string) *MusicDataBase {
    return &MusicDataBase{dbPath: dbPath}
}

// InitializeDatabase inicializa la base de datos en la ruta especificada
func (mdb *MusicDataBase) InitializeDatabase() error {
    usr, err := user.Current() // Obtener el usuario actual para el directorio $HOME
    if err != nil {
        return fmt.Errorf("error obteniendo el usuario actual: %s", err)
    }

    // Cambiar el directorio de la base de datos a $HOME/.local/share/DataBase
    dbDir := filepath.Join(usr.HomeDir, ".local", "share", "DataBase")
    mdb.dbPath = filepath.Join(dbDir, "MusicDataBase.db")

    // Verificar si el directorio para la base de datos existe, si no, crearlo
    if _, err := os.Stat(dbDir); os.IsNotExist(err) {
        err := os.MkdirAll(dbDir, os.ModePerm)
        if err != nil {
            return fmt.Errorf("error al crear el directorio de la base de datos: %s", err)
        }
        fmt.Printf("Directorio creado en: %s\n", dbDir)
    }

    // Verificar si la base de datos existe, si no, crearla
    if _, err := os.Stat(mdb.dbPath); os.IsNotExist(err) {
        fmt.Printf("La base de datos no existe, se creará una nueva en: %s\n", mdb.dbPath)
    }

    // Abrir la base de datos
    db, err := sql.Open("sqlite3", mdb.dbPath)
    if err != nil {
        return fmt.Errorf("error al abrir la base de datos: %s", err)
    }
    defer db.Close()

    // Verificar si el esquema está completo
    if !completeSchemaExists(db) {
        fmt.Println("El esquema no está completo o no existe, se procederá a crear/acompletar el esquema...")
        createSchema(db)
        fmt.Println("Esquema creado exitosamente.")
    } else {
        fmt.Println("El esquema ya está completo.")
    }

    return nil 
}

// Función para verificar si el esquema de la base de datos está completo
func completeSchemaExists(db *sql.DB) bool {
    requiredTables := []string{"types", "performers", "persons", "groups", "in_group", "albums", "rolas"}

    for _, table := range requiredTables {
        if !tableExists(db, table) {
            return false
        }
    }
    return true
}

// Función para verificar si una tabla específica existe en la base de datos
func tableExists(db *sql.DB, tableName string) bool {
    var name string
    query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
    err := db.QueryRow(query, tableName).Scan(&name)
    if err != nil {
        if err == sql.ErrNoRows {
            return false
        }
        log.Fatalf("Error al verificar la tabla %s: %s\n", tableName, err)
    }
    return name == tableName
}

// Función para crear el esquema de la base de datos
func createSchema(db *sql.DB) {
    schema := `
        CREATE TABLE types (
            id_type       INTEGER PRIMARY KEY,
            description   TEXT
        );
        INSERT INTO types VALUES(0,'Person');
        INSERT INTO types VALUES(1,'Group');
        INSERT INTO types VALUES(2,'Unknown');

        CREATE TABLE performers (
            id_performer  INTEGER PRIMARY KEY,
            id_type       INTEGER,
            name          TEXT,
            FOREIGN KEY   (id_type) REFERENCES types(id_type)
        );

        CREATE TABLE persons (
            id_person     INTEGER PRIMARY KEY,
            stage_name    TEXT,
            real_name     TEXT,
            birth_date    TEXT,
            death_date    TEXT
        );

        CREATE TABLE groups (
            id_group      INTEGER PRIMARY KEY,
            name          TEXT,
            start_date    TEXT,
            end_date      TEXT
        );

        CREATE TABLE in_group (
            id_person     INTEGER,
            id_group      INTEGER,
            PRIMARY KEY   (id_person, id_group),
            FOREIGN KEY   (id_person) REFERENCES persons(id_person),
            FOREIGN KEY   (id_group) REFERENCES groups(id_group)
        );

        CREATE TABLE albums (
            id_album      INTEGER PRIMARY KEY,
            path          TEXT,
            name          TEXT,
            year          INTEGER
        );

        CREATE TABLE rolas (
            id_rola       INTEGER PRIMARY KEY,
            id_performer  INTEGER,
            id_album      INTEGER,
            path          TEXT,
            title         TEXT,
            track         INTEGER,
            year          INTEGER,
            genre         TEXT,
            FOREIGN KEY   (id_performer) REFERENCES performers(id_performer),
            FOREIGN KEY   (id_album) REFERENCES albums(id_album)
        );
    `
    _, err := db.Exec(schema)
    if err != nil {
        log.Fatalf("Error al crear el esquema: %s\n", err)
    }
}

