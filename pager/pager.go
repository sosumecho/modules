package pager

import (
	"entgo.io/ent/dialect/sql"
	"fmt"
	"strings"
)

type Filter string

var filters = make(map[FilterType]func(col string, value interface{}) func(s *sql.Selector))

func init() {
	filters[EQ] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(col), value))
		}
	}

	filters[GT] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			s.Where(sql.GT(s.C(col), value))
		}
	}

	filters[GTE] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(col), value))
		}
	}

	filters[LT] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			s.Where(sql.LT(s.C(col), value))
		}
	}

	filters[LTE] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(col), value))
		}
	}

	filters[Contains] = func(col string, value interface{}) func(s *sql.Selector) {
		return func(s *sql.Selector) {
			if rs, ok := value.(string); ok {
				s.Where(sql.Contains(s.C(col), rs))
				return
			}

			if rs, ok := value.(interface{ String() string }); ok {
				s.Where(sql.Contains(s.C(col), rs.String()))
				return
			}

			switch value.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				s.Where(sql.Contains(s.C(col), fmt.Sprintf("%d", value)))
				return
			}
		}
	}
}

func interfaceToString(v interface{}) string {
	if rs, ok := v.(string); ok {
		return rs
	}

	if rs, ok := v.(interface{ String() string }); ok {
		return rs.String()
	}

	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	}

	return ""
}

func (filter Filter) Value() string {
	var v = strings.Split(string(filter), ":")
	if len(v) == 1 {
		return v[0]
	}

	return v[1]
}

func (filter Filter) Type() FilterType {
	var v = strings.Split(string(filter), ":")
	if len(v) == 1 {
		return EQ
	}

	return parseFilterType(v[0])
}

func (filter Filter) Predicate(field string) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		for _, condition := range strings.Split(string(filter), ",") {
			var predicate = strings.SplitN(condition, ":", 2)
			if len(predicate) == 1 {
				EQ.Filter(field, predicate[0])(s)
			} else {
				parseFilterType(predicate[0]).Filter(field, predicate[1])(s)
			}
		}
	}
}

type FilterType uint

const (
	EQ FilterType = iota
	GT
	GTE
	LT
	LTE
	Contains
)

type Pager struct {
	Rows  int    `json:"row" form:"row"`
	Page  int    `json:"page" form:"page"`
	Sort  string `json:"sort" form:"sort"` // -created_at, +updated_at
	Edges string `json:"edges" form:"edges"`
}

func (pager *Pager) Build() {
	if pager.Rows < 1 {
		pager.Rows = 12
	}

	if pager.Rows > 500 {
		pager.Rows = 500
	}

	if pager.Page < 1 {
		pager.Page = 1
	}
}

func (pager *Pager) WithEdge(edge string) bool {
	for _, e := range strings.Split(pager.Edges, ",") {
		if e == edge {
			return true
		}
	}

	return false
}

func (filterType FilterType) Filter(col string, value interface{}) func(s *sql.Selector) {
	if interfaceToString(value) == "" {
		return func(s *sql.Selector) {

		}
	}
	if filter, ok := filters[filterType]; ok {
		return filter(col, value)
	}

	return func(s *sql.Selector) {

	}
}

func parseFilterType(v string) FilterType {
	switch v {
	case "eq":
		return EQ
	case "gt":
		return GT
	case "gte":
		return GTE
	case "lt":
		return LT
	case "lte":
		return LTE
	case "contains":
		return Contains
	}
	return EQ
}
