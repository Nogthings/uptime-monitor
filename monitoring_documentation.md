# Documentación del Paquete de Monitoreo (`monitoring`)

Este documento proporciona una explicación detallada del paquete `monitoring` en el proyecto Uptime Monitor.

## Descripción General

El paquete `monitoring` es responsable de la lógica central de monitoreo de servicios. Contiene los componentes para iniciar, detener y gestionar las tareas de monitoreo para diferentes servicios definidos en el sistema.

## Archivo: `internal/monitoring/monitor.go`

### Estructura `Monitor`

Esta es la estructura principal del paquete. Gestiona todos los workers de monitoreo activos.

```go
type Monitor struct {
    db      *pgxpool.Pool
    workers map[int64]chan struct{}
    mu      sync.Mutex
}
```

- `db (*pgxpool.Pool)`: Un pool de conexiones a la base de datos PostgreSQL, utilizado para interactuar con la base de datos y actualizar el estado de los servicios.
- `workers (map[int64]chan struct{})`: Un mapa que almacena los workers de monitoreo activos.
  - La clave (`int64`) es el ID del servicio que se está monitoreando.
  - El valor (`chan struct{}`) es un canal que se utiliza para enviar una señal de detención al worker.
- `mu (sync.Mutex)`: Un mutex para garantizar el acceso seguro y concurrente al mapa `workers`, evitando condiciones de carrera.

### Función `NewMonitor`

Esta función actúa como el constructor para la estructura `Monitor`.

```go
func NewMonitor(db *pgxpool.Pool) *Monitor
```

- **Propósito:** Crear e inicializar una nueva instancia de `Monitor`.
- **Parámetros:**
  - `db (*pgxpool.Pool)`: El pool de conexiones a la base de datos que utilizará el monitor.
- **Retorno:** Un puntero a la instancia de `Monitor` recién creada, con el mapa de `workers` inicializado.

### Función `StartServiceMonitor`

Esta función (aunque incompleta en el extracto) es la encargada de iniciar un nuevo worker de monitoreo para un servicio específico.

```go
func StartServiceMonitor()
```

- **Propósito:** Iniciar una goroutine (un worker) que monitoreará periódicamente un servicio.
- **Lógica esperada:**
  1.  Registrar un nuevo worker en el mapa `workers`.
  2.  Iniciar una goroutine que se ejecuta en un bucle infinito.
  3.  Dentro del bucle:
      - Realizar una comprobación del estado del servicio (por ejemplo, haciendo una petición HTTP).
      - Actualizar el estado del servicio en la base de datos.
      - Esperar un intervalo de tiempo definido antes de la siguiente comprobación.
      - Escuchar en el canal de detención para finalizar la goroutine si se recibe una señal.
