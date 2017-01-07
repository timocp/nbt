package nbt

import "bufio"
import "compress/gzip"
import "fmt"
import "os"
import "testing"

// test parsing of the sample files referenced from
// http://web.archive.org/web/20110723210920/http://www.minecraft.net/docs/NBT.txt

func TestExample(t *testing.T) {
	f, err := os.Open("examples/test.nbt")
	must(err)
	defer f.Close()
	r := bufio.NewReader(f)
	data := Parse(r)
	root := expectCompoundTag(t, data, "hello world", 1)
	_ = expectStringTag(t, root.ChildByName("name"), "name", "Bananrama")
}

// bigtest.nbt is gzip compressed
func TestBigExample(t *testing.T) {
	f, err := os.Open("examples/bigtest.nbt")
	must(err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	must(err)
	defer r.Close()
	data := Parse(r)
	root := expectCompoundTag(t, data, "Level", 11)
	_ = expectShortTag(t, root.ChildByName("shortTest"), "shortTest", 32767)
	_ = expectLongTag(t, root.ChildByName("longTest"), "longTest", 9223372036854775807)
	_ = expectFloatTag(t, root.ChildByName("floatTest"), "floatTest", 0.49823147)
	_ = expectStringTag(t, root.ChildByName("stringTest"), "stringTest", "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!")
	_ = expectIntTag(t, root.ChildByName("intTest"), "intTest", 2147483647)
	c := expectCompoundTag(t, root.ChildByName("nested compound test"), "nested compound test", 2)
	c1 := expectCompoundTag(t, c.ChildByName("ham"), "ham", 2)
	_ = expectStringTag(t, c1.ChildByName("name"), "name", "Hampus")
	_ = expectFloatTag(t, c1.ChildByName("value"), "value", 0.75)
	c1 = expectCompoundTag(t, c.ChildByName("egg"), "egg", 2)
	_ = expectStringTag(t, c1.ChildByName("name"), "name", "Eggbert")
	_ = expectFloatTag(t, c1.ChildByName("value"), "value", 0.5)
	list := expectListTag(t, root.ChildByName("listTest (long)"), "listTest (long)", 5, tag_long)
	for i := 0; i < len(list.Values); i++ {
		_ = expectLongTag(t, list.Values[i], "", int64(i+11))
	}
	_ = expectByteTag(t, root.ChildByName("byteTest"), "byteTest", 127)
	list = expectListTag(t, root.ChildByName("listTest (compound)"), "listTest (compound)", 2, tag_compound)
	for i := 0; i < len(list.Values); i++ {
		c = expectCompoundTag(t, list.Values[i], "", 2)
		_ = expectStringTag(t, c.ChildByName("name"), "name", fmt.Sprintf("Compound tag #%d", i))
	}
	b := expectByteArrayTag(t, root.ChildByName("byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"), "byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))", 1000)
	for i := 0; i < len(b.Values); i++ {
		expectedByte := byte((i*i*255 + i*7) % 100)
		if b.Values[i] != expectedByte {
			t.Errorf("Expected %dth element to be %d, got %d", expectedByte, b.Values[i])
		}
	}
}

func expectCompoundTag(t *testing.T, tag Tag, expectedName string, expectedSize int) CompoundTag {
	v, ok := tag.(CompoundTag)
	if !ok {
		t.Errorf("Expected CompoundTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if len(v.Values) != expectedSize {
		t.Errorf("Expected compound tag size to be %d, got %d", expectedSize, len(v.Values))
	}
	return v
}

func expectStringTag(t *testing.T, tag Tag, expectedName string, expectedValue string) StringTag {
	v, ok := tag.(StringTag)
	if !ok {
		t.Errorf("Expected StringTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected string tag value to be %q, got %q", expectedValue, v.Value)
	}
	return v
}

func expectShortTag(t *testing.T, tag Tag, expectedName string, expectedValue int16) ShortTag {
	v, ok := tag.(ShortTag)
	if !ok {
		t.Errorf("Expected ShortTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected short tag value to be %d, got %d", expectedValue, v.Value)
	}
	return v
}

func expectLongTag(t *testing.T, tag Tag, expectedName string, expectedValue int64) LongTag {
	v, ok := tag.(LongTag)
	if !ok {
		t.Errorf("Expected LongTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected long tag value to be %d, got %d", expectedValue, v.Value)
	}
	return v
}

func expectFloatTag(t *testing.T, tag Tag, expectedName string, expectedValue float32) FloatTag {
	v, ok := tag.(FloatTag)
	if !ok {
		t.Errorf("Expected FloatTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected float tag value to be %f, got %f", expectedValue, v.Value)
	}
	return v
}

func expectIntTag(t *testing.T, tag Tag, expectedName string, expectedValue int32) IntTag {
	v, ok := tag.(IntTag)
	if !ok {
		t.Errorf("Expected IntTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected int tag value to be %d, got %d", expectedValue, v.Value)
	}
	return v
}

func expectListTag(t *testing.T, tag Tag, expectedName string, expectedSize int, expectedType int) ListTag {
	v, ok := tag.(ListTag)
	if !ok {
		t.Errorf("Expected ListTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if len(v.Values) != expectedSize {
		t.Errorf("Expected list tag size to be %d, got %d", expectedSize, len(v.Values))
	}
	if v.TagType != expectedType {
		t.Errorf("Expected list tag type to be %d, got %d", expectedType, v.TagType)
	}
	return v
}

func expectByteTag(t *testing.T, tag Tag, expectedName string, expectedValue int8) ByteTag {
	v, ok := tag.(ByteTag)
	if !ok {
		t.Errorf("Expected ByteTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if v.Value != expectedValue {
		t.Errorf("Expected byte tag value to be %d, got %d", expectedValue, v.Value)
	}
	return v
}

func expectByteArrayTag(t *testing.T, tag Tag, expectedName string, expectedSize int) ByteArrayTag {
	v, ok := tag.(ByteArrayTag)
	if !ok {
		t.Errorf("Expected ByteArrayTag, got %T", tag)
	}
	expectName(t, v, expectedName)
	if len(v.Values) != expectedSize {
		t.Errorf("Expected byte array tag size to be %d, got %d", expectedSize, len(v.Values))
	}
	return v
}

func expectName(t *testing.T, tag Tag, expectedName string) {
	if tag.GetName() != expectedName {
		t.Errorf("Expected tag name to be %q, got %q", expectedName, tag.GetName())
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
