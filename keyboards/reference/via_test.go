package reference

import (
	"encoding/json"
	"os"
	"testing"
)

func TestVIADefinitionMatchesReferenceIdentityAndMatrix(t *testing.T) {
	data, err := os.ReadFile("via.json")
	if err != nil {
		t.Fatal(err)
	}
	var definition struct {
		Name      string `json:"name"`
		VendorID  string `json:"vendorId"`
		ProductID string `json:"productId"`
		Matrix    struct {
			Rows int `json:"rows"`
			Cols int `json:"cols"`
		} `json:"matrix"`
	}
	if err := json.Unmarshal(data, &definition); err != nil {
		t.Fatal(err)
	}
	if definition.Name != "TGK XIAO BLE Reference" || definition.VendorID != "0x1209" || definition.ProductID != "0x0001" || definition.Matrix.Rows != 1 || definition.Matrix.Cols != 1 {
		t.Fatalf("unexpected VIA definition: %#v", definition)
	}
}
