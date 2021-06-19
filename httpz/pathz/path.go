package pathz

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	NotFound = errors.New("not found")
)

type (
	node struct {
		val string
		handle http.Handler
		childs []*node
	}

	Path struct {
		methodRoot map[string]*node    // method -> path
	}
)

func createNode(val string) *node {
	return &node{
		val: val,
		childs: make([]*node, 0),
	}
}

func NewPath() *Path {
	return &Path{
		methodRoot: make(map[string]*node),
	}
}

func (p *Path) BuildPath(method, path string, handler http.Handler) error {
	pathItems := cleanPath(path)
	if _, ok := p.methodRoot[method]; !ok {
		p.methodRoot[method] = createNode(method)
	} else {
		// root path
		if len(pathItems) == 0 {
			log.Fatal("The route is repeat, method:", method, " path:", path)
		}
	}

	// root path
	if len(pathItems) == 0 {
		p.methodRoot[method].handle = handler
		return nil
	}

	lastNode := p.methodRoot[method]
	pathExist := false
	for _, v := range pathItems {
		lastNode, pathExist = p.addNode(lastNode, v)
	}

	// route is repeat
	if pathExist {
		log.Fatal("The route is repeat, method:", method, " path:", path)
	}

	lastNode.handle = handler

	return nil
}

func (p *Path) addNode(root *node, val string) (*node, bool) {
	for _, v := range root.childs {
		if v.val == val {
			return v, true
		}
	}

	newNode := createNode(val)
	root.childs = append(root.childs, newNode)

	return newNode, false
}


func (p *Path) ParsePath(method, path string) (http.Handler, map[string]interface{}, error) {
	var (
		respMap = make(map[string]interface{})
	)
	if _, ok := p.methodRoot[method]; !ok {
		return nil, nil, NotFound
	}

	pathItems := cleanPath(path)

	// root path
	if len(pathItems) == 0 {
		if p.methodRoot[method].handle == nil {
			return nil, respMap, NotFound
		}
		return p.methodRoot[method].handle, respMap, nil
	}

	// route stack
	r := make([]*node, 0)

	// loop stack
	q := make([]*node, 0)
	q = append(q, p.methodRoot[method])
	idx := 0
	for len(q) != 0 && idx <= len(pathItems) {
		curNode := q[len(q) - 1]
		q = q[:len(q)-1]

		r = append(r, curNode)
		if idx == len(pathItems) && curNode.handle != nil {
			break
		}

		ns := p.parseNode(curNode, pathItems[idx])
		if len(ns) == 0 {
			r = r[len(r)-1:]
			idx--
		} else {
			idx++
		}

		for _, n := range ns {
			q = append(q, n)
		}
	}

	if len(r) == 0 || r[len(r)-1].handle == nil || len(r) - 1 != len(pathItems) {
		return nil, respMap, NotFound
	}

	r = r[1:]

	for i, v := range r {
		if v.val[0] == ':' {
			k := v.val[1:]
			respMap[k] = pathItems[i]
		}
	}

	return r[len(r)-1].handle, respMap, nil
}

func (p *Path) parseNode(root *node, val string) []*node {
	resp := make([]*node, 0)
	for _, v := range root.childs {
		if v.val == val && v.val[0] != ':' {
			resp = append(resp, v)
			break
		}
	}

	for _, v := range root.childs {
		if v.val[0] == ':' {
			resp = append(resp, v)
		}
	}

	return resp
}


func cleanPath(path string) []string {
	ps := strings.Split(path, "?")
	path = ps[0]

	if path[0] == '/' {
		path = path[1:]
	}

	items := strings.Split(path, "/")

	pathItems := make([]string, 0)
	for _, v := range items {
		if len(v) == 0 {
			continue
		}
		pathItems = append(pathItems, v)
	}

	return pathItems
}
