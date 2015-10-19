package metadata

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

type Encoder interface {
	Encode(w io.Writer, prefix string) error
}

func EncodeStringPrefix(e Encoder, prefix string) (string, error) {
	buf := new(bytes.Buffer)
	err := e.Encode(buf, prefix)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func EncodeString(e Encoder) (string, error) {
	return EncodeStringPrefix(e, "")
}

func EncodedFieldTag(t reflect.Type, name, lookup string) (string, error) {
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("field tags only work on structs")
	}
	field, ok := t.FieldByName(name)
	if !ok {
		return "", fmt.Errorf("unable fo find field: %s", strconv.Quote(name))
	}
	tags := strings.Split(field.Tag.Get(lookup), ",")
	if len(tags) > 0 {
		return tags[0], nil
	}
	return "", nil
}

type keys []reflect.Value

func (k keys) Len() int           { return len(k) }
func (k keys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k keys) Less(i, j int) bool { return k.get(i) < k.get(j) }
func (k keys) get(i int) string   { return k[i].String() }

func EncodeElement(w io.Writer, i interface{}, prefix string) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(i); err != nil {
		return err
	}
	for _, l := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if _, err := w.Write([]byte(prefix + l + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func EncodeEmptyMap(w io.Writer, t reflect.Type, label, prefix string) error {
	if t.Implements(reflect.TypeOf(new(Encoder)).Elem()) {
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#[%s.value]\n", label))); err != nil {
			return err
		}
		s, err := EncodeStringPrefix(reflect.New(t).Interface().(Encoder), prefix+"#\t")
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
	}
	return nil
}

func EncodeEmptySlice(w io.Writer, t reflect.Type, label, prefix string) error {
	if t.Implements(reflect.TypeOf(new(Encoder)).Elem()) {
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#[[%s]]\n", label))); err != nil {
			return err
		}
		s, err := EncodeStringPrefix(reflect.New(t).Interface().(Encoder), prefix+"#\t")
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
	}
	return nil
}
func EncodeSliceElement(w io.Writer, i interface{}, label, prefix string) error {
	if _, err := w.Write([]byte(prefix + fmt.Sprintf("[[%s]]\n", label))); err != nil {
		return err
	}
	if err := EncodeStruct(w, i, prefix+"\t"); err != nil {
		return err
	}
	return nil
}

func EncodeSlice(w io.Writer, v reflect.Value, label, prefix string) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return fmt.Errorf("encode slice only works on structs")
	}
	if v.Len() == 0 {
		if err := EncodeEmptySlice(w, v.Type().Elem(), label, prefix); err != nil {
			return err
		}
	}
	for i, n := 0, v.Len(); i < n; i++ {
		if err := EncodeSliceElement(w, v.Index(i).Interface(), label, prefix); err != nil {
			return err
		}
	}
	return nil
}

func EncodeMap(w io.Writer, v reflect.Value, label, prefix string) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return fmt.Errorf("encode map only works on structs")
	}
	if v.Len() == 0 {
		if err := EncodeEmptyMap(w, v.Type().Elem(), label, prefix); err != nil {
			return err
		}
	} else {
		var sv keys = v.MapKeys()
		sort.Sort(sv)
		for _, k := range sv {
			tag := k.String()
			if !(strings.IndexFunc(k.String(), unicode.IsSpace) < 0) {
				tag = strconv.Quote(tag)
			}
			if _, err := w.Write([]byte(prefix + fmt.Sprintf("[%s.%s]\n", label, tag))); err != nil {
				return err
			}
			if err := EncodeStruct(w, v.MapIndex(k).Interface(), prefix+"\t"); err != nil {
				return err
			}
		}
	}
	return nil
}

func EncodeEmptyElement(w io.Writer, v reflect.Value, label, prefix string) error {

	switch v.Type().Elem().Kind() {
	case reflect.Bool:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = true|false\n", label))); err != nil {
			return err
		}
	case reflect.String:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = \"text\"\n", label))); err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = integer\n", label))); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = unsigned integer\n", label))); err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = number\n", label))); err != nil {
			return err
		}
	case reflect.Complex64, reflect.Complex128:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = complex number\n", label))); err != nil {
			return err
		}
	default:
		if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = value\n", label))); err != nil {
			return err
		}
	}

	return nil
}

func EncodeStruct(w io.Writer, i interface{}, prefix string) error {
	if reflect.TypeOf(i).Kind() != reflect.Struct {
		return fmt.Errorf("encode field only works on structs")
	}
	for n := 0; n < reflect.TypeOf(i).NumField(); n++ {
		name := reflect.TypeOf(i).Field(n).Name

		label := strings.ToLower(name)
		if !(strings.IndexFunc(label, unicode.IsSpace) < 0) {
			label = strconv.Quote(label)
		}
		l, err := EncodedFieldTag(reflect.TypeOf(i), name, "toml")
		if err != nil {
			return err
		}
		if l != "" {
			label = l
		}
		s, err := EncodedFieldTag(reflect.TypeOf(i), name, "comment")
		if err != nil {
			return err
		}
		if s != "" {
			if _, err := w.Write([]byte(prefix + fmt.Sprintf("# %s\n", s))); err != nil {
				return err
			}
		}
		v := reflect.ValueOf(i).FieldByName(name)
		switch v.Kind() {
		case reflect.Map, reflect.Slice:
			t := v.Type().Elem()
			if t.Kind() == reflect.Struct && t.Implements(reflect.TypeOf(new(Encoder)).Elem()) {
				switch v.Kind() {
				case reflect.Map:
					if err := EncodeMap(w, v, label, prefix); err != nil {
						return err
					}
				case reflect.Slice:
					if err := EncodeSlice(w, v, label, prefix); err != nil {
						return err
					}
				default:
					if err := EncodeElement(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
						return err
					}
				}
			} else {
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					if err := EncodeElement(w, map[string]interface{}{label: v.Interface()}, prefix+"#"); err != nil {
						return err
					}
				} else {
					if err := EncodeElement(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
						return err
					}
				}
			}
		case reflect.Ptr:
			if v.IsNil() {
				if err := EncodeEmptyElement(w, v, label, prefix); err != nil {
					return nil
				}
			} else {
				if err := EncodeElement(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
					return err
				}
			}
		default:
			if err := EncodeElement(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte(prefix + "\n")); err != nil {
			return err
		}
	}
	return nil
}
