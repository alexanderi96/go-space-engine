# Piano di Implementazione: Gravità, Octree e Multithreading

Questo documento descrive in dettaglio le modifiche necessarie per implementare l'ottimizzazione della gravità utilizzando l'octree e il multithreading nel motore fisico.

## 1. Estensione dell'Octree per il Calcolo del Centro di Massa

### Modifiche a `physics/space/space.go`

Aggiungeremo i seguenti campi alla struttura `Octree`:

```go
type Octree struct {
    bounds     *AABB       // Limiti dell'octree
    maxObjects int         // Numero massimo di oggetti per nodo
    maxLevels  int         // Numero massimo di livelli
    level      int         // Livello corrente
    objects    []body.Body // Oggetti in questo nodo
    children   [8]*Octree  // Figli dell'octree
    divided    bool        // Indica se l'octree è stato diviso
    
    // Nuovi campi per il calcolo della gravità
    totalMass    float64        // Massa totale di tutti i corpi in questo nodo e nei suoi figli
    centerOfMass vector.Vector3 // Centro di massa di tutti i corpi in questo nodo e nei suoi figli
}
```

Modificheremo i metodi `Insert`, `Remove` e `Clear` per aggiornare il centro di massa e la massa totale:

```go
// Insert inserisce un corpo nell'octree
func (ot *Octree) Insert(b body.Body) {
    // Codice esistente...
    
    // Aggiorna il centro di massa e la massa totale
    ot.updateMassAndCenterOfMass(b, true)
}

// Remove rimuove un corpo dall'octree
func (ot *Octree) Remove(b body.Body) {
    // Codice esistente...
    
    // Aggiorna il centro di massa e la massa totale
    ot.updateMassAndCenterOfMass(b, false)
}

// Clear rimuove tutti i corpi dall'octree
func (ot *Octree) Clear() {
    // Codice esistente...
    
    // Resetta il centro di massa e la massa totale
    ot.totalMass = 0
    ot.centerOfMass = vector.Zero3()
}
```

Aggiungeremo un nuovo metodo per aggiornare il centro di massa e la massa totale:

```go
// updateMassAndCenterOfMass aggiorna il centro di massa e la massa totale
func (ot *Octree) updateMassAndCenterOfMass(b body.Body, adding bool) {
    mass := b.Mass().Value()
    position := b.Position()
    
    if adding {
        // Aggiunge il corpo
        oldTotalMass := ot.totalMass
        ot.totalMass += mass
        
        if oldTotalMass > 0 {
            // Aggiorna il centro di massa
            ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Add(position.Scale(mass)).Scale(1.0 / ot.totalMass)
        } else {
            // Se è il primo corpo, il centro di massa è la sua posizione
            ot.centerOfMass = position
        }
    } else {
        // Rimuove il corpo
        if ot.totalMass > mass {
            // Aggiorna il centro di massa
            oldTotalMass := ot.totalMass
            ot.totalMass -= mass
            ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Sub(position.Scale(mass)).Scale(1.0 / ot.totalMass)
        } else {
            // Se era l'ultimo corpo, resetta il centro di massa
            ot.totalMass = 0
            ot.centerOfMass = vector.Zero3()
        }
    }
    
    // Se l'octree è diviso, propaga l'aggiornamento ai figli
    if ot.divided {
        // Ricalcola il centro di massa dai figli
        ot.totalMass = 0
        weightedPosition := vector.Zero3()
        
        for i := 0; i < 8; i++ {
            if ot.children[i] != nil {
                childMass := ot.children[i].totalMass
                ot.totalMass += childMass
                weightedPosition = weightedPosition.Add(ot.children[i].centerOfMass.Scale(childMass))
            }
        }
        
        if ot.totalMass > 0 {
            ot.centerOfMass = weightedPosition.Scale(1.0 / ot.totalMass)
        }
    }
}
```

## 2. Implementazione dell'Algoritmo Barnes-Hut per il Calcolo della Gravità

