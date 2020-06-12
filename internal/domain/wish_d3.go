// Code generated by d3. DO NOT EDIT.

package domain

import "github.com/godzie44/d3/orm/entity"
import "time"
import "database/sql/driver"
import "fmt"

func (w *Wish) D3Token() entity.MetaToken {
	return entity.MetaToken{
		Tpl:       (*Wish)(nil),
		TableName: "lw_wish",
		Tools: entity.InternalTools{
			ExtractField:  w.__d3_makeFieldExtractor(),
			SetFieldVal:   w.__d3_makeFieldSetter(),
			CompareFields: w.__d3_makeComparator(),
			NewInstance:   w.__d3_makeInstantiator(),
			Copy:          w.__d3_makeCopier(),
		},
	}
}

func (w *Wish) __d3_makeFieldExtractor() entity.FieldExtractor {
	return func(s interface{}, name string) (interface{}, error) {
		sTyped, ok := s.(*Wish)
		if !ok {
			return nil, fmt.Errorf("invalid entity type")
		}

		switch name {

		case "id":
			return sTyped.id, nil

		case "content":
			return sTyped.content, nil

		case "createAt":
			return sTyped.createAt, nil

		default:
			return nil, fmt.Errorf("field %s not found", name)
		}
	}
}

func (w *Wish) __d3_makeInstantiator() entity.Instantiator {
	return func() interface{} {
		return &Wish{}
	}
}

func (w *Wish) __d3_makeFieldSetter() entity.FieldSetter {
	return func(s interface{}, name string, val interface{}) error {
		eTyped, ok := s.(*Wish)
		if !ok {
			return fmt.Errorf("invalid entity type")
		}

		switch name {
		case "content":
			eTyped.content = val.(string)
			return nil
		case "createAt":
			eTyped.createAt = val.(time.Time)
			return nil

		case "id":
			if valuer, isValuer := val.(driver.Valuer); isValuer {
				v, err := valuer.Value()
				if err != nil {
					return eTyped.id.Scan(nil)
				}
				return eTyped.id.Scan(v)
			}
			return eTyped.id.Scan(val)
		default:
			return fmt.Errorf("field %s not found", name)
		}
	}
}

func (w *Wish) __d3_makeCopier() entity.Copier {
	return func(src interface{}) interface{} {
		srcTyped, ok := src.(*Wish)
		if !ok {
			return fmt.Errorf("invalid entity type")
		}

		copy := &Wish{}

		copy.id = srcTyped.id
		copy.content = srcTyped.content
		copy.createAt = srcTyped.createAt

		return copy
	}
}

func (w *Wish) __d3_makeComparator() entity.FieldComparator {
	return func(e1, e2 interface{}, fName string) bool {
		if e1 == nil || e2 == nil {
			return e1 == e2
		}

		e1Typed, ok := e1.(*Wish)
		if !ok {
			return false
		}
		e2Typed, ok := e2.(*Wish)
		if !ok {
			return false
		}

		switch fName {

		case "id":
			return e1Typed.id == e2Typed.id
		case "content":
			return e1Typed.content == e2Typed.content
		case "createAt":
			return e1Typed.createAt == e2Typed.createAt
		default:
			return false
		}
	}
}
