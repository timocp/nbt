package nbt

const (
	tag_end        = 0
	tag_byte       = 1
	tag_short      = 2
	tag_int        = 3
	tag_long       = 4
	tag_float      = 5
	tag_double     = 6
	tag_byte_array = 7
	tag_string     = 8
	tag_list       = 9
	tag_compound   = 10
	tag_int_array  = 11
)

// All NBT tag types implement the Tag interface.  Their Value attribute
// contains the payload but should only be accessed after a type check.

type Tag interface {
	GetName() string
}

type EndTag struct{}

func (t EndTag) GetName() string { return "" }

type ByteTag struct {
	name  string
	Value int8
}

func (t ByteTag) GetName() string { return t.name }

type ShortTag struct {
	name  string
	Value int16
}

func (t ShortTag) GetName() string { return t.name }

type IntTag struct {
	name  string
	Value int32
}

func (t IntTag) GetName() string { return t.name }

type LongTag struct {
	name  string
	Value int64
}

func (t LongTag) GetName() string { return t.name }

type FloatTag struct {
	name  string
	Value float32
}

func (t FloatTag) GetName() string { return t.name }

type DoubleTag struct {
	name  string
	Value float64
}

func (t DoubleTag) GetName() string { return t.name }

type StringTag struct {
	name  string
	Value string
}

func (t StringTag) GetName() string { return t.name }

type ByteArrayTag struct {
	name   string
	Values []byte
}

func (t ByteArrayTag) GetName() string { return t.name }

type ListTag struct {
	name    string
	TagType int
	Values  []Tag
}

func (t ListTag) GetName() string { return t.name }

type CompoundTag struct {
	name   string
	Values []Tag
}

func (t CompoundTag) GetName() string { return t.name }

func (t CompoundTag) ChildByName(name string) Tag {
	for _, v := range t.Values {
		if v.GetName() == name {
			return v
		}
	}
	return nil
}

type IntArrayTag struct {
	name   string
	Values []int32
}

func (t IntArrayTag) GetName() string { return t.name }
