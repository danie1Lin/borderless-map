package borderlessmap

//close inCh will terminate FanOut gorutine
type FanOut struct {
	inCh   chan interface{}
	outChs []chan interface{}
}

func (f *FanOut) run() {
	defer func() {
		for i := range f.outChs {
			close(f.outChs[i])
		}
		f.outChs = make([]chan interface{}, 0)
	}()
	for value := range f.inCh {
		for i := range f.outChs {
			select {
			case f.outChs[i] <- value:
				//log.Printf("channel %d is send %d", i, value)
			}
		}
	}
}

func (f *FanOut) Register(buf int) <-chan interface{} {
	outChan := make(chan interface{}, buf)
	f.outChs = append(f.outChs, outChan)
	return outChan
}

func (f *FanOut) Init(buf int) chan<- interface{} {
	f.inCh = make(chan interface{}, buf)
	go f.run()
	return f.inCh
}

type FanIn struct {
	inChs []chan interface{}
	outCh chan interface{}
	done  chan interface{}
}

func (f *FanIn) Register(buf int) chan<- interface{} {
	inCh := make(chan interface{}, buf)
	f.inChs = append(f.inChs, inCh)
	return inCh
}

func (f *FanIn) Init(buf int) (chan interface{}, chan interface{}) {
	f.outCh = make(chan interface{}, buf)
	f.inChs = make([]chan interface{}, 0)
	f.done = make(chan interface{}, 0)
	go f.run()
	return f.outCh, f.done
}

func (f *FanIn) run() {
	defer func() {
		f.inChs = make([]chan interface{}, 0)
		close(f.outCh)
	}()
	for {
		select {
		case <-f.done:
			return
		default:
			for i := range f.inChs {
				select {
				case value := <-f.inChs[i]:
					f.outCh <- value
				default:
				}
			}
		}
	}
}

type FanInOut struct {
	FanIn
	FanOut
	midInCh  chan interface{}
	midOutCh chan interface{}
	Do       func(value interface{}, midOutCh chan interface{})
	done     chan interface{}
}

func (f *FanInOut) RegisterOut(buf int) <-chan interface{} {
	return f.FanOut.Register(buf)
}

func (f *FanInOut) RegisterIn(buf int) chan<- interface{} {
	return f.FanIn.Register(buf)
}

/*
func (f *FanInOut) Init(inBuf, outBuf int) (chan interface{}, chan interface{}) {
	f.FanIn.Init(inBuf)
	f.FanOut.Init(outBuf)
	f.run()
	return
}

func (f *FanInOut) run() {

}
*/
