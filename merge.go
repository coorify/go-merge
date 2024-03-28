package merge

import (
	"fmt"
	"reflect"
)

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

		dft, has := dt.FieldByName(sft.Name)
		if has {
			fds[sft.Name] = struct{}{}
			dfv := dv.FieldByName(sft.Name)
			if sfv.Kind() == reflect.Struct {
				if err := deepMerge(dfv, sfv); err != nil {
					return err
				}
			} else if dft.Type == sft.Type {
				if !sfv.IsZero() {
					dfv.Set(sfv)
				}
			}
		}
	}

	for i := 0; i < sv.NumMethod(); i++ {
		sfmv := sv.Method(i)
		sfmt := st.Method(i)

		_, has := fds[sfmt.Name]
		if has {
			continue
		}

		dft, has := dt.FieldByName(sfmt.Name)
		if !has {
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

		if dft.Type == sfv.Type() {
			if !sfv.IsZero() {
				dfv.Set(sfv)
			}
		}
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
