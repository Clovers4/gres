package container

type CSet struct {
	cm *CMap
}

func NewCSet() *CSet {
	return &CSet{
		NewCMap(),
	}
}

func (cs *CSet) Add(v interface{}) {

}