Aggiungeremo un nuovo metodo `CalculateGravity` all'Octree per implementare l'algoritmo Barnes-Hut:

```go
// CalculateGravity calcola la forza gravitazionale su un corpo utilizzando l'algoritmo Barnes-Hut
func (ot *Octree) CalculateGravity(b body.Body, theta float64) vector.Vector3 {
    force := vector.Zero3()
    ot.calculateGravityRecursive(b, theta, &force)
    return force
}

// calculateGravityRecursive calcola ricorsivamente la forza gravitazionale
func (ot *Octree) calculateGravityRecursive(b body.Body, theta float64, force *vector.Vector3) {
    // Se l'octree non è diviso o non ha corpi, calcola la forza direttamente
    if !ot.divided || ot.totalMass == 0 {
        ot.calculateLeafNodeGravity(b, force)
        return
    }
    
    // Calcola la larghezza del nodo e la distanza dal corpo al centro di massa
    width := ot.bounds.Max.X() - ot.bounds.Min.X()
    deltaPos := ot.centerOfMass.Sub(b.Position())
    distanceSquared := deltaPos.LengthSquared()
    
    // Se il rapporto larghezza/distanza è inferiore a theta, approssima con il centro di massa
    if (width * width) < (theta * theta * distanceSquared) {
        ot.approximateGravityWithCenterOfMass(b, force)
        return
    }
    
    // Altrimenti, calcola ricorsivamente per ogni figlio
    for i := 0; i < 8; i++ {
        if ot.children[i] != nil {
            ot.children[i].calculateGravityRecursive(b, theta, force)
        }
    }
}

// calculateLeafNodeGravity calcola la forza gravitazionale per ogni corpo nel nodo foglia
func (ot *Octree) calculateLeafNodeGravity(b body.Body, force *vector.Vector3) {
    // Costante gravitazionale
    G := constants.G
    
    // Massa del corpo
    bodyMass := b.Mass().Value()
    bodyPos := b.Position()
    
    // Calcola la forza per ogni corpo nel nodo
    for _, obj := range ot.objects {
        // Evita di calcolare la forza su se stesso
        if obj.ID() == b.ID() {
            continue
        }
        
        // Calcola il vettore direzione
        deltaPos := obj.Position().Sub(bodyPos)
        distanceSquared := deltaPos.LengthSquared()
        
        // Evita divisione per zero
        if distanceSquared <= 0 {
            continue
        }
        
        // Calcola la forza gravitazionale
        distance := math.Sqrt(distanceSquared)
        direction := deltaPos.Scale(1.0 / distance)
        
        // F = G * m1 * m2 / r^2
        forceMagnitude := G * bodyMass * obj.Mass().Value() / distanceSquared
        
        // Aggiungi la forza al vettore forza totale
        *force = force.Add(direction.Scale(forceMagnitude))
    }
}

// approximateGravityWithCenterOfMass approssima la forza gravitazionale usando il centro di massa
func (ot *Octree) approximateGravityWithCenterOfMass(b body.Body, force *vector.Vector3) {
    // Costante gravitazionale
    G := constants.G
    
    // Massa del corpo
    bodyMass := b.Mass().Value()
    bodyPos := b.Position()
    
    // Calcola il vettore direzione
    deltaPos := ot.centerOfMass.Sub(bodyPos)
    distanceSquared := deltaPos.LengthSquared()
    
    // Evita divisione per zero
    if distanceSquared <= 0 {
        return
    }
    
    // Calcola la forza gravitazionale
    distance := math.Sqrt(distanceSquared)
    direction := deltaPos.Scale(1.0 / distance)
    
    // F = G * m1 * m2 / r^2
    forceMagnitude := G * bodyMass * ot.totalMass / distanceSquared
    
    // Aggiungi la forza al vettore forza totale
    *force = force.Add(direction.Scale(forceMagnitude))
}
```

## 3. Implementazione del Multithreading per il Calcolo delle Forze e delle Collisioni

