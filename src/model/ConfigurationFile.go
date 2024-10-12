package model

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"
)

// ConfigurationFile define las rutas de configuración, la base de datos por defecto y el directorio de música por defecto.
type ConfigurationFile struct {
    ConfigPath      string // Ruta del archivo de configuración.
    DefaultDBPath   string // Ruta por defecto de la base de datos.
    DefaultMusicDir string // Ruta por defecto del directorio de música.
}

// NewConfigurationFile es el constructor para ConfigurationFile. Establece las rutas por defecto de configuración y base de datos.
func NewConfigurationFile() *ConfigurationFile {
    usr, err := user.Current() // Obtiene el usuario actual del sistema.
    if err != nil {
        fmt.Println("Error obteniendo el usuario actual:", err)
        return nil
    }

    // Ruta del archivo de configuración en $HOME/.config/MusicDataBase
    configDir := filepath.Join(usr.HomeDir, ".config", "MusicDataBase")
    configFilePath := filepath.Join(configDir, "MusicConfig.conf")
    
    // Cambia la ruta de la base de datos a $HOME/.local/share/DataBase
    dbDir := filepath.Join(usr.HomeDir, ".local", "share", "DataBase")
    defaultDBPath := filepath.Join(dbDir, "MusicDataBase.db")

    // Establece el directorio de música. Si no existe "Música", busca "Music".
    musicDir := filepath.Join(usr.HomeDir, "Música")
    if _, err := os.Stat(musicDir); os.IsNotExist(err) {
        musicDir = filepath.Join(usr.HomeDir, "Music")
    }

    // Retorna una nueva instancia de ConfigurationFile con las rutas configuradas.
    return &ConfigurationFile{
        ConfigPath:      configFilePath,
        DefaultDBPath:   defaultDBPath,
        DefaultMusicDir: musicDir,
    }
}

// CreateDefaultConfig crea un archivo de configuración con las rutas por defecto de la base de datos y música.
func (cf *ConfigurationFile) CreateDefaultConfig() error {
    // Crear el directorio .config/MusicDataBase si no existe.
    if err := os.MkdirAll(filepath.Dir(cf.ConfigPath), 0755); err != nil {
        return fmt.Errorf("error creando el directorio de configuración: %v", err)
    }

    // Verifica si el archivo de configuración ya existe.
    if _, err := os.Stat(cf.ConfigPath); err == nil {
        fmt.Println("El archivo de configuración ya existe.")
        return nil
    }

    // Crea el archivo de configuración.
    file, err := os.Create(cf.ConfigPath)
    if err != nil {
        return fmt.Errorf("error creando el archivo de configuración: %v", err)
    }
    defer file.Close() // Cierra el archivo al terminar.

    // Escribe las rutas por defecto en el archivo de configuración.
    if _, err := file.WriteString(fmt.Sprintf("DB_PATH=%s\nMUSIC_DIR=%s\n", cf.DefaultDBPath, cf.DefaultMusicDir)); err != nil {
        return fmt.Errorf("error escribiendo en el archivo de configuración: %v", err)
    }

    fmt.Println("Archivo de configuración creado con rutas por defecto.")
    return nil
}

