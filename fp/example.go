package fp

type A struct {
	Id       int
	Age      int
	Name     string
	Hobbies  []Hobby
	Hobbies2 []*Hobby
}

type Hobby struct {
	Name string
}