### Modifiche a `simulation/world/world.go`

Aggiungeremo un campo per il pool di worker e modificheremo i metodi `applyForces` e `handleCollisions` per utilizzare il multithreading:

```go
// PhysicalWorld implementa l'interfaccia World
type PhysicalWorld struct {
    // Campi esistenti...
    
    // Nuovo campo per il pool di worker
    workerPool *WorkerPool
}

// NewPhysicalWorld crea un nuovo mondo fisico
func NewPhysicalWorld(bounds *space.AABB) *PhysicalWorld {
    // Codice esistente...
    
    // Crea un pool di worker che si adatta al numero di core disponibili
    workerPool := NewWorkerPool(runtime.NumCPU())
    
    return &PhysicalWorld{
        // Campi esistenti...
        workerPool: workerPool,
    }
}
```

Implementeremo un pool di worker che si adatta al numero di core disponibili:

```go
// WorkerPool rappresenta un pool di worker per il calcolo parallelo
type WorkerPool struct {
    numWorkers int
    tasks      chan func()
    wg         sync.WaitGroup
}

// NewWorkerPool crea un nuovo pool di worker
func NewWorkerPool(numWorkers int) *WorkerPool {
    pool := &WorkerPool{
        numWorkers: numWorkers,
        tasks:      make(chan func(), numWorkers*10), // Buffer per le task
    }
    
    // Avvia i worker
    for i := 0; i < numWorkers; i++ {
        go pool.worker()
    }
    
    return pool
}

// worker esegue le task dal canale
func (wp *WorkerPool) worker() {
    for task := range wp.tasks {
        task()
        wp.wg.Done()
    }
}

// Submit invia una task al pool
func (wp *WorkerPool) Submit(task func()) {
    wp.wg.Add(1)
    wp.tasks <- task
}

// Wait attende che tutte le task siano completate
func (wp *WorkerPool) Wait() {
    wp.wg.Wait()
}
```

Modificheremo il metodo `applyForces` per utilizzare il pool di worker:

```go
// applyForces applica tutte le forze a tutti i corpi
func (w *PhysicalWorld) applyForces() {
    bodies := w.GetBodies()
    
    // Applica le forze globali a tutti i corpi in parallelo
    for _, f := range w.forces {
        if f.IsGlobal() {
            for _, b := range bodies {
                b := b // Cattura la variabile per la goroutine
                w.workerPool.Submit(func() {
                    force := f.Apply(b)
                    b.ApplyForce(force)
                })
            }
        }
    }
    
    // Attendi che tutte le task siano completate
    w.workerPool.Wait()
    
    // Applica le forze tra coppie di corpi in parallelo
    for i := 0; i < len(bodies); i++ {
        for j := i + 1; j < len(bodies); j++ {
            i, j := i, j // Cattura le variabili per la goroutine
            w.workerPool.Submit(func() {
                for _, f := range w.forces {
                    if !f.IsGlobal() {
                        forceA, forceB := f.ApplyBetween(bodies[i], bodies[j])
                        bodies[i].ApplyForce(forceA)
                        bodies[j].ApplyForce(forceB)
                    }
                }
            })
        }
    }
    
    // Attendi che tutte le task siano completate
    w.workerPool.Wait()
}
```

Modificheremo il metodo `handleCollisions` per utilizzare il pool di worker:

```go
// handleCollisions rileva e risolve le collisioni
func (w *PhysicalWorld) handleCollisions() {
    bodies := w.GetBodies()
    
    // Rileva e risolvi le collisioni tra coppie di corpi in parallelo
    for i := 0; i < len(bodies); i++ {
        i := i // Cattura la variabile per la goroutine
        w.workerPool.Submit(func() {
            // Usa la struttura spaziale per trovare potenziali collisioni
            radius := bodies[i].Radius().Value()
            nearbyBodies := w.spatialStructure.QuerySphere(bodies[i].Position(), radius*2)
            
            for _, b := range nearbyBodies {
                // Evita di controllare la collisione con se stesso
                if b.ID() == bodies[i].ID() {
                    continue
                }
                
                // Rileva la collisione
                info := w.collider.CheckCollision(bodies[i], b)
                
                // Risolvi la collisione
                if info.HasCollided {
                    w.collisionResolver.ResolveCollision(info)
                }
            }
            
            // Controlla anche le collisioni con i limiti del mondo
            w.handleBoundaryCollisions(bodies[i])
        })
    }
    
    // Attendi che tutte le task siano completate
    w.workerPool.Wait()
}
```

