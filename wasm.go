//go:build js && wasm
package main

import "fmt"
import "syscall/js"
import "math"
import "encoding/json"

var (
	u	[]float64
	v	[]float64
	xDomain []float64
	sinWave []float64
	dx	float64
	dt	float64
	CFL	float64
	n	uint32
)

type DerivFunc func([]float64) []float64

func init() {
	initializeData()
	js.Global().Set("getTimeStep", js.FuncOf(getTimeStep))
	js.Global().Set("runWaveEquation", js.FuncOf(runWaveEquation))
}

func initializeData() {
	n = 100
	dx = float64(2 * math.Pi / float64(n-1))
	u = make([]float64, n)
	v = make([]float64, n)
	xDomain = make([]float64, n)
	sinWave = make([]float64, n)
	for i := range xDomain {
		xDomain[i] = float64(i) * dx
		u[i] = math.Sin(xDomain[i])
		v[i] = 0
	}
	u[0] = 0
	u[n-1] = 0
	CFL = 0.01
	dt = CFL * dx
}

func logToJS(v interface{}){
	msg := fmt.Sprintf("%+v", v)
	js.Global().Get("console").Call("log", msg)
}

func getTimeStep(this js.Value, args[] js.Value) interface{} {
	return js.ValueOf(dt)
}

func runWaveEquation(this js.Value, args []js.Value) interface{} {
	steps := args[0].Int()
	logToJS(fmt.Sprintf("Running simulation for %d steps", steps))
	for step := 0; step < steps; step++ {
		u, v = waveRungeKutta4(u, v, dt)
		u[0] = 0
		u[n-1] = 0
	}
	data := struct {
                  X []float64 `json:"x"`
                  F []float64 `json:"f"`
          }{
                  X: xDomain,
                  F: u,
          }
          jsonData, err := json.Marshal(data)
          if err != nil {
                  fmt.Println("Error marshalling data: ", err)
                  return nil
          }
          return js.ValueOf(string(jsonData))
}

func dudt(v []float64) []float64 {
	return v
}

func dvdt(u []float64) []float64 {
	d2udx2 := secondDerivativeCentralDiff(u, dx)
	return d2udx2
}

func waveRungeKutta4(u []float64, v []float64, dt float64) ([]float64, []float64) {
	n := len(u)
	k1u := make([]float64, n)
	k1v := make([]float64, n)
	k2u := make([]float64, n)
	k2v := make([]float64, n)
	k3u := make([]float64, n)
	k3v := make([]float64, n)
	k4u := make([]float64, n)
	k4v := make([]float64, n)
	uTemp := make([]float64, n)
	vTemp := make([]float64, n)
	uNext := make([]float64, n)
	vNext := make([]float64, n)

	// Calculate k1 = dt * f'(f)
	for i, val := range dudt(v) {
		k1u[i] = dt * val
	}
	for i, val := range dvdt(u) {
		k1v[i] = dt * val
	}

	// Calculate k2 = dt * f'(f+k1/2)
	for i := range u {
		uTemp[i] = u[i] + k1u[i] / 2
	}
	for i := range v {
		vTemp[i] = v[i] + k1v[i] / 2
	}
	for i, val := range dudt(vTemp) {
		k2u[i] = dt * val
	}
	for i, val := range dvdt(uTemp) {
		k2v[i] = dt * val
	}

	// Calculate k3 = dt * f'(f+k2/2)
	for i := range u {
		uTemp[i] = u[i] + k2u[i] / 2
	}
	for i := range v {
		vTemp[i] = v[i] + k2v[i] / 2
	}
	for i, val := range dudt(vTemp) {
		k3u[i] = dt * val
	}
	for i, val := range dvdt(uTemp) {
		k3v[i] = dt * val
	}

	// Calculate k4 = dt * f'(f+k3)
	for i := range u {
		uTemp[i] = u[i] + k3u[i]
	}
	for i := range v {
		vTemp[i] = v[i] + k3v[i]
	}
	for i, val := range dudt(vTemp) {
		k4u[i] = dt * val
	}
	for i, val := range dvdt(uTemp) {
		k4v[i] = dt * val
	}

	for i := range uNext {
		uNext[i] = u[i] + (k1u[i] + 2*k2u[i] + 2*k3u[i] + k4u[i]) / 6
		vNext[i] = v[i] + (k1v[i] + 2*k2v[i] + 2*k3v[i] + k4v[i]) / 6
	}
	return uNext, vNext
}

func rungeKutta4(f []float64, g []float64, dt float64, dgdtFunc DerivFunc) []float64 {
	n := len(f)
	k1 := make([]float64, n)
	k2 := make([]float64, n)
	k3 := make([]float64, n)
	k4 := make([]float64, n)
	fNext := make([]float64, n)
	temp := make([]float64, n)

	// Calculate k1 = dt * f'(f)
	for i, val := range dgdtFunc(g) {
		k1[i] = dt * val
	}
	// Calculate k2 = dt * f'(f+k1/2)
	for i := range f {
		temp[i] = f[i] + k1[i]/2
	}
	for i, val := range dgdtFunc(temp) {
		k2[i] = dt * val
	}
	// Calculate k3 = ft * f'(f+k2/2)
	for i := range f {
                  temp[i] = f[i] + k2[i]/2
        }
        for i, val := range dgdtFunc(temp) {
                  k3[i] = dt * val
	}
	// Calculate k4 = dt * f'(f+k3)
	for i := range f {
        	temp[i] = f[i] + k3[i]
    	}
    	for i, val := range dgdtFunc(temp) {
        	k4[i] = dt * val
    	}
	for i := range fNext {
		fNext[i] = f[i] + (k1[i] + 2*k2[i] + 2*k3[i] + k4[i]) / 6
	}
	return fNext
}

func firstDerivativeCentralDiff(f []float64, h float64) []float64 {
	n := len(f)
	df := make([]float64, n)

	// Central Difference For Interior Points
	for i := 1; i < n-1; i++ {
		df[i] = (f[i+1] - f[i-1]) / (2 * h)
	}

	// Edge Forward/Backward Differences
	df[0] = (-3*f[0] + 4*f[1] - f[2]) / (2*h)
	df[n-1] = (3*f[n-1] - 4*f[n-2] + f[n-3]) / (2*h)
	
	return df
}

func secondDerivativeCentralDiff(f []float64, h float64) []float64 {
	n := len(f)
	d2f := make([]float64, n)
	
	for i :=1; i < n-1; i++ {
		d2f[i] = (f[i-1]-2*f[i]+f[i+1]) / (h*h)
	}
	return d2f
}

func getSineWave(this js.Value, args []js.Value) interface{} {
	data := struct {
		X []float64 `json:"x"`
		SinWave []float64 `json:"sinWave"`
	}{
		X: xDomain,
		SinWave: sinWave,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data: ", err)
		return nil
	}
	return js.ValueOf(string(jsonData))
}

func getFirstDerivative(this js.Value, args []js.Value) interface{} {
	df := firstDerivativeCentralDiff(sinWave, dx)
	data := struct {
		X  []float64 `json:"x"`
		Df []float64 `json:"df"`
	}{
		X:  xDomain,
		Df: df,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return nil
	}
	return js.ValueOf(string(jsonData))
}
