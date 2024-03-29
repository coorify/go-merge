package merge

import (
	"fmt"
	"reflect"
)

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

func deepMerge(dv, sv reflect.Value) error {
	dt := dv.Type()
	st := sv.Type()

	if dv.Kind() != reflect.Struct {
		return fmt.Errorf("invalid field: %s, should be struct,map", dt.Name())
	}

	if sv.Kind() != reflect.Struct {
		return fmt.Errorf("invalid field: %s, should be struct,map", st.Name())
	}

	fds := map[string]struct{}{}

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

		fds[sft.Name] = struct{}{}
		dfv := dv.FieldByName(sft.Name)
		if sfv.Kind() == reflect.Struct {
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

	fromMethod(dv, sv, &fds)
	if sv.CanAddr() {
		fromMethod(dv, sv.Addr(), &fds)
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
