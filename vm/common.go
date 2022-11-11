package vm

func ensureNotEmpty(size int) {
	if size == 0 {
		panic("Stack is empty")
	}
}

func checkIndexBounds(size, index int) {
	if 0 > index || index >= size {
		panic("Index out of bounds")
	}
}

func ensureSpace[T any](slice []T, size, space int) []T {
	newSize := size + space
	if newSize >= len(slice) {
		capacity := max(size<<1, newSize)
		target := make([]T, capacity)
		copy(target, slice)
		return target
	}
	return slice
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
