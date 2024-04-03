package merge

import (
	"fmt"
	"reflect"
)

func i2v(v interface{}, kind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.Bool:
		return reflect.ValueOf(v.(bool))
	case reflect.Int:
		return reflect.ValueOf(v.(int))
	case reflect.Int8:
		return reflect.ValueOf(v.(int8))
	case reflect.Int16:
		return reflect.ValueOf(v.(int16))
	case reflect.Int32:
		return reflect.ValueOf(v.(int32))
	case reflect.Int64:
		return reflect.ValueOf(v.(int64))
	case reflect.Uint:
		return reflect.ValueOf(v.(uint))
	case reflect.Uint8:
		return reflect.ValueOf(v.(uint8))
	case reflect.Uint16:
		return reflect.ValueOf(v.(uint16))
	case reflect.Uint32:
		return reflect.ValueOf(v.(uint32))
	case reflect.Uint64:
		return reflect.ValueOf(v.(uint64))
	case reflect.Float32:
		return reflect.ValueOf(v.(float32))
	case reflect.Float64:
		return reflect.ValueOf(v.(float64))
	case reflect.String:
		return reflect.ValueOf(v.(string))
	}

	return reflect.Value{}
}

func fromMethod(dv, sv reflect.Value, fds *map[string]struct{}) error {
	dt := dv.Type()
	st := sv.Type()

	for i := 0; i < sv.NumMethod(); i++ {
		sfmv := sv.Method(i)
		sfmt := st.Method(i)

		if _, has := (*fds)[sfmt.Name]; has {
			continue
		}

		if _, has := dt.FieldByName(sfmt.Name); has {
			continue
		}

		var sfv reflect.Value
		dfv := dv.FieldByName(sfmt.Name)
		sfvs := sfmv.Call(make([]reflect.Value, 0))
		sfvl := len(sfvs)

		if sfvl == 1 {
			sfv = sfvs[0]
		} else if sfvl == 2 {
			sfv = sfvs[0]
			err := sfvs[1].Interface()
			if err != nil {
				return err.(error)
			}
		}

		if sfv.IsZero() {
			continue
		}

		(*fds)[sfmt.Name] = struct{}{}
		sfv = reflect.Indirect(sfv)
		if dfv.Type() == sfv.Type() {
			dfv.Set(sfv)
		}

		if sfv.CanAddr() && sfv.Addr().Type() == dfv.Type() {
			dfv.Set(sfv.Addr())
		}

		if dfv.CanAddr() && sfv.Type() == dfv.Addr().Type() {
			dfv.Elem().Set(sfv)
		}
	}

	return nil
}

func fromField(dv, sv reflect.Value, fds *map[string]struct{}) error {
	dt := dv.Type()
	st := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		sfv := sv.Field(i)
		sft := st.Field(i)

		if sfv.IsZero() {
			continue
		}

		_, has := dt.FieldByName(sft.Name)
		if !has {
			continue
		}

		(*fds)[sft.Name] = struct{}{}
		dfv := dv.FieldByName(sft.Name)
		if sfv.Kind() == reflect.Struct || sfv.Kind() == reflect.Map {
			if err := deepMerge(dfv, sfv); err != nil {
				return err
			}
		}

		sfv = reflect.Indirect(sfv)
		if dfv.Type() == sfv.Type() {
			dfv.Set(sfv)
		}

		if sfv.CanAddr() && sfv.Addr().Type() == dfv.Type() {
			dfv.Set(sfv.Addr())
		}

		if dfv.CanAddr() && sfv.Type() == dfv.Addr().Type() {
			dfv.Elem().Set(sfv)
		}
	}

	return nil
}

func fromMap(dv, sv reflect.Value, fds *map[string]struct{}) error {
	dt := dv.Type()

	ks := sv.MapKeys()
	for _, k := range ks {
		name, ok := k.Interface().(string)
		if !ok {
			continue
		}
		if _, has := (*fds)[name]; has {
			continue
		}

		_, has := dt.FieldByName(name)
		if !has {
			continue
		}

		sfv := sv.MapIndex(k)
		if sfv.IsZero() {
			continue
		}

		(*fds)[name] = struct{}{}
		dfv := dv.FieldByName(name)
		if sfv.Kind() == reflect.Interface {
			sfv = i2v(sfv.Interface(), dfv.Kind())
		}

		dfv.Set(sfv)
	}

	return nil
}

func deepMerge(dv, sv reflect.Value) error {
	dt := dv.Type()
	st := sv.Type()

	if dv.Kind() != reflect.Struct {
		return fmt.Errorf("invalid field: %s, should be struct", dt.Name())
	}

	if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Map {
		return fmt.Errorf("invalid field: %s, should be struct,map", st.Name())
	}

	fds := map[string]struct{}{}

	if sv.Kind() == reflect.Struct {
		fromField(dv, sv, &fds)
		fromMethod(dv, sv, &fds)
		if sv.CanAddr() {
			fromMethod(dv, sv.Addr(), &fds)
		}
	} else if sv.Kind() == reflect.Map {
		fromMap(dv, sv, &fds)
	}

	return nil
}

func Merge(dst interface{}, src ...interface{}) error {
	dv := reflect.Indirect(reflect.ValueOf(dst))
	for _, s := range src {
		sv := reflect.Indirect(reflect.ValueOf(s))
		if err := deepMerge(dv, sv); err != nil {
			return err
		}
	}
	return nil
}
