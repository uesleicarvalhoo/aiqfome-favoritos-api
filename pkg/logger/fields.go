package logger

type Fields map[string]any

func (f Fields) toKeysAndValues() []any {
	kv := make([]any, 0, len(f)*2)
	for k, v := range f {
		kv = append(kv, k, v)
	}

	return kv
}
