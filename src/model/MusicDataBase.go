package model

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "path/filepath"
)

func InitializeDatabase(dbPath string) {
    configDir := filepath.Dir(dbPath)
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        err := os.MkdirAll(configDir, os.ModePerm)
        if err != nil {
            log.Fatalf("Error al crear el directorio de configuración: %s\n", err)
        }
        fmt.Printf("Directorio creado en: %s\n", configDir)
    }

    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        fmt.Printf("El archivo no existe, se creará una nueva base de datos en: %s\n", dbPath)
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatalf("Error al abrir la base de datos: %s\n", err)
    }
    defer db.Close()

    if !completeSchemaExists(db) {
        fmt.Println("El esquema no está completo o no existe, se procederá a crear/acompletar el esquema...")
        createSchema(db)
        fmt.Println("Esquema creado exitosamente.")
    } else {
        fmt.Println("El esquema ya está completo.")
    }
}

func completeSchemaExists(db *sql.DB) bool {
    requiredTables := []string{"types", "performers", "persons", "groups", "in_group", "albums", "rolas"}

    for _, table := range requiredTables {
        if !tableExists(db, table) {
            return false
        }
    }
    return true
}

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
