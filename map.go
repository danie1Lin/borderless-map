package borderlessmap

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	X_bit = iota
	Y_bit
	Z_bit
)

const (
	N_left = iota
	N_right
	N_front
	N_back
	N_up
	N_down
)

type Element int

const (
	Land Element = iota
	Water
	Obstacle
	Movable
	top
)

func MakeSpaceElement(x, y, z int) SpaceElement {
	s := make([][][]Element, x)
	for x_i := range s {
		s[x_i] = make([][]Element, y)
		for y_i := range s[x_i] {
			s[x_i][y_i] = make([]Element, z)
		}
	}
	return s
}

func NewSpaceElement(lx, ly, lz int, e ...Element) SpaceElement {
	l := len(e)
	if l != lx*ly*lz {
		return nil
	}
	res := MakeSpaceElement(lx, ly, lz)
	asign := func(x, y, z int, v Element) bool {
		e_i := z*lx*ly + y*lx + x
		res[x][y][z] = e[e_i]
		return true
	}
	res.Do(asign)
	return res
}

type SpaceElement [][][]Element

func (s SpaceElement) Do(f func(x, y, z int, value Element) bool) (int, int, int, bool) {
	for x := range s {
		for y := range s[x] {
			for z := range s[x][y] {
				if !f(x, y, z, s[x][y][z]) {
					return x, y, z, false
				}
			}
		}
	}
	return 0, 0, 0, true
}
func (s SpaceElement) Get(x, y, z int) Element {
	return s[x][y][z]
}

func (s SpaceElement) Set(x, y, z int, se SpaceElement) SpaceElement {
	l := se.GetSize()
	for i_x, x := range s[x:l[0]] {
		for i_y, y := range x[y:l[1]] {
			for i_z := range y {
				y[i_z+l[2]] = se[i_x][i_y][i_z]
			}
		}
	}
	return s
}

func (s SpaceElement) Copy() SpaceElement {
	size := s.GetSize()
	dis := make([][][]Element, size[0])
	for x := range dis {
		dis[x] = make([][]Element, size[1])
		for y := range dis[x] {
			dis[x][y] = make([]Element, size[2])
		}
	}
	copyf := func(x, y, z int, value Element) bool {
		dis[x][y][z] = value
		return true
	}

	s.Do(copyf)

	return dis
}

func (s SpaceElement) GetRange(xl, xh, yl, yh, zl, zh int) SpaceElement {
	var block SpaceElement = make([][][]Element, xh-xl)
	for i_x, x := range s[xl:xh] {
		block[i_x] = make([][]Element, yh-yl)
		for i_y, y := range x[yl:yh] {
			block[i_x][i_y] = y[zl:zh]
		}
	}
	return block
}

func (s SpaceElement) GetSize() []int {
	return []int{len(s), len(s[0]), len(s[0][0])}
}

func (s SpaceElement) String() string {
	out := bytes.NewBuffer(make([]byte, 0))
	l := s.GetSize()
	fmt.Fprintf(out, "size:%dX%dX%d:\n", l[0], l[1], l[2])
	for z := 0; z < l[2]; z++ {
		fmt.Fprintf(out, "z:%d----%dX%d\n", z, l[0], l[1])
		for y := 0; y < l[1]; y++ {
			for x := 0; x < l[0]; x++ {
				fmt.Fprintf(out, "%d ", int(s[x][y][z]))
			}
			fmt.Fprintf(out, "\n")
		}
	}
	return out.String()
}

//border orient index
/*
const (
	XpYpZp = b000 //X+Y+Z+
	XnYpZp = b001
	XpYnZp = b010
	XnYnZp = b011
	XpYpZn = b100
	XnYpZn = b101
	XpYnZn = b110
	XnYnZn = b111
)
*/

func SyncObject(obj Object, from, to Submap) {

}

func ChangeObjectOwner(obj Object, form, to Submap) {
	//if failed change back
}

type MapId string

type MapData interface {
	GetSize() ([]int, error)
	GetBlock(xl, xh, yl, yh, zl, zh int) (SpaceElement, error)
}

type mapdata struct {
	data SpaceElement
	size []int
}

func (md *mapdata) GetSize() ([]int, error) {
	if md.size == nil {
		md.size = make([]int, 3)
		md.size[0] = len(md.data)
		md.size[1] = len(md.data[0])
		md.size[2] = len(md.data[0][0])
		for xid, x_data := range md.data {
			if len(x_data) != md.size[1] {
				return nil, fmt.Errorf("missing row in  x axis %d", xid)
			}
			for yid, y_data := range x_data {
				if len(y_data) != md.size[2] {
					return nil, fmt.Errorf("data missing in  x:%d y:%d", xid, yid)
				}
			}
		}
	}
	return md.size, nil
}

func (md *mapdata) GetBlock(xl, xh, yl, yh, zl, zh int) (SpaceElement, error) {
	if xl < 0 || yl < 0 || zl < 0 || xh > md.size[0] || yh > md.size[2] || zh > md.size[2] {
		return nil, fmt.Errorf("error index")
	}
	var block SpaceElement = md.data.GetRange(xl, xh, yl, yh, zl, zh)
	return block, nil
}

