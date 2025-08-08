# Raylib Adapter

Questo adapter fornisce un'implementazione dell'interfaccia `RenderAdapter` utilizzando [Raylib](https://www.raylib.com/), una libreria semplice e facile da usare per lo sviluppo di giochi e applicazioni grafiche.

## Caratteristiche

- **Rendering 3D**: Supporto completo per il rendering 3D con controlli camera manuali
- **Rendering di corpi**: Rendering automatico di sfere e cubi basato sul tipo di materiale
- **Vettori di debug**: Visualizzazione di vettori di velocità e accelerazione
- **Bounding boxes**: Rendering di AABB per il debug
- **Controlli camera**: Controlli camera manuali con mouse e tastiera
- **Colori personalizzati**: Mappatura automatica dei colori basata sui materiali
- **UI integrata**: Istruzioni e stato dei controlli visualizzati a schermo

## Controlli

### Mouse
- **Click destro**: Attiva/disattiva la cattura del mouse
- **Movimento mouse** (quando catturato): Ruota la vista
- **Scroll**: Zoom in/out

### Tastiera
- **WASD**: Movimento della camera (avanti/indietro/sinistra/destra)
- **Q/E**: Movimento verticale (su/giù)
- **ESC**: Chiude l'applicazione

## Utilizzo

```go
package main

import (
    "log"
    "time"
    
    "github.com/alexanderi96/go-space-engine/render/raylib"
    "github.com/alexanderi96/go-space-engine/core/vector"
    // ... altri import
)

func main() {
    // Crea l'adapter Raylib
    adapter := raylib.NewRaylibAdapter(1200, 800, "Titolo Finestra")
    
    // Inizializza l'adapter
    if err := adapter.Initialize(); err != nil {
        log.Fatalf("Errore nell'inizializzazione: %v", err)
    }
    
    // Imposta la posizione della camera
    adapter.SetCameraPosition(vector.NewVector3(0, 50, 150))
    adapter.SetCameraTarget(vector.Zero3())
    
    // Abilita funzionalità di debug
    adapter.SetDebugMode(true)
    adapter.SetRenderVelocities(true)
    
    // Avvia il loop di rendering
    adapter.Run(func(deltaTime time.Duration) {
        // Aggiorna la simulazione
        world.Update(deltaTime)
        
        // Renderizza il mondo
        adapter.RenderWorld(world)
    })
}
```

## Configurazione

### Parametri del costruttore
- `width`: Larghezza della finestra in pixel
- `height`: Altezza della finestra in pixel  
- `title`: Titolo della finestra

### Impostazioni camera
- `cameraSpeed`: Velocità di movimento della camera (default: 50.0)
- `mouseSensitivity`: Sensibilità del mouse (default: 0.003)

## Funzionalità di Debug

L'adapter supporta diverse modalità di debug:

- **Debug Mode**: Mostra i confini del mondo
- **Render Velocities**: Visualizza i vettori di velocità (blu)
- **Render Accelerations**: Visualizza i vettori di accelerazione (rosso)
- **Render Bounding Boxes**: Mostra le bounding box dei corpi (verde)
- **Render Octree**: Visualizza la struttura dell'octree (giallo)

## Rendering dei Materiali

L'adapter mappa automaticamente i materiali ai colori:

- **Sun**: Giallo
- **Mercury**: Grigio chiaro
- **Venus**: Giallo-arancio
- **Earth**: Blu
- **Mars**: Rosso-arancio
- **Jupiter**: Marrone chiaro
- **Saturn**: Giallo pallido
- **Uranus**: Azzurro
- **Neptune**: Blu scuro
- **Spacecraft**: Bianco (renderizzato come cubo)
- **Altri materiali**: Colore generato automaticamente dall'ID

## Dipendenze

- [raylib-go](https://github.com/gen2brain/raylib-go): Binding Go per Raylib
- Raylib: Libreria grafica C (installata automaticamente con raylib-go)

## Note Tecniche

- L'adapter utilizza controlli camera manuali invece dell'orbita automatica di Raylib
- Il rendering è ottimizzato per simulazioni spaziali con sfondo scuro
- La UI mostra sempre le istruzioni e lo stato dei controlli
- Il mouse può essere catturato/rilasciato dinamicamente per un controllo flessibile

## Esempi

Vedi `examples/raylib/solar_system/main.go` per un esempio completo di utilizzo con una simulazione del sistema solare.
