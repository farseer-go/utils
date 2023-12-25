package pointer

import "reflect"

// Of 返回一个指向值' v '的指针
func Of[T any](v T) *T {
	return &v
}

// Unwrap 返回指针的值，如果指针为nil则返回val
func Unwrap[T any](p *T, val ...T) T {
	var d T
	if len(val) > 0 {
		d = val[0]
	}
	if p == nil {
		return d
	}
	return *p
}

// Extract 递归拆解指针，获取指向值
func Extract(value any) any {
	if value == nil {
		return nil
	}
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if t.Kind() != reflect.Pointer {
		return value
	}
	return Extract(v.Elem().Interface())
}
