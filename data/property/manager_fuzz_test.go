package property

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

const fuzzKeys = 8

// referenceFinalize is Finalize as it was before FLOGO-19663 (single-threaded, in place): the model that
// Finalize must still match.
func referenceFinalize(properties map[string]Property, processors ...PostProcessor) {
	valueMap := make(map[string]interface{}, len(properties))
	for k, v := range properties {
		valueMap[k] = v.Value
	}
	for _, processor := range processors {
		_ = processor(valueMap)
	}
	for k, v := range valueMap {
		if prop, exists := properties[k]; exists {
			prop.Value = v
			properties[k] = prop
		}
	}
}

func fuzzInitial() map[string]Property {
	props := make(map[string]Property, fuzzKeys)
	for i := 0; i < fuzzKeys; i++ {
		props["p"+strconv.Itoa(i)] = Property{Value: "v" + strconv.Itoa(i), Overridable: i%2 == 0}
	}
	return props
}

// fuzzProcessors turns ops into up to 3 processors. Each op is a byte pair (code, arg). Bits 0-1 of code pick the
// action (set an existing key, add a new key, delete a key, set a key to nil), bits 2-3 the processor, and bit 4
// makes that processor return an error after its actions.
func fuzzProcessors(ops []byte) []PostProcessor {
	type action struct{ code, arg byte }
	var actions [3][]action
	var fails [3]bool
	used := 0
	for i := 0; i+1 < len(ops); i += 2 {
		code, arg := ops[i], ops[i+1]
		p := int(code>>2&0x03) % 3
		actions[p] = append(actions[p], action{code, arg})
		if code&0x10 != 0 {
			fails[p] = true
		}
		if p >= used {
			used = p + 1
		}
	}

	processors := make([]PostProcessor, 0, used)
	for p := 0; p < used; p++ {
		acts, fail := actions[p], fails[p]
		processors = append(processors, func(properties map[string]interface{}) error {
			for _, a := range acts {
				key := "p" + strconv.Itoa(int(a.arg)%fuzzKeys)
				switch a.code & 0x03 {
				case 0:
					properties[key] = "x" + strconv.Itoa(int(a.arg))
				case 1:
					properties["new"+strconv.Itoa(int(a.arg))] = a.arg
				case 2:
					delete(properties, key)
				default:
					properties[key] = nil
				}
			}
			if fail {
				return errors.New("processor failed")
			}
			return nil
		})
	}
	return processors
}

// FuzzManager_FinalizeMatchesReference checks that Finalize gives the same properties as the old algorithm for
// any processor behaviour, and that it never modifies the snapshot it replaces.
func FuzzManager_FinalizeMatchesReference(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x00, 0x01})
	f.Add([]byte{0x01, 0x02, 0x02, 0x03, 0x03, 0x04})
	f.Add([]byte{0x14, 0x05, 0x08, 0x05, 0x02, 0x05})
	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) > 128 {
			ops = ops[:128]
		}
		processors := fuzzProcessors(ops)

		want := fuzzInitial()
		referenceFinalize(want, processors...)

		m := NewManager(fuzzInitial())
		published := m.snapshot()
		if err := m.Finalize(processors...); err != nil {
			t.Fatalf("Finalize returned %v", err)
		}

		if got := m.snapshot(); !reflect.DeepEqual(got, want) {
			t.Fatalf("Finalize gave %v, the reference model %v", got, want)
		}
		if !reflect.DeepEqual(published, fuzzInitial()) {
			t.Fatalf("Finalize modified the snapshot it replaced: %v", published)
		}
		values := make(map[string]interface{}, len(want))
		for k, v := range want {
			values[k] = v.Value
		}
		if got := m.GetProperties(); !reflect.DeepEqual(got, values) {
			t.Fatalf("GetProperties gave %v, want %v", got, values)
		}
	})
}
