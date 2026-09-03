package mapUtil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Tomatosky/jo-util/sliceUtil"
	"github.com/Tomatosky/jo-util/strUtil"
)

const (
	Wildcard  = "*"
	KeySepStr = "."
)

// DeepGet value by key path. eg "top" "top.sub"
func DeepGet(mp map[string]any, path string) (val any) {
	val, _ = GetByPath(mp, path)
	return
}

// GetByPath get value by key path from a map(map[string]any). eg "top" "top.sub"
func GetByPath(mp map[string]any, path string) (val any, ok bool) {
	if len(path) == 0 {
		return mp, true
	}
	if val, ok = mp[path]; ok {
		return val, true
	}

	// no sub key
	if len(mp) == 0 || strings.IndexByte(path, '.') < 1 {
		return nil, false
	}

	// key is path. eg: "top.sub"
	return GetByPathKeys(mp, strings.Split(path, "."))
}

// GetByPathKeys get value by path keys from a map(map[string]any). eg "top" "top.sub"
//
// Example:
//
//	mp := map[string]any{
//		"top": map[string]any{
//			"sub": "value",
//		},
//	}
//	val, ok := GetByPathKeys(mp, []string{"top", "sub"}) // return "value", true
func GetByPathKeys(mp map[string]any, keys []string) (val any, ok bool) {
	kl := len(keys)
	if kl == 0 {
		return mp, true
	}

	// find top item data use top key
	var item any
	topK := keys[0]
	if item, ok = mp[topK]; !ok {
		return
	}

	// find sub item data use sub key
	return getByPathKeys(item, keys[1:])
}

// indirect like reflect.Indirect(), but can also indirect reflect.Interface. otherwise, will return self
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}

func getByPathKeys(item any, keys []string) (val any, ok bool) {
	kl := len(keys)

	for i, k := range keys {
		switch tData := item.(type) {
		case map[string]string: // is string map
			if item, ok = tData[k]; !ok {
				return
			}
		case map[string]any: // is map(decode from toml/json/yaml)
			if item, ok = tData[k]; !ok {
				return
			}
		case map[any]any: // is map(decode from yaml.v2)
			if item, ok = tData[k]; !ok {
				return
			}
		case []map[string]any: // is an any-map slice
			if k == Wildcard {
				if kl == i+1 { // * is last key
					return tData, true
				}

				// * is not last key, find sub item data
				sl := make([]any, 0, len(tData))
				for _, v := range tData {
					if val, ok = getByPathKeys(v, keys[i+1:]); ok {
						sl = append(sl, val)
					}
				}

				if len(sl) > 0 {
					return sl, true
				}
				return nil, false
			}

			// k is index number
			idx, err := strconv.Atoi(k)
			if err != nil || idx < 0 || idx >= len(tData) {
				return nil, false
			}
			item = tData[idx]
		default:
			if k == Wildcard && kl == i+1 { // * is last key
				return tData, true
			}

			rv := reflect.ValueOf(tData)
			// check is slice
			if rv.Kind() == reflect.Slice {
				if k == Wildcard {
					// * is not last key, find sub item data
					sl := make([]any, 0, rv.Len())
					for si := 0; si < rv.Len(); si++ {
						el := indirect(rv.Index(si))
						if el.Kind() != reflect.Map {
							return nil, false
						}

						// el is map value.
						if val, ok = getByPathKeys(el.Interface(), keys[i+1:]); ok {
							sl = append(sl, val)
						}
					}

					if len(sl) > 0 {
						return sl, true
					}
					return nil, false
				}

				// check k is index number
				ii, err := strconv.Atoi(k)
				if err != nil || ii < 0 || ii >= rv.Len() {
					return nil, false
				}

				item = rv.Index(ii).Interface()
				continue
			}

			// as error
			return nil, false
		}

		// next is last key and it is *
		if kl == i+2 && keys[i+1] == Wildcard {
			return item, true
		}
	}

	return item, true
}

func DeepSet(mp *map[string]any, path string, val any) {
	_ = SetByPath(mp, path, val)
}

// SetByPath set sub-map value by key path.
// Supports dot syntax to set deep values.
//
// For example:
//
//	SetByPath("name.first", "Mat")
func SetByPath(mp *map[string]any, path string, val any) error {
	return SetByKeys(mp, strings.Split(path, KeySepStr), val)
}

