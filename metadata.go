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

const DateTimeFormat = "2006-01-02 15:04:05"

type Encoder interface {
	Encode(w io.Writer, prefix string) error
}

func EncodeStringPrefix(i Encoder, prefix string) (string, error) {
	buf := new(bytes.Buffer)
	err := i.Encode(buf, prefix)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func EncodeString(i Encoder) (string, error) {
	return EncodeStringPrefix(i, "")
}

func FieldTag(t reflect.Type, name, lookup string) (string, error) {
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

func EncodeTOML(w io.Writer, i interface{}, prefix string) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(i); err != nil {
		return err
	}
	for _, l := range strings.Split(buf.String(), "\n") {
		if _, err := w.Write([]byte(prefix + l + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func EncodeField(w io.Writer, i interface{}, prefix string) error {
	if reflect.TypeOf(i).Kind() != reflect.Struct {
		return fmt.Errorf("encode field only works on structs")
	}
	for n := 0; n < reflect.TypeOf(i).NumField(); n++ {
		name := reflect.TypeOf(i).Field(n).Name

		label := strings.ToLower(name)
		if !(strings.IndexFunc(label, unicode.IsSpace) < 0) {
			label = strconv.Quote(label)
		}
		n, err := FieldTag(reflect.TypeOf(i), name, "toml")
		if err != nil {
			return err
		}
		if n != "" {
			label = n
		}
		s, err := FieldTag(reflect.TypeOf(i), name, "comment")
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
						if err := EncodeField(w, v.MapIndex(k).Interface(), prefix+"\t"); err != nil {
							return err
						}
					}
				case reflect.Slice:
					for i, n := 0, v.Len(); i < n; i++ {
						if _, err := w.Write([]byte(prefix + fmt.Sprintf("[[%s]]\n", label))); err != nil {
							return err
						}
						if err := EncodeField(w, v.Index(i).Interface(), prefix+"\t"); err != nil {
							return err
						}
					}
				default:
					if err := EncodeTOML(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
						return err
					}
				}
			} else {
				if v.Kind() == reflect.Slice && v.Len() == 0 {
					if err := EncodeTOML(w, map[string]interface{}{label: v.Interface()}, prefix+"#"); err != nil {
						return err
					}
				} else {
					if err := EncodeTOML(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
						return err
					}
				}
			}
		case reflect.Ptr:
			if v.IsNil() {
				switch v.Type().Elem().Kind() {
				case reflect.Bool:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = true|false\n\n", label))); err != nil {
						return err
					}
				case reflect.String:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = \"text\"\n\n", label))); err != nil {
						return err
					}
				case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = integer\n\n", label))); err != nil {
						return err
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = unsigned integer\n\n", label))); err != nil {
						return err
					}
				case reflect.Float32, reflect.Float64:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = number\n\n", label))); err != nil {
						return err
					}
				case reflect.Complex64, reflect.Complex128:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = complex number\n\n", label))); err != nil {
						return err
					}
				default:
					if _, err := w.Write([]byte(prefix + fmt.Sprintf("#%s = value\n\n", label))); err != nil {
						return err
					}
				}
			} else {
				if err := EncodeTOML(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
					return err
				}
			}
		default:
			if err := EncodeTOML(w, map[string]interface{}{label: v.Interface()}, prefix); err != nil {
				return err
			}
		}
	}
	return nil
}
