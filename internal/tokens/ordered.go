package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// node là giá trị JSON giữ thứ tự khóa; encoding/json map không giữ.
type node struct {
	kind kind
	str  string
	num  json.Number
	b    bool
	keys []string
	obj  map[string]*node
	arr  []*node
}

type kind int

const (
	kindNull kind = iota
	kindString
	kindNumber
	kindBool
	kindObject
	kindArray
)

func (n *node) get(k string) *node {
	if n == nil || n.kind != kindObject {
		return nil
	}
	return n.obj[k]
}

// plain đổi node sang giá trị Go thường: string, json.Number, bool, map, []any.
func (n *node) plain() any {
	switch n.kind {
	case kindString:
		return n.str
	case kindNumber:
		return n.num
	case kindBool:
		return n.b
	case kindObject:
		m := map[string]any{}
		for _, k := range n.keys {
			m[k] = n.obj[k].plain()
		}
		return m
	case kindArray:
		out := make([]any, len(n.arr))
		for i, x := range n.arr {
			out[i] = x.plain()
		}
		return out
	}
	return nil
}

func parseOrdered(b []byte) (*node, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	n, err := readValue(dec)
	if err != nil {
		return nil, err
	}
	if _, err := dec.Token(); err == nil {
		return nil, fmt.Errorf("dữ liệu thừa sau object")
	}
	return n, nil
}

func readValue(dec *json.Decoder) (*node, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			n := &node{kind: kindObject, obj: map[string]*node{}}
			for dec.More() {
				kt, err := dec.Token()
				if err != nil {
					return nil, err
				}
				k := kt.(string)
				v, err := readValue(dec)
				if err != nil {
					return nil, err
				}
				if _, dup := n.obj[k]; !dup {
					n.keys = append(n.keys, k)
				}
				n.obj[k] = v
			}
			_, err = dec.Token() // }
			return n, err
		case '[':
			n := &node{kind: kindArray}
			for dec.More() {
				v, err := readValue(dec)
				if err != nil {
					return nil, err
				}
				n.arr = append(n.arr, v)
			}
			_, err = dec.Token() // ]
			return n, err
		}
	case string:
		return &node{kind: kindString, str: t}, nil
	case json.Number:
		return &node{kind: kindNumber, num: t}, nil
	case bool:
		return &node{kind: kindBool, b: t}, nil
	case nil:
		return &node{kind: kindNull}, nil
	}
	return nil, fmt.Errorf("token lạ %v", tok)
}
