// // Package g3n fornisce un'implementazione dell'interfaccia Renderer utilizzando G3N
package g3n

// import (
// 	"time"

// 	"github.com/alexanderi96/go-space-engine/core/vector"
// 	"github.com/alexanderi96/go-space-engine/physics/body"
// 	"github.com/alexanderi96/go-space-engine/physics/space"
// 	"github.com/alexanderi96/go-space-engine/render/adapter"

// 	"github.com/g3n/engine/app"
// 	"github.com/g3n/engine/camera"
// 	"github.com/g3n/engine/core"
// 	"github.com/g3n/engine/geometry"
// 	"github.com/g3n/engine/gls"
// 	"github.com/g3n/engine/graphic"
// 	"github.com/g3n/engine/light"
// 	"github.com/g3n/engine/material"
// 	"github.com/g3n/engine/math32"
// 	"github.com/g3n/engine/renderer"
// 	"github.com/g3n/engine/util/helper"

// 	"sync"
// )

// // NewColor crea un nuovo colore
// func NewColor(r, g, b, a float64) adapter.Color {
// 	return adapter.NewColor(r, g, b, a)
// }

// // G3NRenderer implementa l'interfaccia Renderer utilizzando G3N
// type G3NRenderer struct {
// 	app        *app.Application
// 	scene      *core.Node
// 	camera     *camera.Camera
// 	cameraCtrl *camera.OrbitControl

// 	// Mappa per tenere traccia degli oggetti grafici associati ai corpi fisici
// 	bodyNodes map[body.ID]*BodyMesh
// 	nodeMutex sync.RWMutex

// 	// Altre proprietà necessarie
// 	width, height int
// 	running       bool
// 	bgColor       adapter.Color
// }

// // BodyMesh rappresenta un mesh con una luce puntuale associata
// type BodyMesh struct {
// 	Mesh  *graphic.Mesh
// 	Light *light.Point
// }

// // NewG3NRenderer crea un nuovo renderer G3N
// func NewG3NRenderer() *G3NRenderer {
// 	return &G3NRenderer{
// 		bodyNodes: make(map[body.ID]*BodyMesh),
// 		width:     800,
// 		height:    600,
// 		running:   false,
// 		bgColor:   adapter.NewColor(0.7, 0.7, 0.7, 1.0), // Colore di sfondo grigio chiaro
// 	}
// }

// // Initialize inizializza il renderer
// func (r *G3NRenderer) Initialize() error {
// 	// Crea l'applicazione G3N
// 	r.app = app.App()

// 	// Ottieni le dimensioni della finestra
// 	r.width, r.height = r.app.GetSize()

// 	// Imposta il colore di sfondo
// 	r.app.Gls().ClearColor(float32(r.bgColor.R), float32(r.bgColor.G), float32(r.bgColor.B), float32(r.bgColor.A))

// 	// Crea la scena
// 	r.scene = core.NewNode()

// 	// Crea la camera
// 	r.camera = camera.New(1)         // 1 = Perspective
// 	r.camera.SetPosition(0, 15, 150) // Posiziona la camera come nell'esempio diretto
// 	r.camera.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
// 	aspect := float32(r.width) / float32(r.height)
// 	r.camera.SetAspect(aspect)

// 	// Aggiungi la camera alla scena
// 	r.scene.Add(r.camera)

// 	// Crea il controllo orbitale della camera
// 	r.cameraCtrl = camera.NewOrbitControl(r.camera)

// 	// Aggiungi luci con maggiore intensità
// 	ambientLight := light.NewAmbient(&math32.Color{0.8, 0.8, 0.8}, 1.0) // Stessa intensità dell'esempio diretto
// 	r.scene.Add(ambientLight)

// 	// Aggiungi più luci puntuali per illuminare la scena da diverse angolazioni
// 	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
// 	pointLight1.SetPosition(10, 10, 10)
// 	r.scene.Add(pointLight1)

// 	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
// 	pointLight2.SetPosition(-10, 10, 10)
// 	r.scene.Add(pointLight2)

// 	pointLight3 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
// 	pointLight3.SetPosition(0, 10, -10)
// 	r.scene.Add(pointLight3)

// 	// Aggiungi assi e griglia come nell'esempio diretto
// 	axes := helper.NewAxes(2)
// 	r.scene.Add(axes)

// 	grid := helper.NewGrid(20, 1, &math32.Color{0.4, 0.4, 0.4})
// 	r.scene.Add(grid)

// 	r.running = true

// 	return nil
// }

// // Shutdown chiude il renderer
// func (r *G3NRenderer) Shutdown() error {
// 	r.running = false
// 	r.app.Exit()
// 	return nil
// }

// // BeginFrame inizia un nuovo frame
// func (r *G3NRenderer) BeginFrame() {
// 	// Il frame viene iniziato automaticamente da G3N
// }

