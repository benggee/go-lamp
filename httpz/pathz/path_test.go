package pathz

import (
	"fmt"
	"net/http"
	"testing"
)

func getHttpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func TestPath(t *testing.T) {
	p := NewPathz()

	p.BuildPath("POST", "/user/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AAAAA")
	}))

	p.BuildPath("POST", "/user/info/:userId", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AAAAA")
	}))

	p.BuildPath("POST", "/user/info/:userId/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AAAAA")
	}))


	fmt.Println(p.methodRoot["POST"].val)
	fmt.Println(p.methodRoot["POST"].childs[0].val)
	fmt.Println(p.methodRoot["POST"].childs[0].childs[0].val)
	fmt.Println(p.methodRoot["POST"].childs[0].childs[0].childs[0].val)
	fmt.Println(p.methodRoot["POST"].childs[0].childs[0].childs[0].childs[0].val)


	h, param, e := p.ParsePath("POST", "/user/info/aaaaa/test")

	fmt.Println("HH:", h, " PP:", param, " EE:",e)

}
