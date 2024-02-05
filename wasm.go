//go:build js && wasm
package main

import "fmt"
import "syscall/js"
import "math"
import "encoding/json"


func init() {
	js.Global().Set("getSineWave", js.FuncOf(getSineWave))
}

func getSineWave(this js.Value, args []js.Value) interface{} {
	const n = 100
	xDomain := make([]float64, n)
	sinWave := make([]float64, n)
	dx := float64(2 * math.Pi / float64(n-1))
	for i := range xDomain {
		xDomain[i] = float64(i) * dx
		sinWave[i] = float64(math.Sin(float64((xDomain[i]))))
	}
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
