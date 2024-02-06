//go:build js && wasm
package main

import "fmt"
import "syscall/js"
import "math"
import "encoding/json"


var (
	f []float64
	fOld []float64
	xDomain []float64
	sinWave []float64
	dx	float64
	dt	float64
	CFL	float64
)

func init() {
	js.Global().Set("getSineWave", js.FuncOf(getSineWave))
	js.Global().Set("getFirstDerivative", js.FuncOf(getFirstDerivative))
	initializeData()
}

func initializeData() {
	const n = 100
	dx = float64(2 * math.Pi / float64(n-1))
	xDomain = make([]float64, n)
	sinWave = make([]float64, n)
	for i := range xDomain {
		xDomain[i] = float64(i) * dx
		sinWave[i] = math.Sin(xDomain[i])
	}
	f = sinWave
	fOld = sinWave
	CFL = 0.2
	dt = CFL * dx
}

func rungeKutta4(dfdt []float64, f []float64, dt float64) []float64 {
	
}

func applyTimeStep(f []float64, dt float64) []float64 {
	
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
	// Index adjustment
	n = n - 1
	df[n] = (3*f[n] - 4*f[n-1] + f[n-2]) / (2*h)
	
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