// // EndFrame termina il frame corrente
// func (r *G3NRenderer) EndFrame() {
// 	// Assicurati che la scena e la camera siano state inizializzate
// 	if r.scene == nil || r.camera == nil {
// 		return
// 	}

// 	// Nota: quando si usa il metodo Run, il rendering viene gestito internamente
// 	// da G3N, quindi non è necessario chiamare Render qui.
// 	// Questo metodo è mantenuto per compatibilità con l'interfaccia Renderer.
// }

// // RenderBody renderizza un corpo
// func (r *G3NRenderer) RenderBody(b body.Body) {
// 	r.nodeMutex.Lock()
// 	defer r.nodeMutex.Unlock()

// 	// Controlla se il corpo è già stato renderizzato
// 	if _, exists := r.bodyNodes[b.ID()]; !exists {
// 		// Crea un nuovo mesh per il corpo
// 		bodyMesh := r.createBodyMesh(b)
// 		r.bodyNodes[b.ID()] = bodyMesh
// 	} else {
// 		// Aggiorna la posizione del mesh
// 		r.updateMeshFromBody(b)
// 	}
// }

// // RenderBodies renderizza tutti i corpi
// func (r *G3NRenderer) RenderBodies(bodies []body.Body) {
// 	// Crea un set di ID dei corpi correnti
// 	currentIDs := make(map[body.ID]bool)
// 	for _, b := range bodies {
// 		currentIDs[b.ID()] = true
// 		r.RenderBody(b)
// 	}

// 	// Rimuovi i mesh dei corpi che non esistono più
// 	r.nodeMutex.Lock()
// 	defer r.nodeMutex.Unlock()

// 	for id, bodyMesh := range r.bodyNodes {
// 		if !currentIDs[id] {
// 			r.scene.Remove(bodyMesh.Mesh)
// 			if bodyMesh.Light != nil {
// 				r.scene.Remove(bodyMesh.Light)
// 			}
// 			delete(r.bodyNodes, id)
// 		}
// 	}
// }

// // createBodyMesh crea un mesh per un corpo fisico
// func (r *G3NRenderer) createBodyMesh(b body.Body) *BodyMesh {
// 	// Crea una sfera per rappresentare il corpo
// 	radius := float32(b.Radius().Value())
// 	geom := geometry.NewSphere(float64(radius), 32, 16)

// 	// Crea un materiale in base al materiale del corpo fisico
// 	mat := material.NewStandard(&math32.Color{0.8, 0.8, 0.8})

// 	// Se il materiale del corpo ha un colore, usalo
// 	var bodyColor math32.Color
// 	if b.Material() != nil {
// 		// Qui dovresti mappare il materiale fisico a un colore G3N
// 		// Per semplicità, usiamo un colore predefinito per ogni tipo di materiale
// 		switch b.Material().Name() {
// 		case "Iron":
// 			bodyColor = math32.Color{0.6, 0.6, 0.6}
// 		case "Rock":
// 			bodyColor = math32.Color{0.5, 0.3, 0.2}
// 		case "Ice":
// 			bodyColor = math32.Color{0.8, 0.9, 1.0}
// 		case "Copper":
// 			bodyColor = math32.Color{0.8, 0.5, 0.2}
// 		default:
// 			// Colore casuale basato sull'ID del corpo (che è una stringa)
// 			id := string(b.ID())
// 			hash := 0
// 			for i := 0; i < len(id); i++ {
// 				hash = 31*hash + int(id[i])
// 			}
// 			if hash < 0 {
// 				hash = -hash
// 			}
// 			r := float32(hash%255) / 255.0
// 			g := float32((hash/255)%255) / 255.0
// 			b := float32((hash/(255*255))%255) / 255.0
// 			bodyColor = math32.Color{r, g, b}
// 		}
// 		mat.SetColor(&bodyColor)
// 	} else {
// 		bodyColor = math32.Color{0.8, 0.8, 0.8}
// 	}

// 	// Crea un mesh con la geometria e il materiale
// 	mesh := graphic.NewMesh(geom, mat)

// 	// Imposta la posizione del mesh
// 	pos := b.Position()
// 	mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

// 	// Crea una luce puntuale per il corpo
// 	bodyLight := light.NewPoint(&bodyColor, 0.5)
// 	bodyLight.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
// 	bodyLight.SetLinearDecay(1.0)
// 	bodyLight.SetQuadraticDecay(1.0)
// 	r.scene.Add(bodyLight)

// 	// Aggiungi il mesh alla scena
// 	r.scene.Add(mesh)

// 	// Crea un BodyMesh
// 	bodyMesh := &BodyMesh{
// 		Mesh:  mesh,
// 		Light: bodyLight,
// 	}

// 	return bodyMesh
// }

