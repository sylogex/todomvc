package utils

import (
	"encoding/json"
	"fmt"
	"github.com/neelance/gopherjs/js"
	"math/rand"
	"time"
)

func Store(key string, val interface{}) {
	byteArr, _ := json.Marshal(val)
	str := string(byteArr)
	js.Global("localStorage").Call("setItem", key, str)
}
func Retrieve(key string, val interface{}) {
	item := js.Global("localStorage").Call("getItem", key)
	if item.IsNull() {
		val = nil
		return
	}
	str := item.String()
	json.Unmarshal([]byte(str), &val)
}
func Pluralize(count int, word string) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func Uuid() string {
	uuid := ""
	for i := 0; i < 32; i++ {
		rand := int(js.Global("Math").Call("random").Float()*16) | 0
		switch i {
		case 8, 12, 16, 20:
			uuid += "-"
		}
		switch i {
		case 12:
			uuid += "4"
		case 16:
			uuid += js.Global("Number").New(rand&3|8).Call("toString", 16).String()
		default:
			uuid += js.Global("Number").New(rand).Call("toString", 16).String()
		}
	}
	return uuid
}
func UuidGo() (uuid string) { //"pure" Go, but slower than native JS bindings
	for i := 0; i < 32; i++ {
		rand.Seed(time.Now().UnixNano() + int64(i))
		random := rand.Intn(16)
		switch i {
		case 8, 12, 16, 20:
			uuid += "-"
		}
		switch i {
		case 12:
			uuid += fmt.Sprintf("%X", 4)
		case 16:
			uuid += fmt.Sprintf("%X", random&3|8)
		default:
			uuid += fmt.Sprintf("%X", random)
		}
	}
	return
}

//handlebar templates
type Handlebar struct {
	js.Object
}

func CompileHandlebar(template string) *Handlebar {
	h := js.Global("Handlebars").Call("compile", template)
	return &Handlebar{h}
}
func RenderHandlebar(hb *Handlebar, i interface{}) string {
	return hb.Object.Invoke(i).String()
}
func RegisterHandlebarsHelper() {
	fn := func(a, b, options js.Object) js.Object {
		this := js.This()
		if a.String() == b.String() {
			return options.Call("fn", this)
		} else {
			return options.Call("inverse", this)
		}
	}
	js.Global("Handlebars").Call("registerHelper", "eq", fn)
}

//router (Director.js)
type Router struct {
	js.Object
}

func NewRouter() Router {
	return Router{Object: js.Global("Router").New()}
}
func (r Router) On(path string, handler func(string)) {
	r.Call("on", path, handler)
}

func (r Router) Init(path string) {
	r.Call("init", path)
}
