package borderlessmap

import (
	"math/rand"
	"sync"
	"testing"
)

type Pkg struct {
	id   int
	data []byte
}

func benchmarkFanOut(outChNum, dataNum, bufNum, pkgSize int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		fanout := &FanOut{}
		in := fanout.Init(bufNum)
		out := make([]<-chan interface{}, outChNum)
		result := make(map[int]int)
		getPackageID := func(pkg interface{}) int {
			m := pkg.(Pkg)
			return m.id
		}
		rec := make(chan interface{}, dataNum*outChNum)
		for i := 0; i < outChNum; i++ {
			out[i] = fanout.Register(bufNum)
		}
		wg := sync.WaitGroup{}
		for i := range out {
			wg.Add(dataNum) // receive
			x := i
			go func() {
				for o := range out[x] {
					rec <- o
					wg.Done()
				}
			}()
		}
		wg.Add(1)
		pkg := []Pkg{}

		data := make([]byte, pkgSize, pkgSize)
		for i := 0; i < dataNum; i++ {
			for j := 0; j < pkgSize; j++ {
				data[j] = byte(rand.Int() % 64)
			}
			pkg = append(pkg, Pkg{id: i, data: data})
		}

		go func() {
			for i := 0; i < dataNum; i++ {
				in <- pkg[i]
			}
			wg.Done()
		}()
		wg.Wait()
		close(in)
		close(rec)
		for i := range rec {
			v := getPackageID(i)
			if _, ok := result[v]; ok {
				result[v]++
			} else {
				result[v] = 1
			}
		}

		for i := 0; i < dataNum; i++ {
			if result[i] != outChNum {
				b.Error("want", i, "receive", outChNum, "times not", result[i])
			}
		}
	}
}

//BenchmarkFanOut[outChNum]_[dataNum}_[bufNum]_[dataSize]
func BenchmarkFanOut6_100_100_10(b *testing.B) {
	benchmarkFanOut(6, 100, 100, 10, b)
}

func BenchmarkFanOut6_100_100_100(b *testing.B) {
	benchmarkFanOut(6, 100, 100, 100, b)
}

func BenchmarkFanOut6_100_100_1000(b *testing.B) {
	benchmarkFanOut(6, 100, 100, 1000, b)
}