// SetByKeys set sub-map value by path keys.
// Supports dot syntax to set deep values.
//
// For example:
//
//	SetByKeys([]string{"name", "first"}, "Mat")
func SetByKeys(mp *map[string]any, keys []string, val any) (err error) {
	if mp == nil {
		return fmt.Errorf("input parameter#mp must not be nil")
	}
	for _, key := range keys {
		pos := strings.IndexRune(key, '[')
		if pos < 1 || !strings.HasSuffix(key, "]") {
			continue
		}
		idx, parseErr := strconv.Atoi(key[pos+1 : len(key)-1])
		if parseErr == nil && idx < 0 {
			return fmt.Errorf("slice index must not be negative: %d", idx)
		}
	}

	kln := len(keys)
	if kln == 0 {
		return nil
	}

	mpv := *mp
	topK := keys[0]
	_, _, topHasIndex := parseArrKeyIndex(topK)
	if len(mpv) == 0 {
		if topHasIndex {
			rv := reflect.ValueOf(mp).Elem()
			return setMapByKeys(rv, keys, reflect.ValueOf(val))
		}
		*mp = MakeByKeys(keys, val)
		return nil
	}

	if topHasIndex {
		topK, _, _ = parseArrKeyIndex(topK)
	}
	if kln == 1 && !topHasIndex {
		mpv[topK] = val
		return nil
	}

	if _, ok := mpv[topK]; !ok {
		if topHasIndex {
			rv := reflect.ValueOf(mp).Elem()
			return setMapByKeys(rv, keys, reflect.ValueOf(val))
		}
		mpv[topK] = MakeByKeys(keys[1:], val)
		return nil
	}

	rv := reflect.ValueOf(mp).Elem()
	return setMapByKeys(rv, keys, reflect.ValueOf(val))
}