// // updateMeshFromBody aggiorna un mesh in base allo stato del corpo fisico
// func (r *G3NRenderer) updateMeshFromBody(b body.Body) {
// 	bodyMesh := r.bodyNodes[b.ID()]
// 	if bodyMesh == nil {
// 		return
// 	}

// 	// Aggiorna la posizione
// 	pos := b.Position()
// 	bodyMesh.Mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
// 	if bodyMesh.Light != nil {
// 		bodyMesh.Light.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
// 	}
// }

// // RenderAABB renderizza un AABB
// func (r *G3NRenderer) RenderAABB(aabb *space.AABB, color adapter.Color) {
// 	// Implementazione della renderizzazione di un AABB
// 	min := aabb.Min
// 	max := aabb.Max

// 	// Crea un cubo per rappresentare l'AABB
// 	width := float32(max.X() - min.X())
// 	height := float32(max.Y() - min.Y())
// 	depth := float32(max.Z() - min.Z())

// 	geom := geometry.NewBox(width, height, depth)
// 	mat := material.NewStandard(&math32.Color{float32(color.R), float32(color.G), float32(color.B)})
// 	mat.SetOpacity(float32(color.A))
// 	mat.SetTransparent(true)

// 	mesh := graphic.NewMesh(geom, mat)
// 	mesh.SetPosition(float32(min.X()+max.X())/2, float32(min.Y()+max.Y())/2, float32(min.Z()+max.Z())/2)

// 	// Aggiungi il mesh alla scena temporaneamente (verrà rimosso alla fine del frame)
// 	r.scene.Add(mesh)

// 	// Aggiungi una funzione di pulizia per rimuovere il mesh alla fine del frame
// 	// Questo è un po' un hack, ma funziona per una demo
// 	go func() {
// 		// Aspetta che il frame sia completato
// 		time.Sleep(time.Millisecond * 100)
// 		r.scene.Remove(mesh)
// 	}()
// }

// // RenderOctree renderizza un octree
// func (r *G3NRenderer) RenderOctree(octree *space.Octree, maxDepth int) {
// 	// Implementazione della renderizzazione di un octree
// 	// Per semplicità, renderizziamo solo l'AABB dell'octree
// 	// Non possiamo accedere direttamente al campo bounds dell'octree
// 	// Per ora, non renderizziamo l'octree
// 	if octree != nil {
// 		// Creiamo un AABB che copre l'intera scena
// 		min := vector.NewVector3(-10, -10, -10)
// 		max := vector.NewVector3(10, 10, 10)
// 		aabb := space.NewAABB(min, max)
// 		r.RenderAABB(aabb, adapter.NewColor(0.5, 0.5, 0.5, 0.3))
// 	}
// }

// // RenderLine renderizza una linea
// func (r *G3NRenderer) RenderLine(start, end vector.Vector3, color adapter.Color) {
// 	// Implementazione della renderizzazione di una linea
// 	// Crea una geometria per la linea
// 	geom := geometry.NewGeometry()

// 	// Aggiungi i vertici
// 	positions := math32.NewArrayF32(0, 6)
// 	positions.Append(float32(start.X()), float32(start.Y()), float32(start.Z()))
// 	positions.Append(float32(end.X()), float32(end.Y()), float32(end.Z()))

// 	// Imposta gli attributi della geometria
// 	geom.AddVBO(
// 		gls.NewVBO(positions).
// 			AddAttrib(gls.VertexPosition),
// 	)

// 	// Crea un materiale per la linea con il colore specificato
// 	mat := material.NewStandard(&math32.Color{float32(color.R), float32(color.G), float32(color.B)})
// 	mat.SetTransparent(true)
// 	mat.SetOpacity(float32(color.A))

// 	// Crea una linea con la geometria e il materiale
// 	line := graphic.NewLines(geom, mat)

// 	// Aggiungi la linea alla scena temporaneamente
// 	r.scene.Add(line)

// 	// Aggiungi una funzione di pulizia per rimuovere la linea alla fine del frame
// 	go func() {
// 		// Aspetta che il frame sia completato
// 		time.Sleep(time.Millisecond * 100)
// 		r.scene.Remove(line)
// 	}()
// }

// // RenderSphere renderizza una sfera
// func (r *G3NRenderer) RenderSphere(center vector.Vector3, radius float64, color adapter.Color) {
// 	// Implementazione della renderizzazione di una sfera
// 	geom := geometry.NewSphere(radius, 16, 8)
// 	mat := material.NewStandard(&math32.Color{float32(color.R), float32(color.G), float32(color.B)})
// 	mat.SetOpacity(float32(color.A))
// 	mat.SetTransparent(true)

// 	mesh := graphic.NewMesh(geom, mat)
// 	mesh.SetPosition(float32(center.X()), float32(center.Y()), float32(center.Z()))

