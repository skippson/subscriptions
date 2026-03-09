package models

func PtrToOptional[T any](p *T) Optional[T] {
	if p == nil {
		return Optional[T]{Valid: false}
	}
	return Optional[T]{Value: *p, Valid: true}
}

func OptionalToPtr[T any](o Optional[T]) *T {
	if !o.Valid {
		return nil
	}
	return &o.Value
}

func ValueToOption[T any](v T) Optional[T]{
	return Optional[T]{
		Valid: true,
		Value: v,
	}
}