func setMapByKeys(rv reflect.Value, keys []string, nv reflect.Value) (err error) {
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("input parameter#rv must be a Map, but was %s", rv.Kind())
	}

	// If the map is nil, make a new map
	if rv.IsNil() {
		mapType := reflect.MapOf(rv.Type().Key(), rv.Type().Elem())
		rv.Set(reflect.MakeMap(mapType))
	}

	var ok bool
	maxI := len(keys) - 1
	for i, key := range keys {
		idx := -1
		isMap := rv.Kind() == reflect.Map
		isSlice := rv.Kind() == reflect.Slice
		isLast := i == len(keys)-1

		// slice index key must be ended on the keys.
		// eg: "top.arr[2]" -> "arr[2]"
		if pos := strings.IndexRune(key, '['); pos > 0 {
			var realKey string
			if realKey, idx, ok = parseArrKeyIndex(key); ok {
				// update value
				key = realKey
				if !isMap {
					err = fmt.Errorf(
						"current value#%s type is %s, cannot get sub-value by key: %s",
						strings.Join(keys[i:], "."),
						rv.Kind(),
						key,
					)
					break
				}

				rftK, valueErr := reflectValueForType(reflect.ValueOf(key), rv.Type().Key())
				if valueErr != nil {
					return valueErr
				}
				tmpV := rv.MapIndex(rftK)
				if !tmpV.IsValid() {
					if isLast {
						elemType := reflect.TypeOf((*any)(nil)).Elem()
						if nv.IsValid() {
							elemType = nv.Type()
						}
						sliVal := reflect.MakeSlice(reflect.SliceOf(elemType), idx+1, idx+1)
						if err = setReflectValue(sliVal.Index(idx), nv); err != nil {
							return err
						}
						if err = setReflectMapValue(rv, rftK, sliVal); err != nil {
							return err
						}
					} else {
						// deep make map by keys
						newVal := MakeByKeys(keys[i+1:], reflectInterface(nv))
						mpVal := reflect.ValueOf(newVal)

						sliVal := reflect.MakeSlice(reflect.SliceOf(mpVal.Type()), idx+1, idx+1)
						sliVal.Index(idx).Set(mpVal)

						if err = setReflectMapValue(rv, rftK, sliVal); err != nil {
							return err
						}
					}
					break
				}

				// get real type: any -> map
				if tmpV.Kind() == reflect.Interface {
					tmpV = tmpV.Elem()
				}

				if tmpV.Kind() != reflect.Slice {
					err = fmt.Errorf(
						"current value#%s type is %s, cannot set sub by index: %d",
						strings.Join(keys[i:], "."),
						tmpV.Kind(),
						idx,
					)
					break
				}

				wantLen := idx + 1
				sliLen := tmpV.Len()
				elemTyp := tmpV.Type().Elem()

				if wantLen > sliLen {
					newAdd := reflect.MakeSlice(tmpV.Type(), 0, wantLen-sliLen)
					for m := 0; m < wantLen-sliLen; m++ {
						newAdd = reflect.Append(newAdd, reflect.New(elemTyp).Elem())
					}

					tmpV = reflect.AppendSlice(tmpV, newAdd)
				}

				if !isLast {
					if elemTyp.Kind() == reflect.Map {
						err = setMapByKeys(tmpV.Index(idx), keys[i+1:], nv)
						if err != nil {
							return err
						}

						// tmpV.Index(idx).Set(elemV)
						if err = setReflectMapValue(rv, rftK, tmpV); err != nil {
							return err
						}
					} else {
						err = fmt.Errorf(
							"key %s[%d] elem must be map for set sub-value by remain path: %s",
							key,
							idx,
							strings.Join(keys[i:], "."),
						)
					}
				} else {
					// last - set value
					if err = setReflectValue(tmpV.Index(idx), nv); err != nil {
						return err
					}
					if err = setReflectMapValue(rv, rftK, tmpV); err != nil {
						return err
					}
				}
				break
			}
		}

		// set value on last key
		if isLast {
			if isMap {
				if rv.IsNil() {
					if !rv.CanSet() {
						return fmt.Errorf("cannot initialize map of type %s", rv.Type())
					}
					rv.Set(reflect.MakeMap(rv.Type()))
				}
				rftK, valueErr := reflectValueForType(reflect.ValueOf(key), rv.Type().Key())
				if valueErr != nil {
					return valueErr
				}
				if err = setReflectMapValue(rv, rftK, nv); err != nil {
					return err
				}
				break
			}

			if isSlice {
				// key is slice index
				if strUtil.IsInt(key) {
					idx, _ = strconv.Atoi(key)
				}

				if idx > -1 {
					wantLen := idx + 1
					sliLen := rv.Len()

					if wantLen > sliLen {
						elemTyp := rv.Type().Elem()
						newAdd := reflect.MakeSlice(rv.Type(), 0, wantLen-sliLen)

						for m := 0; m < wantLen-sliLen; m++ {
							newAdd = reflect.Append(newAdd, reflect.New(elemTyp).Elem())
						}

						if !rv.CanAddr() {
							err = fmt.Errorf("cannot set value to a cannot addr slice, key: %s", key)
							break
						}

						rv.Set(reflect.AppendSlice(rv, newAdd))
					}

					if err = setReflectValue(rv.Index(idx), nv); err != nil {
						return err
					}
				} else {
					err = fmt.Errorf("cannot set slice value by named key %q", key)
				}
			} else {
				err = fmt.Errorf(
					"cannot set sub-value for type %q(path %q, key %q)",
					rv.Kind(),
					strings.Join(keys[:i], "."),
					key,
				)
			}

			break
		}

		if isMap {
			rftK, valueErr := reflectValueForType(reflect.ValueOf(key), rv.Type().Key())
			if valueErr != nil {
				return valueErr
			}
			if tmpV := rv.MapIndex(rftK); tmpV.IsValid() {
				var isPtr bool
				// get real type: any -> map
				tmpV, isPtr = getRealVal(tmpV)
				if tmpV.Kind() == reflect.Map {
					rv = tmpV
					continue
				}

				// sub is slice and is not ptr
				if tmpV.Kind() == reflect.Slice {
					if isPtr {
						rv = tmpV
						continue // to (E)
					}

					// next key is index number.
					nxtKey := keys[i+1]
					if strUtil.IsInt(nxtKey) {
						idx, _ = strconv.Atoi(nxtKey)
						if idx < 0 {
							return fmt.Errorf("slice index must not be negative: %d", idx)
						}
						sliLen := tmpV.Len()
						wantLen := idx + 1

						if wantLen > sliLen {
							elemTyp := tmpV.Type().Elem()
							newAdd := reflect.MakeSlice(tmpV.Type(), 0, wantLen-sliLen)
							for m := 0; m < wantLen-sliLen; m++ {
								newAdd = reflect.Append(newAdd, reflect.New(elemTyp).Elem())
							}

							tmpV = reflect.AppendSlice(tmpV, newAdd)
						}

						// rv = tmpV.Index(idx) // TODO
						if i+1 == maxI {
							if err = setReflectValue(tmpV.Index(idx), nv); err != nil {
								return err
							}
						} else {
							err = setMapByKeys(tmpV.Index(idx), keys[i+2:], nv)
							if err != nil {
								return err
							}
						}

						if err = setReflectMapValue(rv, rftK, tmpV); err != nil {
							return err
						}
					} else {
						err = fmt.Errorf("cannot set slice value by named key %s(parent: %s)", nxtKey, key)
					}
				} else {
					err = fmt.Errorf(
						"map item type is %s(path:%q), cannot set sub-value by path %q",
						tmpV.Kind(),
						strings.Join(keys[0:i+1], "."),
						strings.Join(keys[i+1:], "."),
					)
				}
			} else {
				// deep make map by keys
				newVal := MakeByKeys(keys[i+1:], reflectInterface(nv))
				if err = setReflectMapValue(rv, rftK, reflect.ValueOf(newVal)); err != nil {
					return err
				}
			}

			break
		} else if isSlice && strUtil.IsInt(key) { // (E). slice from ptr slice
			idx, _ = strconv.Atoi(key)
			if idx < 0 {
				return fmt.Errorf("slice index must not be negative: %d", idx)
			}
			sliLen := rv.Len()
			wantLen := idx + 1

			if wantLen > sliLen {
				elemTyp := rv.Type().Elem()
				newAdd := reflect.MakeSlice(rv.Type(), 0, wantLen-sliLen)
				for m := 0; m < wantLen-sliLen; m++ {
					newAdd = reflect.Append(newAdd, reflect.New(elemTyp).Elem())
				}

				if !rv.CanSet() {
					return fmt.Errorf("cannot expand slice of type %s", rv.Type())
				}
				rv.Set(reflect.AppendSlice(rv, newAdd))
			}

			rv = rv.Index(idx)
		} else {
			err = fmt.Errorf(
				"map item type is %s, cannot set sub-value by path %q",
				rv.Kind(),
				strings.Join(keys[i:], "."),
			)
		}
	}
	return
}

