package parser

import (
	"errors"
	"reflect"
)

var (
	BindPointerErr = errors.New("target value must be a pointer to struct")
	BindStructErr  = errors.New("target value must be a pointer to struct")
	BindFieldErr   = errors.New("target field must be a pointer to struct, slice of structs or plain struct")
)

type PegBindRef struct {
	Start, End int
}

func element(t reflect.Type) (reflect.Type, error) {
	base := t
	if t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		base = t.Elem()
	}
	if base.Kind() != reflect.Struct {
		return nil, BindFieldErr
	}
	return base, nil
}

func reference(v reflect.Value) (reflect.Value, error) {
	if v.Kind() == reflect.Pointer && v.Type().Elem().Kind() == reflect.Struct {
		return v, nil
	} else if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Struct {
		return v.Index(v.Len() - 1).Addr(), nil
	} else if v.Kind() == reflect.Struct {
		return v.Addr(), nil
	}
	return reflect.Value{}, BindFieldErr
}

func bindNode(n ParsingNode, fields map[string]reflect.Value) error {
	symbol := n.Symbol()
	field, fieldOk := fields[symbol]
	if fieldOk {
		t, err := element(field.Type())
		if err != nil {
			return err
		}
		if field.Kind() == reflect.Slice {
			field.Set(reflect.Append(field, reflect.Zero(t)))
		} else if field.Kind() == reflect.Pointer {
			field.Set(reflect.New(t))
		} else {
			field.Set(reflect.Zero(t))
		}
		r, err := reference(field)
		if err != nil {
			return err
		}
		return bind(n, r)
	}
	for _, c := range n.Children() {
		if err := bindNode(c, fields); err != nil {
			return err
		}
	}
	return nil
}

func bind(n ParsingNode, p reflect.Value) error {
	if p.Kind() != reflect.Pointer {
		return BindPointerErr
	}
	v := p.Elem()
	if v.Kind() != reflect.Struct {
		return BindStructErr
	}
	t := v.Type()
	_, startOk := t.FieldByName("Start")
	_, endOk := t.FieldByName("End")
	if startOk && endOk {
		start, end := n.Range()
		v.FieldByName("Start").Set(reflect.ValueOf(start))
		v.FieldByName("End").Set(reflect.ValueOf(end))
	}

	fields := make(map[string]reflect.Value, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		symbol := field.Tag.Get("peg")
		if symbol != "" {
			fields[symbol] = v.Field(i)
		}
	}
	return bindNode(n, fields)
}

func Bind(n ParsingNode, v any) error { return bind(n, reflect.ValueOf(v)) }
