package buffer

// Buffer is a simple in-memory buffer that can be used to read and write data.
// The idea about using buffer and not string, array or linked list is that
// for string: // 1. it is immutable, so every time we change it, we create a new string 2. it is not efficient for large data
// for array: // 1. it is fixed size, so we need to create a new array every time we want to change the size 2. it is not efficient for large data, or we deal with shifting data
// for linked list: // 1. it is not efficient for large data, as we need to traverse the list every time we want to access an element 2. it is not efficient for random access, as we need to traverse the list every time we want to access an element + takes more memory in case we have double pointers

// GappedTextBuffer is a text buffer that uses a gap to efficiently insert and delete text.
// like the one used in text editors like vim, emacs, etc.
type GappedTextBuffer struct {
	data     []byte
	gapStart int
	gapEnd   int
}

func NewGappedTextBuffer(cap int) *GappedTextBuffer {
	return &GappedTextBuffer{
		data:     make([]byte, cap), // Initialize the buffer with the full capacity
		gapStart: 0,
		gapEnd:   cap,
	}
}

// Get the buffer's length
func (b *GappedTextBuffer) Len() int {
	return len(b.data) - (b.gapEnd - b.gapStart)
}

// Get the buffer's capacity
func (b *GappedTextBuffer) Cap() int {
	return cap(b.data)
}

// Insert a byte at the given position
func (b *GappedTextBuffer) Insert(pos int, value byte) {
	if pos < 0 || pos > b.Len() {
		panic("Insert position out of bounds")
	}
	if b.gapStart == b.gapEnd {
		b.grow()
	}
	b.MoveCursorTo(pos)
	b.data[b.gapStart] = value
	b.gapStart++
}

// Moving the cursor to a pos
func (b *GappedTextBuffer) MoveCursorTo(pos int) {
	if pos < 0 || pos > b.Len() {
		return
	}

	if pos < b.gapStart {
		n := b.gapStart - pos
		copy(b.data[b.gapEnd-n:b.gapEnd], b.data[pos:b.gapStart])
		b.gapStart -= n
		b.gapEnd -= n
	} else if pos > b.gapStart {
		n := pos - b.gapStart
		copy(b.data[b.gapStart:b.gapStart+n], b.data[b.gapEnd:b.gapEnd+n])
		b.gapStart += n
		b.gapEnd += n
	}
}

func (b *GappedTextBuffer) grow() {
	oldData := b.data
	oldTailSize := len(oldData) - b.gapEnd

	newCap := b.Cap() * 2
	newData := make([]byte, newCap)

	// Copy head
	copy(newData, oldData[:b.gapStart])

	// Copy tail
	copy(newData[newCap-oldTailSize:], oldData[b.gapEnd:])

	b.data = newData
	b.gapEnd = newCap - oldTailSize
}

// Delete a byte at the given position
func (b *GappedTextBuffer) Delete(pos int) {
	if pos < 0 || pos >= b.Len() {
		panic("Delete position out of bounds")
	}
	b.MoveCursorTo(pos)
	b.gapEnd++
}

func (b *GappedTextBuffer) Flatten() string {
	return string(b.data[:b.gapStart]) + string(b.data[b.gapEnd:])
}
