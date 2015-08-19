package newton

import (
	"net/url"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ReaderTestFixture struct {
	*gunit.Fixture

	sources []Source
	reader  *Reader
}

func (this *ReaderTestFixture) Setup() {
	this.sources = []Source{
		&FakeSource{},
		&FakeSource{key: "string", value: []string{"asdf"}},
		&FakeSource{key: "string", value: []string{"qwer"}},
		&FakeSource{key: "string-no-values", value: []string{}},
		&FakeSource{key: "int", value: []string{"42"}},
		&FakeSource{key: "int", value: []string{"-1"}},
		&FakeSource{key: "int-bad", value: []string{"not an integer"}},
		&FakeSource{key: "bool", value: []string{"true"}},
		&FakeSource{key: "bool-bad", value: []string{"not a bool"}},
		&FakeSource{key: "url", value: []string{"http://www.google.com"}},
		&FakeSource{key: "url-bad", value: []string{"%%%%%%"}}, // not a url
	}

	this.reader = NewReader(this.sources...)
}

////////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestInitializeSources() {
	for _, source := range this.sources {
		this.So(source.(*FakeSource).initialized, should.Equal, 1)
	}
}

////////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestNilSourcesAreSkipped() {
	source1 := &FakeSource{key: "1"}
	source2 := &FakeSource{key: "2"}
	this.sources = []Source{source1, nil, nil, nil, source2}

	this.So(func() { this.reader = NewReader(this.sources...) }, should.NotPanic)
	this.So(this.reader.sources, should.Resemble, []Source{source1, source2})
}

////////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestStrings_Found() {
	value := this.reader.Strings("string")

	this.So(value, should.Resemble, []string{"asdf"})
}

