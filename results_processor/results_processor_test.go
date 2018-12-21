package results_processor

import "testing"

type TestType struct {
	Id int
	A  int `gen_versions:""`
	B  int `gen_versions:""`
}

func TestSomething(t *testing.T) {
	t1 := &TestType{
		1,
		2,
		3,
	}
	t2 := &TestType{
		1,
		2,
		4,
	}

	_, err := processModelFieldsVersions(t1, t2)
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(1)
}
