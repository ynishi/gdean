package service

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"os/exec"
)

type MaxEmvInput struct {
	P1  float32 `json:"p1"`
	DS1 []int32 `json:"ds1"`
	DS2 []int32 `json:"ds2"`
}

func calcMaxEmv(p1 float32, dataP1 []int32, dataP2 []int32) (int32, error) {
	mei := MaxEmvInput{P1: p1, DS1: dataP1, DS2: dataP2}
	d, _ := json.Marshal(mei)
	var meo map[string]int32
	ret := runPython("max_emv", string(d))
	if err := json.Unmarshal([]byte(ret), &meo); err != nil {
		fmt.Println(err)
		return -1, err
	}
	return meo["ans"], nil
}

func calcMaxEmvGo(p1 float32, dataP1 []int32, dataP2 []int32) (int32, error) {
	p2 := 1 - p1
	l := len(dataP1)
	dataP12 := make([]float32, l)
	dataP22 := make([]float32, l)
	sums := make([]float32, l)
	for i := 0; i < l; i++ {
		dataP12[i] = float32(dataP1[i]) * p1
		dataP22[i] = float32(dataP2[i]) * p2
		sums[i] = dataP12[i] + dataP22[i]
	}
	maxSum := funk.MaxFloat32(sums)
	maxSumI := funk.IndexOf(sums, maxSum)

	return int32(maxSumI), nil

}

func runPython(name string, payload string) string {
	// expect valid python environment exists
	out1, _ := exec.Command("ls", "calc/calc.py").Output()
	fmt.Println(string(out1))

	out, err := exec.Command("python", "calc/calc.py", "{\"name\":\""+name+"\",\"data\":"+payload+"}").Output()
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(out)
}
