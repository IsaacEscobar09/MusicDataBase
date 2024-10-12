# Music Data Base
## Autor: Escobar González Isaac Giovani

### Descripcion
Este repositorio es un gestor de Música en donde podras visualizar la música que tenga el usuario en una interfaz grafica (fyne) donde te proporcionara la información de cada cancion en dicha interfaz.

### Warning 
Es necesario que le usuario primero especifique la ruta de sus canciones, antes de minar.

### Requisitos
Es necesario que el usuario tenga instalado `Go` `(Golang)` en su maquina, exactamente la version `1.23.1`, para esto puedes descargarlo directamente de su pagina oficial la cual es:  
`https://go.dev/dl/`  
Ademas necesitara tener instalado `sqlite3`

### Compilacion
1. Para poder correr el programa debes de clonar este repositorio con el comando:  
`$ https://github.com/IsaacEscobar09/MusicDataBase.git`  
2. Una vez clonado te colocas dentro de la carpeta:  
`$ cd MusicDataBase`  
3. Despues puedes proceder a correr el programa de dos maneras ya sea con el comando:  
`$ go run src/main.go`  
O bien crear el ejecutable con el comando:  
`$ go build -o <NombreDelEjecutable> src/main.go`  
y para correrlo usar lo siguiente:  
`$ ./<NombreDelEjecutable>`
4. La primera vez que ejecutes el programa, la interfaz puede que llegue a tardar en aparecer o mostrarse ante el usuario pero tarde o temprano se mostrara, solo es la primera vez, ya después al ejecutarlo por segunda vez y en adelante, esta se mostrara rapido.  
5. Disfrutar el programa.

### Funciones de la Interfaz
La interfaz cuenta con los siguientes botones:
1. `Miner`: Este boton "minara" las canciones que tenga las canciones en el directorio que tenga elegido en "Settings" y al finalizar dicha operacion, le mostrara las canciones en la interfaz que mino.
2. `Inicio` : Este boton lo que hara es volver a poner TODAS las canciones que hayan sido minadas por el "minero".
3. `Setting`: Este boton te desplegara una ventana en la cual podras cambiar la ruta/path tanto de tu directorio en donde se encuentren tus canciones .mp3 (por defecto es Music o Musica si el sistema esta en idioma español) y tambien tu directorio de tu base de datos (por defecto es en $HOME/.local/share/DataBase).
4. `Help`: Este boton te mandara al repositorio de GitHub para encontrar más información.

### Barra de Busqueda
El usuario podra realizar busquedas con filtros o sin filtros y despues pulsando la tecla `Enter`.  
La busqueda con filtros es de la siguiente forma:  
`p: <performer>` para buscar por nombre de artista.  
`a: <album>` para buscar por albúm:  
`c: <canción>` para buscar por titulo de la canción.  
`g: <genero>` para buscar por genero de canción.  
`y: <año>` para buscar por año.

Puedes hacer uso de una `,` para poder buscar con más de un filtro.  
### Ejemplo (con filtros)  
`p: Beyonce, y: 2008`  
Este busqueda deberia darte todas las canciones de `Beyonce` del año `2008`.

Ahora bien, si quieres hacer una busqueda sin filtros, o bien, una busqueda general, simplemente busca mediante una palabra `<clave>` de dicha cancion, la interfaz se encargara de buscar todas las coincidencias en la base de datos para despues mostrarle las canciones correspondientes a dicha consulta.

### Ejemplo (sin filtros)
`Luis Miguel`  
Lo que hara la interfaz es mostrarle todas las coincidencias que encuentre en TODA la base de datos y tengan dicha palabra para mostrarselas al usuario.

### Extra 
Se pusieron los botones por default para `X`, `☐` y `−` para cerrar la aplicacion, pantalla completa y minimizar. Esto dado que si el usuario cuenta con un entorno de escritorio que no sea capaz de mostrarle dichos botones en la barra de la ventana, entonces estos botones le seran de utilidad (ademas de que hice el programa usando Hyprland y no podia visualizar dichos botones).
