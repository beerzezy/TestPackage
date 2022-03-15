package restutil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beerzezy/TestPackage/null"
)

func WriteResponse(w http.ResponseWriter, obj interface{}) {
	b, _ := json.Marshal(obj)
	bs := null.StripNullJSON(b)
	output := string(bs)
	fmt.Fprint(w, output)
}