## 4. Modifiche alla Forza Gravitazionale per Utilizzare l'Octree

### Modifiche a `physics/force/force.go`

Modificheremo la classe `GravitationalForce` per utilizzare l'octree per il calcolo ottimizzato della gravità:

```go
// GravitationalForce implementa la forza gravitazionale
type GravitationalForce struct {
    G     float64 // Costante gravitazionale
    Theta float64 // Parametro di approssimazione per l'algoritmo Barnes-Hut
}

// NewGravitationalForce crea una nuova forza gravitazionale
func NewGravitationalForce() *GravitationalForce {
    return &GravitationalForce{
        G:     constants.G,
        Theta: 0.5, // Valore di theta che bilancia precisione ed efficienza
    }
}

// Apply applica la forza gravitazionale a un corpo (non fa nulla per un singolo corpo)
func (gf *GravitationalForce) Apply(b body.Body) vector.Vector3 {
    // La gravità richiede due corpi per essere applicata
    return vector.Zero3()
}

// ApplyBetween applica la forza gravitazionale tra due corpi
func (gf *GravitationalForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
    // Calcola il vettore direzione da a a b
    direction := b.Position().Sub(a.Position())
    
    // Calcola la distanza al quadrato
    distanceSquared := direction.LengthSquared()
    
    // Evita divisione per zero o forze troppo grandi
    if distanceSquared < 1e-10 {
        return vector.Zero3(), vector.Zero3()
    }
    
    // Normalizza la direzione
    distance := math.Sqrt(distanceSquared)
    normalizedDirection := direction.Scale(1.0 / distance)
    
    // Calcola la forza secondo la legge di gravitazione universale
    // F = G * m1 * m2 / r^2
    massA := a.Mass().Value()
    massB := b.Mass().Value()
    forceMagnitude := gf.G * massA * massB / distanceSquared
    
    // Calcola i vettori forza (direzioni opposte)
    forceOnA := normalizedDirection.Scale(forceMagnitude)
    forceOnB := normalizedDirection.Scale(-forceMagnitude)
    
    return forceOnA, forceOnB
}

// IsGlobal restituisce true perché la gravità è una forza globale
func (gf *GravitationalForce) IsGlobal() bool {
    return true
}

// SetTheta imposta il parametro di approssimazione per l'algoritmo Barnes-Hut
func (gf *GravitationalForce) SetTheta(theta float64) {
    gf.Theta = theta
}

// GetTheta restituisce il parametro di approssimazione per l'algoritmo Barnes-Hut
func (gf *GravitationalForce) GetTheta() float64 {
    return gf.Theta
}
```

## 5. Modifiche all'Esempio G3N per Utilizzare le Nuove Funzionalità

### Modifiche a `examples/g3n/main.go`

Modificheremo l'esempio per utilizzare l'octree e il multithreading:

```go
func main() {
    log.Println("Inizializzazione dell'esempio G3N Physics con Adapter Diretto")
    
    // Crea la configurazione della simulazione
    cfg := config.NewSimulationBuilder().
        WithTimeStep(0.01).
        WithMaxBodies(1000). // Aumentato il numero massimo di corpi
        WithGravity(true).
        WithCollisions(true).
        WithBoundaryCollisions(true).
        WithWorldBounds(
            vector.NewVector3(-100, -100, -100),
            vector.NewVector3(100, 100, 100),
        ).
        WithOctreeConfig(10, 8). // Configurazione ottimale per l'octree
        Build()
    
    // Crea il mondo della simulazione
    w := world.NewPhysicalWorld(cfg.GetWorldBounds())
    
    // Aggiungi la forza gravitazionale
    gravityForce := force.NewGravitationalForce()
    gravityForce.SetTheta(0.5) // Imposta il valore di theta
    w.AddForce(gravityForce)
    
    // Crea un sistema solare realistico
    createSolarSystem(w)
    
    // Crea l'adapter G3N diretto
    adapter := g3n.NewG3NAdapter()
    
    // Configura l'adapter
    adapter.SetBackgroundColor(g3n.NewColor(0.0, 0.0, 0.1, 1.0)) // Sfondo blu scuro per lo spazio
    
    // Variabili per il timing
    lastUpdateTime := time.Now()
    simulationTime := 0.0
    
    // Avvia il loop di rendering
    adapter.Run(func(deltaTime time.Duration) {
        // Calcola il delta time
        currentTime := time.Now()
        dt := currentTime.Sub(lastUpdateTime).Seconds()
        lastUpdateTime = currentTime
        
        // Limita il delta time per evitare instabilità
        if dt > 0.1 {
            dt = 0.1
        }
        
        // Esegui un passo della simulazione
        w.Step(dt)
        simulationTime += dt
        
        // Renderizza il mondo
        adapter.RenderWorld(w)
    })
    
    log.Println("Esempio completato")
}

// createSolarSystem crea un sistema solare realistico
func createSolarSystem(w world.World) {
    log.Println("Creazione del sistema solare")
    
    // Crea il sole
    sun := body.NewRigidBody(
        units.NewQuantity(1.989e30, units.Kilogram), // Massa reale del sole
        units.NewQuantity(696340000, units.Meter),   // Raggio reale del sole (scala ridotta)
        vector.NewVector3(0, 0, 0),
        vector.NewVector3(0, 0, 0),
        physMaterial.NewBasicMaterial(
            "Sun",
            units.NewQuantity(1408, units.Kilogram),
            units.NewQuantity(1000, units.Joule),
            units.NewQuantity(200, units.Watt),
            0.9,
            0.5,
            [3]float64{1.0, 0.8, 0.0}, // Colore giallo
        ),
    )
    sun.SetStatic(true) // Il sole è statico (non si muove)
    w.AddBody(sun)
    log.Printf("Sole creato: ID=%v, Posizione=%v", sun.ID(), sun.Position())
    
    // Crea i pianeti
    createPlanet(w, "Mercury", 0.33e24, 2440, 57.9e9, 47.4, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Venus", 4.87e24, 6052, 108.2e9, 35.0, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Earth", 5.97e24, 6371, 149.6e9, 29.8, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Mars", 0.642e24, 3390, 227.9e9, 24.1, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Jupiter", 1898e24, 69911, 778.6e9, 13.1, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Saturn", 568e24, 58232, 1433.5e9, 9.7, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Uranus", 86.8e24, 25362, 2872.5e9, 6.8, vector.NewVector3(0, 1, 0))
    createPlanet(w, "Neptune", 102e24, 24622, 4495.1e9, 5.4, vector.NewVector3(0, 1, 0))
    
    // Crea un campo di asteroidi
    createAsteroidBelt(w, 1000, 300e9, 500e9)
}

// createPlanet crea un pianeta con parametri realistici
func createPlanet(w world.World, name string, mass, radius, distance, speed float64, orbitPlane vector.Vector3) body.Body {
    // Scala i valori per la simulazione
    scaleFactor := 1e-9 // Scala le distanze
    massScale := 1e-20  // Scala le masse
    
    // Calcola la posizione iniziale
    position := vector.NewVector3(distance*scaleFactor, 0, 0)
    
    // Calcola la velocità orbitale (perpendicolare alla posizione)
    velocity := orbitPlane.Cross(position).Normalize().Scale(speed * scaleFactor)
    
    // Crea il pianeta
    planet := body.NewRigidBody(
        units.NewQuantity(mass*massScale, units.Kilogram),
        units.NewQuantity(radius*scaleFactor*10, units.Meter), // Aumenta il raggio per la visualizzazione
        position,
        velocity,
        physMaterial.NewBasicMaterial(
            name,
            units.NewQuantity(5000, units.Kilogram),
            units.NewQuantity(800, units.Joule),
            units.NewQuantity(1.5, units.Watt),
            0.7,
            0.5,
            getPlanetColor(name),
        ),
    )
    
    // Aggiungi il pianeta al mondo
    w.AddBody(planet)
    log.Printf("Pianeta %s aggiunto: ID=%v, Posizione=%v, Velocità=%v", name, planet.ID(), planet.Position(), planet.Velocity())
    
    return planet
}

// createAsteroidBelt crea un campo di asteroidi
func createAsteroidBelt(w world.World, count int, minDistance, maxDistance float64) {
    log.Printf("Creazione di %d asteroidi", count)
    
    // Scala i valori per la simulazione
    scaleFactor := 1e-9 // Scala le distanze
    massScale := 1e-24  // Scala le masse
    
    for i := 0; i < count; i++ {
        // Genera una posizione casuale nel campo di asteroidi
        distance := minDistance + rand.Float64()*(maxDistance-minDistance)
        angle := rand.Float64() * 2 * math.Pi
        
        x := distance * math.Cos(angle) * scaleFactor
        z := distance * math.Sin(angle) * scaleFactor
        y := (rand.Float64()*2 - 1) * 10e9 * scaleFactor // Distribuzione verticale
        
        position := vector.NewVector3(x, y, z)
        
        // Calcola la velocità orbitale
        speed := math.Sqrt(constants.G * 1.989e30 * massScale / distance) * 1e5
        velocity := vector.NewVector3(-z, 0, x).Normalize().Scale(speed * scaleFactor)
        
        // Crea l'asteroide
        asteroid := body.NewRigidBody(
            units.NewQuantity(rand.Float64()*1e16*massScale, units.Kilogram),
            units.NewQuantity(rand.Float64()*1000*scaleFactor, units.Meter),
            position,
            velocity,
            physMaterial.Rock,
        )
        
        // Aggiungi l'asteroide al mondo
        w.AddBody(asteroid)
    }
    
    log.Printf("Campo di asteroidi creato")
}

// getPlanetColor restituisce un colore appropriato per il pianeta
func getPlanetColor(name string) [3]float64 {
    switch name {
    case "Mercury":
        return [3]float64{0.7, 0.7, 0.7} // Grigio
    case "Venus":
        return [3]float64{0.9, 0.7, 0.0} // Giallo-arancio
    case "Earth":
        return [3]float64{0.0, 0.3, 0.8} // Blu
    case "Mars":
        return [3]float64{0.8, 0.3, 0.0} // Rosso
    case "Jupiter":
        return [3]float64{0.8, 0.6, 0.4} // Marrone chiaro
    case "Saturn":
        return [3]float64{0.9, 0.8, 0.5} // Giallo-marrone
    case "Uranus":
        return [3]float64{0.5, 0.8, 0.9} // Azzurro
    case "Neptune":
        return [3]float64{0.0, 0.0, 0.8} // Blu scuro
    default:
        return [3]float64{0.5, 0.5, 0.5} // Grigio
    }
}
```

## Conclusione

Queste modifiche implementeranno l'ottimizzazione della gravità utilizzando l'octree e il multithreading nel motore fisico. L'algoritmo Barnes-Hut ridurrà la complessità del calcolo della gravità da O(n²) a O(n log n), permettendo di simulare un numero molto maggiore di corpi. Il multithreading permetterà di sfruttare tutti i core della CPU, migliorando ulteriormente le prestazioni.

L'esempio G3N mostrerà un sistema solare realistico con pianeti e un campo di asteroidi, dimostrando le capacità del motore fisico.