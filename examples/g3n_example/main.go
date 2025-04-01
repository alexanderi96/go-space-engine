// Package main fornisce un esempio di utilizzo di G3N
package main

import (
	"log"
	"time"

	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
)

func main() {
	// Crea il renderer G3N
	r := g3n.NewG3NRenderer()
	err := r.Initialize()
	if err != nil {
		log.Fatalf("Errore durante l'inizializzazione del renderer: %v", err)
	}

	// Ottieni la scena dal renderer
	scene := r.GetScene()

	// Crea una luce ambientale
	ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	scene.Add(ambLight)

	// Crea una luce direzionale
	dirLight := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	dirLight.SetPosition(1, 0, 0)
	scene.Add(dirLight)

	// Crea una sfera
	geom := geometry.NewSphere(0.5, 32, 16)
	mat := material.NewStandard(&math32.Color{0.5, 0.5, 0.8})
	mat.SetShininess(100)
	sphere := graphic.NewMesh(geom, mat)
	scene.Add(sphere)

	// Crea un cubo
	geom2 := geometry.NewBox(0.5, 0.5, 0.5)
	mat2 := material.NewStandard(&math32.Color{0.8, 0.5, 0.5})
	mat2.SetShininess(100)
	cube := graphic.NewMesh(geom2, mat2)
	cube.SetPosition(1, 0, 0)
	scene.Add(cube)

	// Crea gli assi
	axes := helper.NewAxes(1)
	scene.Add(axes)

	// Crea l'interfaccia utente
	gui.Manager().Set(scene)

	// Avvia il loop di rendering
	r.Run(func(deltaTime time.Duration) {
		// Ruota la sfera
		sphere.RotateY(0.01)

		// Ruota il cubo
		cube.RotateY(-0.01)
	})
}
