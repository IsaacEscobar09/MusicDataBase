package model

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"
)

type ConfigurationFile struct {
    ConfigPath      string
    DefaultDBPath   string
    DefaultMusicDir string
}

// NewConfigurationFile es el constructor para ConfigurationFile
func NewConfigurationFile() *ConfigurationFile {
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error obteniendo el usuario actual:", err)
        return nil
    }

    configDir := filepath.Join(usr.HomeDir, ".config", "MusicDataBase")
    configFilePath := filepath.Join(configDir, "MusicConfig.conf")
    
    // Cambiar la ruta de la base de datos a $HOME/.local/share/DataBase
    dbDir := filepath.Join(usr.HomeDir, ".local", "share", "DataBase")
    defaultDBPath := filepath.Join(dbDir, "MusicDataBase.db")

    musicDir := filepath.Join(usr.HomeDir, "Música")
    if _, err := os.Stat(musicDir); os.IsNotExist(err) {
        musicDir = filepath.Join(usr.HomeDir, "Music")
    }

    return &ConfigurationFile{
        ConfigPath:      configFilePath,
        DefaultDBPath:   defaultDBPath,
        DefaultMusicDir: musicDir,
    }
}

// CreateDefaultConfig crea el archivo de configuración con las rutas por defecto
func (cf *ConfigurationFile) CreateDefaultConfig() error {
    // Crear directorio .config/MusicDataBase si no existe
    if err := os.MkdirAll(filepath.Dir(cf.ConfigPath), 0755); err != nil {
        return fmt.Errorf("error creando el directorio de configuración: %v", err)
    }

    // Verificar si el archivo de configuración ya existe
    if _, err := os.Stat(cf.ConfigPath); err == nil {
        fmt.Println("El archivo de configuración ya existe.")
        return nil
    }

    // Crear y escribir el archivo de configuración
    file, err := os.Create(cf.ConfigPath)
    if err != nil {
        return fmt.Errorf("error creando el archivo de configuración: %v", err)
    }
    defer file.Close()

    // Escribir las rutas por defecto (nueva ruta para la base de datos)
    if _, err := file.WriteString(fmt.Sprintf("DB_PATH=%s\nMUSIC_DIR=%s\n", cf.DefaultDBPath, cf.DefaultMusicDir)); err != nil {
        return fmt.Errorf("error escribiendo en el archivo de configuración: %v", err)
    }

    fmt.Println("Archivo de configuración creado con rutas por defecto.")
    return nil
}