func setReflectValue(dst, value reflect.Value) error {
	if !dst.CanSet() {
		return fmt.Errorf("cannot set value of type %s", dst.Type())
	}
	prepared, err := reflectValueForType(value, dst.Type())
	if err != nil {
		return err
	}
	dst.Set(prepared)
	return nil
}

func setReflectMapValue(dst, key, value reflect.Value) error {
	if dst.IsNil() {
		return fmt.Errorf("cannot set entry on nil map of type %s", dst.Type())
	}
	prepared, err := reflectValueForType(value, dst.Type().Elem())
	if err != nil {
		return err
	}
	dst.SetMapIndex(key, prepared)
	return nil
}

func reflectValueForType(value reflect.Value, target reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		switch target.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			return reflect.Zero(target), nil
		default:
			return reflect.Value{}, fmt.Errorf("cannot assign nil to %s", target)
		}
	}
	if !value.Type().AssignableTo(target) {
		return reflect.Value{}, fmt.Errorf("cannot assign %s to %s", value.Type(), target)
	}
	return value, nil
}

func reflectInterface(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

func MakeByKeys(keys []string, val any) (mp map[string]any) {
	pathKeys := append([]string(nil), keys...)
	size := len(pathKeys)
	if size == 0 {
		return map[string]any{}
	}

	// if last key contains slice index, make slice wrap the val
	lastKey := pathKeys[size-1]
	if newK, idx, ok := parseArrKeyIndex(lastKey); ok {
		elemType := reflect.TypeOf((*any)(nil)).Elem()
		if val != nil {
			elemType = reflect.TypeOf(val)
		}
		sliTyp := reflect.SliceOf(elemType)
		sliVal := reflect.MakeSlice(sliTyp, idx+1, idx+1)
		if val != nil {
			sliVal.Index(idx).Set(reflect.ValueOf(val))
		}

		// update val and last key
		val = sliVal.Interface()
		pathKeys[size-1] = newK
	}

	if size == 1 {
		return map[string]any{pathKeys[0]: val}
	}

	// multi nodes
	sliceUtil.Reverse(pathKeys)
	for _, p := range pathKeys {
		if mp == nil {
			mp = map[string]any{p: val}
		} else {
			mp = map[string]any{p: mp}
		}
	}
	return
}

// "arr[2]" => "arr", 2, true
func parseArrKeyIndex(key string) (string, int, bool) {
	pos := strings.IndexRune(key, '[')
	if pos < 1 || !strings.HasSuffix(key, "]") {
		return key, 0, false
	}

	var idx int
	var err error

	idxStr := key[pos+1 : len(key)-1]
	if idxStr != "" {
		idx, err = strconv.Atoi(idxStr)
		if err != nil {
			return key, 0, false
		}
	}

	key = key[:pos]
	return key, idx, true
}

func getRealVal(rv reflect.Value) (reflect.Value, bool) {
	// get real type: any -> map
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	isPtr := false
	if rv.Kind() == reflect.Ptr {
		isPtr = true
		rv = rv.Elem()
	}

	return rv, isPtr
}