type Map struct {
	MapData
	Submaps [][][]Submap
}

func (m *Map) Split(split ...int) error {
	mapData := m.MapData
	size, err := mapData.GetSize()
	if err != nil {
		log.Warn(err)
		return err
	}
	lastBlock := []int{0, 0, 0}
	otherBlock := []int{0, 0, 0}
	l := len(split)
	if l != len(size) {
		if l == 0 {
			panic(fmt.Errorf("can't split into 0 splits"))
		} else if l == 1 {
			split = append(split, split[0], split[0])
		} else if l == 2 {
			split = []int{split[0], split[0], split[1]}
		} else {
			panic(fmt.Errorf("too many dimensions"))
		}
	}
	loop := []int{0, 0, 0}
	for i := range size {
		if split[i] <= 1 {
			loop[i] = 1
			break
		}
		lastBlock[i] = size[i] % (split[i])
		otherBlock[i] = (size[i] - lastBlock[i]) / (split[i] - 1) // size = (split-1) * otherblock + lastBlock

		if otherBlock[i] == 0 {
			panic(fmt.Errorf("split failed in axis %d", i))
		}
		loop[i] = split[i]
	}
	m.Submaps = make([][][]Submap, loop[0])
	for x_s := 0; x_s < loop[0]; x_s++ {
		m.Submaps[x_s] = make([][]Submap, loop[1])
		for y_s := 0; y_s < loop[1]; y_s++ {
			m.Submaps[x_s][y_s] = make([]Submap, loop[2])
			for z_s := 0; z_s < loop[2]; z_s++ {
				xh, yh, zh := otherBlock[0]*(x_s+1), otherBlock[1]*(y_s+1), otherBlock[2]*(z_s+1)
				if lastBlock[0] != 0 && x_s+1 == loop[0] {
					xh = otherBlock[0]*x_s + lastBlock[0]
				}
				if lastBlock[1] != 0 && y_s+1 == loop[1] {
					yh = otherBlock[1]*y_s + lastBlock[1]
				}
				if lastBlock[2] != 0 && z_s+1 == loop[2] {
					zh = otherBlock[2]*z_s + lastBlock[2]
				}
				block, _ := mapData.GetBlock(otherBlock[0]*x_s, xh, otherBlock[1]*y_s, yh, otherBlock[2]*z_s, zh)
				smb := &SubmapBase{}
				smb.init(0, 2, block)
				m.Submaps[x_s][y_s][z_s] = smb
				if err != nil {
					panic(fmt.Errorf("block %d %d %d split failed", x_s, y_s, z_s))
				}
			}
		}
	}
	return nil
}

type Submap interface {
	init(buf, band int, data SpaceElement) error
	ObjectSignal() ([]chan Object, error)
	AddObj(objs ...Object)
}

type SubmapBase struct {
	Id         MapId
	Neighbor   []*Submap //左右前後上下
	Border     [][]int   //border[0..7] = [x,y] the minimun cube hold
	BandUnit   int       //the Unit two submap overlaping
	Data       SpaceElement
	size       []int
	Objects    map[Object]struct{}
	ObjectSync []chan Object
	Scaler     float64 // unit: m/1
}

func (smb *SubmapBase) GetObjs() (map[Object]struct{}, error) {
	return smb.Objects, nil
}

func (smb *SubmapBase) SyncObj(obj Object) error {
	return nil
}

func (smb *SubmapBase) ObjectSignal() ([]chan Object, error) {
	return smb.ObjectSync, nil
}

func (smb *SubmapBase) AddObj(objs ...Object) {
	for _, obj := range objs {
		if _, ok := smb.Objects[obj]; !ok {
			smb.Objects[obj] = struct{}{}
		} else {
			log.Warn("already have this obj")
		}
	}
}

func (smb *SubmapBase) init(buf, band int, data SpaceElement) error {
	smb.BandUnit = band
	smb.Objects = make(map[Object]struct{})
	smb.ObjectSync = make([]chan Object, 6)
	for i := 0; i < 6; i++ {
		smb.ObjectSync[i] = make(chan Object, buf)
	}
	smb.Data = data
	smb.size = []int{len(data), len(data[0]), len(data[0][0])}
	return nil
}

func (smb *SubmapBase) CheckBandArea() error {
	for Obj := range smb.Objects {

		if obj, ok := Obj.(*ObjectBase); ok {
			//left
			dim := 0
			for _, i := range []int{0, 2, 4} {
				if obj.posInSubmap[dim] <= smb.BandUnit {
					log.Infof("border %d send notify\n", i)
					smb.ObjectSync[i] <- obj
				}
				dim++
			}
			dim = 0
			for _, i := range []int{1, 3, 5} {
				if obj.posInSubmap[dim] >= smb.size[dim]-smb.BandUnit {
					log.Infof("border %d send notify\n", i)
					smb.ObjectSync[i] <- obj
				}
				dim++
			}
		} else {
			err := fmt.Errorf("Not objectBase typ")
			log.Info(err)
			return err
		}
	}
	return nil
}

type Object interface {
}

type ObjectBase struct {
	Orient   []float64
	position []float64
	//In
	box         [][][]int
	posInSubmap []int     //
	AOI         [][][]int //AOI should <= BandUnit
}
