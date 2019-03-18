package gen

import (
	"fmt"
	"reflect"
	"strings"
)

/**
orm.Model 自定义内容
*/

const modelTpl = `
// init 注册数据引擎
func init() {
	orm.RegisterAlias(&{{.Typ}}{})
}

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

	items = *val
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
// Is{{.FieldName}}Mark {{.FieldName}}是否已赋值（赋值标识）
func (v *{{.Typ}}) Is{{.FieldName}}Mark() bool {
	return v.HasFieldMark(v.Get{{.FieldName}}FieldMarkKey())
}

// Set{{.FieldName}} 设置{{.FieldName}}}的值，并将赋值标识设为:true
func (v *{{.Typ}}) Set{{.FieldName}}(val {{.FieldTyp}}) {
	v.{{.FieldName}} = val
	v.Mark{{.FieldName}}()
}

// UnMark{{.FieldName}} 取消{{.FieldName}}}的赋值标识，设为:false
func (v *{{.Typ}}) UnMark{{.FieldName}}() {
	v.SetFieldMark(v.Get{{.FieldName}}FieldMarkKey(), false)
}

// Mark{{.FieldName}} 设置{{.FieldName}}}的赋值标识，设为:true
func (v *{{.Typ}}) Mark{{.FieldName}}() {
	v.SetFieldMark(v.Get{{.FieldName}}FieldMarkKey(), true)
}

// Mark{{.FieldName}}FieldMarkKey 获取FieldMarkKey
func (v *{{.Typ}}) Get{{.FieldName}}FieldMarkKey() string {
	return "{{.FieldMarkKey}}"
}
`

const fieldIgnoreTpl = `
// Is{{.FieldName}}Mark {{.FieldName}}是否已赋值（赋值标识）
func (v *{{.Typ}}) Is{{.FieldName}}Mark() bool {
	return false
}

// Set{{.FieldName}} 设置{{.FieldName}}}的值，并将赋值标识设为:true
func (v *{{.Typ}}) Set{{.FieldName}}(val {{.FieldTyp}}) {
	v.{{.FieldName}} = val
	v.Mark{{.FieldName}}()
}

// UnMark{{.FieldName}} 取消{{.FieldName}}}的赋值标识，设为:false(xorm:"-"时，该方法不操作赋值标识)
func (v *{{.Typ}}) UnMark{{.FieldName}}() {
	//v.fieldMark[v.Get{{.FieldName}}FieldMarkKey()] = false
}

// Mark{{.FieldName}} 设置{{.FieldName}}}的赋值标识，设为:true(xorm:"-"时，该方法不操作赋值标识)
func (v *{{.Typ}}) Mark{{.FieldName}}() {
	//v.fieldMark[v.Get{{.FieldName}}FieldMarkKey()] = true
}

// Mark{{.FieldName}}FieldMarkKey 获取FieldMarkKey(xorm:"-"时，FieldMarkKey为空)
func (v *{{.Typ}}) Get{{.FieldName}}FieldMarkKey() string {
	return ""
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

		var fieldStr string

		xormField, _ := parseXormFieldName(f)
		if xormField == "-" { // xorm:"-" 忽略fileMarks操作标识

			fieldStr = strings.Replace(fieldIgnoreTpl, "{{.FieldName}}", f.Name, -1)
			fieldStr = strings.Replace(fieldStr, "{{.FieldTyp}}", g.getType(f.Type), -1)

		} else {
			if xormField == "" {
				xormField = f.Name // 没有设置xorm字段名，取结构体属性名
			}

			fieldMarkKey := xormField

			fieldStr = strings.Replace(fieldTpl, "{{.FieldName}}", f.Name, -1)          // {{.FieldName}}: 对外显示的方法(属性)名
			fieldStr = strings.Replace(fieldStr, "{{.FieldMarkKey}}", fieldMarkKey, -1) // {{.FieldMarkKey}}: fieldMark实际的keyName
			fieldStr = strings.Replace(fieldStr, "{{.FieldTyp}}", g.getType(f.Type), -1)
		}

		modelStr += fieldStr
	}

	modelStr = strings.Replace(modelStr, "{{.Typ}}", typ, -1)

	fmt.Fprintln(g.out, modelStr)
}