// 	// Aggiungi il mesh alla scena temporaneamente
// 	r.scene.Add(mesh)

// 	// Aggiungi una funzione di pulizia per rimuovere il mesh alla fine del frame
// 	go func() {
// 		// Aspetta che il frame sia completato
// 		time.Sleep(time.Millisecond * 100)
// 		r.scene.Remove(mesh)
// 	}()
// }

// // SetCamera imposta la posizione e l'orientamento della camera
// func (r *G3NRenderer) SetCamera(position, target, up vector.Vector3) {
// 	r.camera.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))
// 	r.camera.LookAt(&math32.Vector3{float32(target.X()), float32(target.Y()), float32(target.Z())}, &math32.Vector3{float32(up.X()), float32(up.Y()), float32(up.Z())})
// }

// // SetCameraFOV imposta il campo visivo della camera
// func (r *G3NRenderer) SetCameraFOV(fov float64) {
// 	r.camera.SetFov(float32(fov))
// }

// // SetBackgroundColor imposta il colore di sfondo
// func (r *G3NRenderer) SetBackgroundColor(color adapter.Color) {
// 	r.bgColor = color
// 	// Applica il colore di sfondo all'applicazione G3N
// 	if r.app != nil {
// 		r.app.Gls().ClearColor(float32(color.R), float32(color.G), float32(color.B), float32(color.A))
// 	}
// }

// // GetWidth restituisce la larghezza della finestra di rendering
// func (r *G3NRenderer) GetWidth() int {
// 	return r.width
// }

// // GetHeight restituisce l'altezza della finestra di rendering
// func (r *G3NRenderer) GetHeight() int {
// 	return r.height
// }

// // IsRunning restituisce true se il renderer è in esecuzione
// func (r *G3NRenderer) IsRunning() bool {
// 	return r.running
// }

// // ProcessEvents processa gli eventi del renderer
// func (r *G3NRenderer) ProcessEvents() {
// 	// Processa gli eventi della finestra
// 	// Questo è un po' un hack, ma funziona per una demo
// 	time.Sleep(time.Millisecond * 16) // Circa 60 FPS
// }

// // Run avvia il loop di rendering utilizzando il loop interno di G3N
// func (r *G3NRenderer) Run(updateFunc func(deltaTime time.Duration)) {
// 	if !r.running {
// 		if err := r.Initialize(); err != nil {
// 			panic(err)
// 		}
// 	}

// 	// Avvia il loop di rendering
// 	r.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
// 		// Chiama la funzione di aggiornamento fornita
// 		if updateFunc != nil {
// 			updateFunc(deltaTime)
// 		}

// 		// La renderizzazione della scena viene gestita dall'adapter
// 		// che chiama esplicitamente le operazioni OpenGL
// 		// Qui non chiamiamo direttamente renderer.Render
// 	})
// }

// // RenderScene renderizza la scena
// func (r *G3NRenderer) RenderScene() {
// 	if r.app == nil || r.scene == nil || r.camera == nil {
// 		return
// 	}

// 	// Ottieni il renderer G3N
// 	renderer := r.app.Renderer()
// 	if renderer == nil {
// 		return
// 	}

// 	// Renderizza la scena utilizzando il renderer G3N standard
// 	// Questo è lo stesso approccio utilizzato nell'esempio diretto
// 	renderer.Render(r.scene, r.camera)
// }

// // GetScene restituisce la scena G3N
// func (r *G3NRenderer) GetScene() *core.Node {
// 	return r.scene
// }

// // GetApp restituisce l'applicazione G3N
// func (r *G3NRenderer) GetApp() *app.Application {
// 	return r.app
// }

// // GetCamera restituisce la camera G3N
// func (r *G3NRenderer) GetCamera() *camera.Camera {
// 	return r.camera
// }

// // GetCameraControl restituisce il controllo della camera G3N
// func (r *G3NRenderer) GetCameraControl() *camera.OrbitControl {
// 	return r.cameraCtrl
// }

// // AddAxes aggiunge gli assi alla scena
// func (r *G3NRenderer) AddAxes(size float32) {
// 	axes := helper.NewAxes(size)
// 	r.scene.Add(axes)
// }

// // AddGrid aggiunge una griglia alla scena
// func (r *G3NRenderer) AddGrid(size float32, step float32, color adapter.Color) {
// 	grid := helper.NewGrid(size, step, &math32.Color{float32(color.R), float32(color.G), float32(color.B)})
// 	r.scene.Add(grid)
// }

// // AddLight aggiunge una luce puntuale alla scena
// func (r *G3NRenderer) AddLight(position vector.Vector3, color adapter.Color, intensity float32) {
// 	pointLight := light.NewPoint(&math32.Color{float32(color.R), float32(color.G), float32(color.B)}, intensity)
// 	pointLight.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))
// 	r.scene.Add(pointLight)
// }
