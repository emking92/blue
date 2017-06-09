package entities

type FilterGroup []Filter

type Filter struct {
	Signal Signal `json:"signal"`
	Count  int64  `json:"count"`
	Index  int    `json:"index"`
}

type Signal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

type Operation string

const (
	OperationAddition       = Operation("+")
	OperationSubtraction    = Operation("-")
	OperationMultiplication = Operation("*")
	OperationDivision       = Operation("/")
	OperationModulo         = Operation("%")
	OperationBitwiseAnd     = Operation("&")
	OperationBitwiseOr      = Operation("|")
	OperationBitwiseXor     = Operation("^")
)

type Comparator string

const (
	ComparatorEqual            = Comparator("=")
	ComparatorLessThan         = Comparator("<")
	ComparatorGreaterThan      = Comparator(">")
	ComparatorLessThanEqual    = Comparator("<=")
	ComparatorGreaterThanEqual = Comparator(">=")
	ComparatorNotEqual         = Comparator("?")
)

func SignalItem(item string) Signal {
	return Signal{Type: "item", Name: item}
}

var (
	SignalVirtualEverything = Signal{Type: "virtual", Name: "signal-everything"}
	SignalVirtualAnything   = Signal{Type: "virtual", Name: "signal-anything"}
	SignalVirtualEach       = Signal{Type: "virtual", Name: "signal-each"}
)

func NewFilterGroup(signalCountPairs ...interface{}) FilterGroup {
	group := []Filter{}

	for i := 0; i < len(signalCountPairs)/2; i++ {
		cntInterface := signalCountPairs[2*i+1]
		var cnt int64
		switch cntInterface.(type) {
		case int:
			cnt = int64(cntInterface.(int))
		case byte:
			cnt = int64(cntInterface.(byte))
		default:
			cnt = cntInterface.(int64)
		}

		if cnt == 0 {
			continue
		}

		group = append(group, Filter{
			Signal: signalCountPairs[2*i].(Signal),
			Count:  cnt,
			Index:  i + 1,
		})
	}

	return group
}
