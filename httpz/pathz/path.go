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
	innerResult struct {
		key string
		val string
		named bool
		found bool
	}

	node struct {
		val string
		handle http.Handler
		childs []*node
	}

	Path struct {
		methodRoot map[string]*node    // method -> path
	}

	Result struct {
		Handle http.Handler
		Params map[string]string
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
	_ , err := p.ParsePath(method, path)
	if err != NotFound {
		log.Fatalf("the route rule is exists method:%s, path:%s", method, path)
	}

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


func (p *Path) next(n *node, path string, result *Result) bool {
	if len(path) == 0 && n.handle != nil {
		result.Handle = n.handle
		return true
	}

	for i := range path {
		if path[i] == '/' {
			token := path[:i]
			for _, child := range n.childs {
				if r := match(child.val, token); r.found {
					if p.next(child, path[i+1:], result) {
						if r.named {
							addResult(result, r.key, r.val)
						}
						return true
					}
				}
			}
			return false
		}
	}

	for _, child := range n.childs {
		r := match(child.val, path)
		if r.found && child.handle != nil {
			result.Handle = child.handle
			if r.named {
				addResult(result, r.key, r.val)
			}
			return true
		}
	}

	return false
}

func addResult(result *Result, k, v string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}
	result.Params[k] = v
}


func match(k, token string) innerResult {
	if k[0] == ':' {
		return innerResult{
			key: k[1:],
			val: token,
			named: true,
			found: true,
		}
	}

	return innerResult{
		found: k == token,
	}
}


func (p *Path) ParsePath(method, path string) (Result, error) {
	var (
		resp = Result{}
	)
	if _, ok := p.methodRoot[method]; !ok {
		return resp, NotFound
	}

	pathItems := cleanPath(path)
	// root path
	if len(pathItems) == 0 {
		if p.methodRoot[method].handle == nil {
			return resp, NotFound
		}
		resp.Handle = p.methodRoot[method].handle
		return resp, nil
	}

	newPath := strings.Join(pathItems, "/")
	if !p.next(p.methodRoot[method], newPath, &resp) {
		return resp, NotFound
	}

	return resp, nil
}

func (p *Path) parseNode(root *node, val string) []*node {
	resp := make([]*node, 0)

	for _, v := range root.childs {
		if v.val == val || v.val[0] == ':'{
			resp = append(resp, v)
		}
	}

	return resp
}


func cleanPath(path string) []string {
	pathBytes := []byte(path)

	l := 0
	pathLen := len(pathBytes)
	isFirst := true

	newPath := make([]byte, 0)
	for l < pathLen {
		if pathBytes[l] == '?' {
			break
		}

		r := l + 1
		if r >= pathLen {
			r = pathLen - 1
		}

		if pathBytes[l] == pathBytes[r] && pathBytes[l] == '/' {
			l++
			continue
		}

		if isFirst && pathBytes[l] == '/' {
			l++
			continue
		}

		newPath = append(newPath, pathBytes[l])
		isFirst = false
		l++
	}

	if len(newPath) == 0 {
		return []string{}
	}

	if newPath[len(newPath)-1] == '/' {
		newPath = newPath[:len(newPath)-1]
	}

	tmpPath := string(newPath)
	return strings.Split(tmpPath, "/")
}
