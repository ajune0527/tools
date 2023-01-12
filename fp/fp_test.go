package fp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFp_FormatPrint(t *testing.T) {
	data := A{
		Id:   1,
		Name: "J",
		Hobbies: []Hobby{
			{
				Name: "basketball",
			},
			{
				Name: "football",
			},
			{
				Name: "",
			},
		},
		Hobbies2: []*Hobby{
			{
				Name: "basketball",
			},
			{
				Name: "football",
			},
			{
				Name: "",
			},
		},
	}
	f := &Fp{
		ConsoleOut: false,
		Prefix:     " ",
		SkipZero:   true,
		ReplacePkg: "main",
	}
	// out := f.FormatPrint(data, make([]string, 0), 0, reflect.Struct)
	// fmt.Printf("%s", strings.Join(out, ""))
	// _ = A{Id: 1, Name: "J", Hobbies: []Hobby{{Name: "basketball"}, {Name: "football"}}, Hobbies2: []*Hobby{{Name: "basketball"}, {Name: "football"}}}

	datas := []A{data, data}
	out := f.FormatPrint(datas, make([]string, 0), 0, reflect.Struct)
	// fmt.Printf("%s\n", strings.Join(out, ""))
	s := `[]fp.A{{Id: 1,Name: "J",Hobbies: []fp.Hobby{{Name: "basketball"},{Name: "football"}},Hobbies2: []*fp.Hobby{{Name: "basketball"},{Name: "football"}}},{Id: 1,Name: "J",Hobbies: []fp.Hobby{{Name: "basketball"},{Name: "football"}},Hobbies2: []*fp.Hobby{{Name: "basketball"},{Name: "football"}}}}`
	assert.Equal(t, strings.Join(out, ""), s)
}
