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
	xDomain = make([]float64, n)
	sinWave = make([]float64, n)
	for i := range xDomain {
		xDomain[i] = float64(i) * dx
		sinWave[i] = math.Sin(xDomain[i])
	}
	u = sinWave
	v = sinWave
	u[0] = 0
	u[n-1] = 0
	CFL = 0.15
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
		u = rungeKutta4(u, dt, dudt)
		v = rungeKutta4(v, dt, dvdt)
		u[0] = 0
		u[len(u)-1] = 0
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
	c2 := 1.0
	d2udx2 := firstDerivativeCentralDiff(firstDerivativeCentralDiff(u, dx), dx)
	for i := range d2udx2 {
		d2udx2[i] *=c2
	}
	return d2udx2
}

func rungeKutta4(f []float64, dt float64, dfdtFunc DerivFunc) []float64 {
	n := len(f)
	k1 := make([]float64, n)
	k2 := make([]float64, n)
	k3 := make([]float64, n)
	k4 := make([]float64, n)
	fNext := make([]float64, n)
	temp := make([]float64, n)

	// Calculate k1 = dt * f'(f)
	for i, val := range dfdtFunc(f) {
		k1[i] = dt * val
	}
	// Calculate k2 = dt * f'(f+k1/2)
	for i := range f {
		temp[i] = f[i] + k1[i]/2
	}
	for i, val := range dfdtFunc(temp) {
		k2[i] = dt * val
	}
	// Calculate k3 = ft * f'(f+k2/2)
	for i := range f {
                  temp[i] = f[i] + k2[i]/2
        }
        for i, val := range dfdtFunc(temp) {
                  k3[i] = dt * val
	}
	// Calculate k4 = dt * f'(f+k3)
	for i := range f {
        	temp[i] = f[i] + k3[i]
    	}
    	for i, val := range dfdtFunc(temp) {
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
