package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "text/template"
    "time"
)

type PageVariables struct {
    Hostname string
    Categoria string
    Imagenes []string
}

func main() {
    // Verificar si se ha proporcionado un puerto
    if len(os.Args) < 2 {
        fmt.Println("Por favor proporciona un número de puerto como argumento.")
        return
    }

    port := os.Args[1]
    IniciarServidor(port)
}

func IniciarServidor(port string) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Obtener el nombre del host
        hostname, err := os.Hostname()
        if err != nil {
            fmt.Println("Error al obtener el nombre del host:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Obtener una categoría aleatoria
        categorias := []string{"futbol", "artist", "landscape", "people"}
        rand.Seed(time.Now().UnixNano())
        categoria := categorias[rand.Intn(len(categorias))]

        // Obtener imágenes aleatorias de la categoría
        imagenes := obtenerImagenesAleatorias(filepath.Join("static", categoria), 6)
		fmt.Println(imagenes)
        // Preparar los datos para la plantilla
        data := PageVariables{
            Hostname: hostname,
            Categoria: categoria,
            Imagenes: imagenes,
        }

        // Cargar la plantilla
        tmpl, err := template.ParseFiles("static/index.html")
        if err != nil {
            fmt.Println("Error al cargar la plantilla:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Ejecutar la plantilla con los datos
        err = tmpl.Execute(w, data)
        if err != nil {
            fmt.Println("Error al ejecutar la plantilla:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    })

    // Iniciar el servidor en el puerto proporcionado
    fmt.Printf("Servidor escuchando en el puerto %s\n", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }
}

func obtenerImagenesAleatorias(directorio string, cantidad int) []string {
    files, err := os.ReadDir(directorio)
    if err != nil {
        fmt.Println("Error al leer el directorio:", err)
        return nil
    }

    var imagenes []string
    for _, file := range files {
        if !file.IsDir() {
            imagenes = append(imagenes, file.Name())
        }
    }

    // Mezclar y seleccionar aleatoriamente las imágenes
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(imagenes), func(i, j int) {
        imagenes[i], imagenes[j] = imagenes[j], imagenes[i]
    })

    // Asegurarse de no superar la cantidad solicitada
    if len(imagenes) > cantidad {
        imagenes = imagenes[:cantidad]
    }

    return imagenes
}