func (this *ReaderTestFixture) TestStrings_NotFound() {
	value := this.reader.Strings("blahblah")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestStringsError_Found() {
	value, err := this.reader.StringsError("string")

	this.So(value, should.Resemble, []string{"asdf"})
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestStringsError_NotFound() {
	value, err := this.reader.StringsError("81")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringsPanic_Found() {
	value := this.reader.StringsPanic("string")

	this.So(value, should.Resemble, []string{"asdf"})
}

func (this *ReaderTestFixture) TestStringsPanic_NotFound() {
	this.So(func() { this.reader.StringsPanic("blahblah") }, should.Panic)
}

func (this *ReaderTestFixture) TestStringsFatal_Found() {
	value := this.reader.StringsFatal("string")

	this.So(value, should.Resemble, []string{"asdf"})
}

func (this *ReaderTestFixture) TestStringsFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.StringsFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringsDefault_Found() {
	value := this.reader.StringsDefault("string", []string{"default"})

	this.So(value, should.Resemble, []string{"asdf"})
}

func (this *ReaderTestFixture) TestStringsDefault_NotFound() {
	value := this.reader.StringsDefault("blahblah", []string{"default"})

	this.So(value, should.Resemble, []string{"default"})
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestString_Found() {
	value := this.reader.String("string")

	this.So(value, should.Equal, "asdf")
}
func (this *ReaderTestFixture) TestString_NotFound() {
	value := this.reader.String("blahblah")

	this.So(value, should.BeEmpty)
}

func (this *ReaderTestFixture) TestStringError_Found() {
	value, err := this.reader.StringError("string")

	this.So(value, should.Resemble, "asdf")
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestStringError_NotFound() {
	value, err := this.reader.StringError("81")

	this.So(value, should.Equal, "")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringError_FoundButNoValuesProvided() {
	value, err := this.reader.StringError("string-no-values")

	this.So(value, should.Equal, "")
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestStringPanic_Found() {
	value := this.reader.StringPanic("string")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringPanic_NotFound() {
	this.So(func() { this.reader.StringPanic("blahblah") }, should.Panic)
}

func (this *ReaderTestFixture) TestStringFatal_Found() {
	value := this.reader.StringFatal("string")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.StringFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringDefault_Found() {
	value := this.reader.StringDefault("string", "default")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringDefault_NotFound() {
	value := this.reader.StringDefault("blahblah", "default")

	this.So(value, should.Resemble, "default")
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestIntsError_Found() {
	value, err := this.reader.IntsError("int")

	this.So(value, should.Resemble, []int{42})
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestIntsError_NotFound() {
	value, err := this.reader.IntsError("asdf")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestIntsError_MalformedValue() {
	value, err := this.reader.IntsError("int-bad")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestInts_Found() {
	value := this.reader.Ints("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestInts_NotFound() {
	value := this.reader.Ints("qrew")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestInts_MalformedValue() {
	value := this.reader.Ints("int-bad")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestIntsPanic_Found() {
	value := this.reader.IntsPanic("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsPanic_NotFound() {
	this.So(func() { this.reader.IntsPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntsPanic_MalformedValue() {
	this.So(func() { this.reader.IntsPanic("int-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntsFatal_Found() {
	value := this.reader.IntsFatal("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntsFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestIntsFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntsFatal("int-bad")
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestIntsDefault_Found() {
	value := this.reader.IntsDefault("int", []int{84})

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsDefault_NotFound() {
	value := this.reader.IntsDefault("missing", []int{84})

	this.So(value, should.Resemble, []int{84})
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestIntError_Found() {
	value, err := this.reader.IntError("int")

	this.So(value, should.Resemble, 42)
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestIntError_NotFound() {
	value, err := this.reader.IntError("asdf")

	this.So(value, should.Equal, 0)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestIntError_MalformedValue() {
	value, err := this.reader.IntError("int-bad")

	this.So(value, should.Equal, 0)
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestInt_Found() {
	value := this.reader.Int("int")

	this.So(value, should.Resemble, 42)
}

func (this *ReaderTestFixture) TestInt_NotFound() {
	value := this.reader.Int("qrew")

	this.So(value, should.Equal, 0)
}

func (this *ReaderTestFixture) TestInt_MalformedValue() {
	value := this.reader.Int("int-bad")

	this.So(value, should.Equal, 0)
}

func (this *ReaderTestFixture) TestIntPanic_Found() {
	value := this.reader.IntPanic("int")

	this.So(value, should.Resemble, 42)
}

func (this *ReaderTestFixture) TestIntPanic_NotFound() {
	this.So(func() { this.reader.IntPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntPanic_MalformedValue() {
	this.So(func() { this.reader.IntPanic("int-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntFatal_Found() {
	value := this.reader.IntFatal("int")

	this.So(value, should.Resemble, 42)
}

func (this *ReaderTestFixture) TestIntFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestIntFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntFatal("int-bad")
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestIntDefault_Found() {
	value := this.reader.IntDefault("int", 84)

	this.So(value, should.Resemble, 42)
}

func (this *ReaderTestFixture) TestIntDefault_NotFound() {
	value := this.reader.IntDefault("missing", 84)

	this.So(value, should.Resemble, 84)
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestBoolError_Found() {
	value, err := this.reader.BoolError("bool")

	this.So(value, should.BeTrue)
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestBoolError_NotFound() {
	value, err := this.reader.BoolError("asdf")

	this.So(value, should.BeFalse)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestBoolError_MalformedValue() {
	value, err := this.reader.BoolError("bool-bad")

	this.So(value, should.BeFalse)
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestBool_Found() {
	value := this.reader.Bool("bool")

	this.So(value, should.BeTrue)
}

func (this *ReaderTestFixture) TestBool_NotFound() {
	value := this.reader.Bool("qrew")

	this.So(value, should.BeFalse)
}

func (this *ReaderTestFixture) TestBool_MalformedValue() {
	value := this.reader.Bool("bool-bad")

	this.So(value, should.BeFalse)
}

func (this *ReaderTestFixture) TestBoolPanic_Found() {
	value := this.reader.BoolPanic("bool")

	this.So(value, should.BeTrue)
}

func (this *ReaderTestFixture) TestBoolPanic_NotFound() {
	this.So(func() { this.reader.BoolPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestBoolPanic_MalformedValue() {
	this.So(func() { this.reader.BoolPanic("bool-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestBoolFatal_Found() {
	value := this.reader.BoolFatal("bool")

	this.So(value, should.BeTrue)
}

func (this *ReaderTestFixture) TestBoolFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.BoolFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestBoolFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.BoolFatal("bool-bad")
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestBoolDefault_Found() {
	value := this.reader.BoolDefault("bool", false)

	this.So(value, should.BeTrue)
}

func (this *ReaderTestFixture) TestBoolDefault_NotFound() {
	value := this.reader.BoolDefault("missing", true)

	this.So(value, should.BeTrue)
}

//////////////////////////////////////////////////////////////

var _validURL, _ = url.Parse("http://www.google.com")
var validURL = *_validURL

func (this *ReaderTestFixture) TestURLsError_Found() {
	value, err := this.reader.URLsError("url")

	this.So(value, should.Resemble, []url.URL{validURL})
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestURLsError_NotFound() {
	value, err := this.reader.URLsError("asdf")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestURLsError_MalformedValue() {
	value, err := this.reader.URLsError("url-bad")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestURLs_Found() {
	value := this.reader.URLs("url")

	this.So(value, should.Resemble, []url.URL{validURL})
}

func (this *ReaderTestFixture) TestURLs_NotFound() {
	value := this.reader.URLs("qrew")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestURLs_MalformedValue() {
	value := this.reader.URLs("url-bad")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestURLsPanic_Found() {
	value := this.reader.URLsPanic("url")

	this.So(value, should.Resemble, []url.URL{validURL})
}

func (this *ReaderTestFixture) TestURLsPanic_NotFound() {
	this.So(func() { this.reader.URLsPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestURLsPanic_MalformedValue() {
	this.So(func() { this.reader.URLsPanic("url-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestURLsFatal_Found() {
	value := this.reader.URLsFatal("url")

	this.So(value, should.Resemble, []url.URL{validURL})
}

func (this *ReaderTestFixture) TestURLsFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.URLsFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestURLsFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.URLsFatal("url-bad")
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestURLsDefault_Found() {
	value := this.reader.URLsDefault("url", []url.URL{})

	this.So(value, should.Resemble, []url.URL{validURL})
}

func (this *ReaderTestFixture) TestURLsDefault_NotFound() {
	value := this.reader.URLsDefault("missing", []url.URL{url.URL{}})

	this.So(value, should.Resemble, []url.URL{url.URL{}})
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestURLError_Found() {
	value, err := this.reader.URLError("url")

	this.So(value, should.Resemble, validURL)
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestURLError_NotFound() {
	value, err := this.reader.URLError("asdf")

	this.So(value, should.Resemble, url.URL{})
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestURLError_MalformedValue() {
	value, err := this.reader.URLError("url-bad")

	this.So(value, should.Resemble, url.URL{})
	this.So(err, should.Equal, MalformedValueError)
}

func (this *ReaderTestFixture) TestURL_Found() {
	value := this.reader.URL("url")

	this.So(value, should.Resemble, validURL)
}

func (this *ReaderTestFixture) TestURL_NotFound() {
	value := this.reader.URL("qrew")

	this.So(value, should.Resemble, url.URL{})
}

func (this *ReaderTestFixture) TestURL_MalformedValue() {
	value := this.reader.URL("url-bad")

	this.So(value, should.Resemble, url.URL{})
}

func (this *ReaderTestFixture) TestURLPanic_Found() {
	value := this.reader.URLPanic("url")

	this.So(value, should.Resemble, validURL)
}

func (this *ReaderTestFixture) TestURLPanic_NotFound() {
	this.So(func() { this.reader.URLPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestURLPanic_MalformedValue() {
	this.So(func() { this.reader.URLPanic("url-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestURLFatal_Found() {
	value := this.reader.URLFatal("url")

	this.So(value, should.Resemble, validURL)
}

func (this *ReaderTestFixture) TestURLFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.URLFatal("balhaafslk")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestURLFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.URLFatal("url-bad")
	this.So(err, should.Equal, MalformedValueError)
}

var _defaultURL, _ = url.Parse("http://bing.com")
var defaultURL = *_defaultURL

func (this *ReaderTestFixture) TestURLDefault_Found() {
	value := this.reader.URLDefault("url", defaultURL)

	this.So(value, should.Resemble, validURL)
}

func (this *ReaderTestFixture) TestURLDefault_NotFound() {
	value := this.reader.URLDefault("missing", defaultURL)

	this.So(value, should.Resemble, defaultURL)
}

//////////////////////////////////////////////////////////////

type FakeSource struct {
	key         string
	value       []string
	initialized int
}

func (this *FakeSource) Initialize() {
	this.initialized++
}

func (this *FakeSource) Strings(key string) ([]string, error) {
	if key == this.key {
		return this.value, nil
	}
	return nil, KeyNotFoundError
}
