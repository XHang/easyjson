package gen

import (
	"fmt"
	"reflect"
	"strings"
)

const modelTpl = `
// New{{.Typ}} 创建新的{{.Typ}}
func New{{.Typ}}() *{{.Typ}} {
	v := &{{.Typ}}{}
	v.ResetFieldMark()
	return v
}

// NewItems 创建对应的切片指针对象，配合 go-component/orm 组件
func (*{{.Typ}}) NewItems() interface{} {
	items := new([]{{.Typ}})
	*items = make([]{{.Typ}}, 0)
	return items
}

// Get{{.Typ}}Slice 从 ModelList 中获取items切片对象
func Get{{.Typ}}Slice(ml *orm.ModelList) (items []{{.Typ}}, ok bool) {
	var val *[]{{.Typ}}
	val, ok = ml.Items.(*[]{{.Typ}})
	if !ok {
		return
	}

	items := *val
	return
}

// BindReader 从 reader 中读取内容映射绑定
func (v *{{.Typ}}) BindReader(reader io.Reader) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return v.UnmarshalJSON(body)
}
`

const fieldTpl = `
// {{.FieldName}}Mark {{.FieldName}}是否已赋值（赋值标识）
func (v *{{.Typ}}) {{.FieldName}}Mark() bool {
	return v.HasFieldMark("{{.FieldName}}")
}

// Set{{.FieldName}} 设置{{.FieldName}}}的值，并将赋值标识设为:true
func (v *{{.Typ}}) Set{{.FieldName}}(val {{.FieldTyp}}) {
	v.{{.FieldName}} = val
	v.SetFieldMark("{{.FieldName}}")
}
`

const resetEOFTpl = `	if in.Error() == io.EOF {
		in.ResetError(nil)
	}`

const resetErrorTpl = `		if err := in.Error(); err != nil {
			msg := ""
			if strings.Contains(err.Error(), "unknown field") {
				msg = "不存在的参数：" + key
			} else {
				msg = key + "格式错误"
			}
			in.ResetError(&jlexer.LexerError{
				Data: msg,
			})
			return
		}`

func (g *Generator) genModel(t reflect.Type, fs []reflect.StructField, typ string) {
	modelStr := modelTpl

	for _, f := range fs {
		// jsonName := g.fieldNamer.GetJSONFieldName(t, f)
		tags := parseFieldTags(f)

		if tags.omit {
			continue
		}

		fieldStr := strings.Replace(fieldTpl, "{{.FieldName}}", f.Name, -1)
		fieldStr = strings.Replace(fieldStr, "{{.FieldTyp}}", g.getType(f.Type), -1)

		modelStr += fieldStr
	}

	modelStr = strings.Replace(modelStr, "{{.Typ}}", typ, -1)

	fmt.Fprintln(g.out, modelStr)
}
