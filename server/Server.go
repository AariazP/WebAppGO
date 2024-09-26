package main

import (
    "fmt"
    "html/template"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "time"
)

// Estructura para almacenar los datos que se pasan a la plantilla
type PageVariables struct {
    Hostname  string
    Categoria string
    Imagenes  []string
}

var hostname string

func main() {
    // Verificar si se ha proporcionado un puerto
    if len(os.Args) < 2 {
        fmt.Println("Por favor proporciona un número de puerto como argumento.")
        return
    }

    // Obtener el puerto desde los argumentos
    port := os.Args[1]

    // Obtener el nombre del host
    var err error
    hostname, err = os.Hostname()
    if err != nil {
        fmt.Println("Error al obtener el nombre del host:", err)
        return
    }

    // Iniciar el servidor
    IniciarServidor(port)
}

// Método que inicializa el servidor y lo pone a escuchar en el puerto especificado
func IniciarServidor(port string) {
    // Definir el manejador para la ruta raíz
    http.HandleFunc("/", handleRequest)

    // Iniciar el servidor en el puerto proporcionado
    fmt.Printf("Servidor escuchando en el puerto %s\n", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }
}

// Maneja las solicitudes a la ruta raíz
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Obtener la categoría aleatoria
    categoria := getRandomCategory()

    // Obtener las imágenes aleatorias de la categoría
    imagenes := getRandomImages(categoria)

    // Preparar los datos para la plantilla
    data := PageVariables{
        Hostname:  hostname,
        Categoria: categoria,
        Imagenes:  imagenes,
    }

    // Cargar la plantilla
    tmpl, err := template.ParseFiles("./static/index.html")
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
}

// getRandomCategory obtiene un nombre de directorio aleatorio de la ruta ./static
func getRandomCategory() string {
    // Leer los directorios en ./static
    entries, err := os.ReadDir("./static")
    if err != nil {
        fmt.Println("Error al leer el directorio:", err)
        return "Sin categoría"
    }

    // Filtrar solo los directorios
    var categories []string
    for _, entry := range entries {
        if entry.IsDir() {
            categories = append(categories, entry.Name())
        }
    }

    // Si no hay categorías, retornar "Sin categoría"
    if len(categories) == 0 {
        return "Sin categoría"
    }

    // Inicializar la semilla para el generador aleatorio
    rand.Seed(time.Now().UnixNano())
    // Seleccionar una categoría aleatoria
    return categories[rand.Intn(len(categories))]
}

// getRandomImages obtiene hasta 6 imágenes aleatorias de la categoría especificada
func getRandomImages(categoria string) []string {
    // Definir la ruta de la categoría
    categoryPath := filepath.Join("./static", categoria)

    // Leer los archivos en el directorio de la categoría
    entries, err := os.ReadDir(categoryPath)
    if err != nil {
        fmt.Println("Error al leer el directorio:", err)
        return []string{}
    }

    // Filtrar solo los archivos de imagen
    var images []string
    for _, entry := range entries {
        if !entry.IsDir() {
            // Puedes agregar validaciones para las extensiones de imagen si lo necesitas
            images = append(images, filepath.Join(categoria, entry.Name()))
        }
    }

    // Si no hay imágenes, retornar un arreglo vacío
    if len(images) == 0 {
        return []string{}
    }

    // Inicializar la semilla para el generador aleatorio
    rand.Seed(time.Now().UnixNano())

    // Asegurarse de que no se seleccionen más de 6 imágenes
    max := 6
    if len(images) < max {
        max = len(images)
    }

    // Barajar imágenes aleatoriamente
    rand.Shuffle(len(images), func(i, j int) {
        images[i], images[j] = images[j], images[i]
    })

    return images[:max]
}
