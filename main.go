package main

import (
	"fmt"
	"net/http"

	"github.com/beerzezy/TestPackage/null"
	"github.com/beerzezy/TestPackage/restutil"
)

type Account struct {
	Firstname null.String `json:"first_name"`
	Lastname  null.String `json:"last_name"`
	Code      null.Int64  `json:"code"`
	Code2     null.Int64  `json:"code2"`
	Phone     null.String `json:"phone"`
	//Detail    []detail    `json:"detail"`
}

// type detail struct {
// 	Doc null.String `json:"doc"`
// }

func getAccount(w http.ResponseWriter, r *http.Request) {

	var resBody Account
	defer restutil.WriteResponse(w, &resBody)

	resBody = Account{
		Firstname: null.NewString("Chaiyarin"),
		//Lastname: null.NewString("Chaiyarin"),
		Code: null.NewInt64(0),
		//Code2:    null.NewInt64(11),
		//Phone: "1111",
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/getAccount", getAccount)
	http.ListenAndServe(":8083", nil)
}

func main() {
	handleRequest()
}
