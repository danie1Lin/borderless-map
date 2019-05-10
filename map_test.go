package borderlessmap

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

var ele [][][]Element = [][][]Element{
	[][]Element{
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
	},
	[][]Element{
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
	},

	[][]Element{
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
		[]Element{1, 1, 1, 1, 1, 1, 1, 1},
	},
	[][]Element{
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
	},
	[][]Element{
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
	},
	[][]Element{
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
	},
	[][]Element{
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
		[]Element{0, 0, 0, 0, 0, 0, 0, 0},
	},
	[][]Element{
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
		[]Element{5, 5, 5, 5, 5, 5, 5, 5},
	},
}

func TestNewSpaceElement(t *testing.T) {
	se := NewSpaceElement(2, 2, 4, 1, 3, 2, 4, 5, 7, 6, 8, 9, 11, 10, 12, 13, 15, 14, 16)
	var expect SpaceElement = [][][]Element{
		[][]Element{
			[]Element{1, 5, 9, 13},
			[]Element{2, 6, 10, 14},
		},
		[][]Element{
			[]Element{3, 7, 11, 15},
			[]Element{4, 8, 12, 16},
		},
	}
	if !reflect.DeepEqual(se, expect) {
		t.Fatalf("want %s but %s", expect, se)
	}
}

func TestSpaceElement_Copy(t *testing.T) {

	e := SpaceElement(ele)
	dis := e.Copy()
	check := func(x, y, z int, value Element) bool {
		return (dis[x][y][z] == e[x][y][z]) && (&e[x][y][z] != &dis[x][y][z])
	}
	if x, y, z, ok := e.Do(check); !ok {
		t.Fatal(x, y, z, ": ", dis[x][y][z], &dis[x][y][z], e[x][y][z], &e[x][y][z])
	}
}

func TestSpaceElement_Set(t *testing.T) {
	copy := SpaceElement(ele).Copy()
	copy.Set(7, 7, 7, [][][]Element{})
}

func TestGetNotify(t *testing.T) {
	map1 := &SubmapBase{}
	map1.init(0, 2, ele)
	fmt.Print(map1.Data)
	obj1 := &ObjectBase{
		posInSubmap: []int{0, 0, 0},
	}
	map1.AddObj(obj1)
	fmt.Print("map size", map1.size)
	go map1.CheckBandArea()
	signals, _ := map1.ObjectSignal()
	wg := sync.WaitGroup{}
	fmt.Println("start")
	for i := 0; i < 6; i++ {
		wg.Add(1)
		id := i
		sig := signals[i]

		go func() {
			select {
			case obj := <-sig:
				t.Logf("chan %d get : %v\n", id, obj)
			case <-time.After(time.Millisecond):
				t.Logf("chan %d get nothing\n", id)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestMap_Split(t *testing.T) {
	data := &mapdata{
		data: ele,
	}
	w_map := &Map{
		MapData: data,
	}
	//fmt.Println(w_map.MapData.(*mapdata).data)
	w_map.Split(3, 3, 3)

	out := w_map.Submaps[2][2][0].(*SubmapBase).Data.String()
	matchStr := "" +
		"size:2X2X3:\n" +
		"z:0----2X2\n" +
		"0 5 \n" +
		"0 5 \n" +
		"z:1----2X2\n" +
		"0 5 \n" +
		"0 5 \n" +
		"z:2----2X2\n" +
		"0 5 \n" +
		"0 5 \n"
	if out != matchStr {
		t.Fatalf("\nwant : \n%s\ngot : \n%s", matchStr, out)
	}
	matchMap := SpaceElement{
		[][]Element{
			[]Element{0, 0, 0},
			[]Element{0, 0, 0},
		},
		[][]Element{
			[]Element{5, 5, 5},
			[]Element{5, 5, 5},
		},
	}

	if !reflect.DeepEqual(matchMap, w_map.Submaps[2][2][0].(*SubmapBase).Data) {
		t.Fatalf("want : \n%v\ngot : \n%v", matchMap, w_map.Submaps[2][2][0].(*SubmapBase).Data)
	}
}
