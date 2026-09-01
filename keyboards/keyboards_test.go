package keyboards_test

import (
	"reflect"
	"testing"

	"github.com/tgk-project/tgk/keyboards/reference"
	"github.com/tgk-project/tgk/keyboards/template"
)

func TestReferenceAndTemplateComposeWithoutCoreChanges(t *testing.T) {
	referenceFactory, err := reference.FactoryConfig()
	if err != nil {
		t.Fatal(err)
	}
	templateFactory, err := template.FactoryConfig()
	if err != nil {
		t.Fatal(err)
	}
	if referenceFactory.Identity().LayoutID == templateFactory.Identity().LayoutID {
		t.Fatal("separate keyboard packages must use distinct LayoutID values")
	}
	if got, want := referenceFactory.Positions(), []uint16{0}; !reflect.DeepEqual(got, want) {
		t.Fatalf("reference positions = %v, want %v", got, want)
	}
	if got, want := templateFactory.Positions(), []uint16{10, 20}; !reflect.DeepEqual(got, want) {
		t.Fatalf("template positions = %v, want %v", got, want)
	}
}
