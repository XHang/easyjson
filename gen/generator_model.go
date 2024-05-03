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
	v.ResetMark()
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
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return v.UnmarshalJSON(body)
}

// DBFields 数据库字段
func (v *{{.Typ}}) DBFields() []string {
	fields := "{{.DBFields}}"
	return strings.Split(fields, ",")
}

// DBFieldJSONTags 数据库字段对应的json tag name
func (v *{{.Typ}}) DBFieldJSONTags() []string {
	names := "{{.DBFieldJSONTags}}"
	return strings.Split(names, ",")
}

// DBFieldsIndex 数据库字段索引 大小写敏感
func (v *{{.Typ}}) DBFieldsIndex() map[string]struct{} {
	fields := v.DBFields()
	l := len(fields)
	index := make(map[string]struct{}, l)
	for i := 0; i < l; i++{
		index[fields[i]] = struct{}{}
	}
	return index
}

// DBFieldJSONNamesIndex 数据库字段对应的json name 索引 大小写敏感
func (v *{{.Typ}}) DBFieldJSONNamesIndex() map[string]struct{} {
	names := v.DBFieldJSONTags()
	l := len(names)
	index := make(map[string]struct{}, l)
	for i := 0; i < l; i++{
		index[names[i]] = struct{}{}
	}
	return index
}

// DBFieldJSONNameIndex 数据库字段为Key, json name 为 value 索引; 大小写敏感
func (v *{{.Typ}}) DBFieldJSONNameIndex() map[string]string {
	fields := v.DBFields()
	names := v.DBFieldJSONTags()
	l := len(names)
	index := make(map[string]string, l)
	for i := 0; i < l; i++{
		index[fields[i]] = names[i]
	}
	return index
}

// JSONNameDBFieldIndex json name 为Key, 数据库字段 为 value 索引; 大小写敏感
func (v *{{.Typ}}) JSONNameDBFieldIndex() map[string]string {
	fields := v.DBFields()
	names := v.DBFieldJSONTags()
	l := len(names)
	index := make(map[string]string, l)
	for i := 0; i < l; i++{
		index[names[i]] = fields[i]
	}
	return index
}
`

const fieldTpl = `
// Is{{.FieldName}}Mark {{.FieldName}}是否已赋值（赋值标识）
func (v *{{.Typ}}) Is{{.FieldName}}Mark() bool {
	return v.HasPropertyMark(v.Get{{.FieldName}}MarkKey())
}

// Set{{.FieldName}} 设置{{.FieldName}}}的值，并将赋值标识设为:true
func (v *{{.Typ}}) Set{{.FieldName}}(val {{.FieldTyp}}) {
	v.{{.FieldName}} = val
	v.Mark{{.FieldName}}()
}

// UnMark{{.FieldName}} 取消{{.FieldName}}}的赋值标识，设为:false
func (v *{{.Typ}}) UnMark{{.FieldName}}() {
	v.SetPropertyMark(v.Get{{.FieldName}}MarkKey(), false)
	v.SetFieldMark(v.Get{{.FieldName}}MarkKey(), false)
}

// Mark{{.FieldName}} 设置{{.FieldName}}}的赋值标识，设为:true
func (v *{{.Typ}}) Mark{{.FieldName}}() {
	v.SetPropertyMark(v.Get{{.FieldName}}MarkKey(), true)
	v.SetFieldMark(v.Get{{.FieldName}}MarkKey(), true)
}

// Get{{.FieldName}}MarkKey 获取MarkKey
func (v *{{.Typ}}) Get{{.FieldName}}MarkKey() string {
	return "{{.MarkKey}}"
}
`

const fieldIgnoreTpl = `
// Is{{.FieldName}}Mark {{.FieldName}}是否已赋值（赋值标识）
func (v *{{.Typ}}) Is{{.FieldName}}Mark() bool {
	return v.HasPropertyMark(v.Get{{.FieldName}}MarkKey())
}

// Set{{.FieldName}} 设置{{.FieldName}}}的值，并将赋值标识设为:true
func (v *{{.Typ}}) Set{{.FieldName}}(val {{.FieldTyp}}) {
	v.{{.FieldName}} = val
	v.Mark{{.FieldName}}()
}

// UnMark{{.FieldName}} 取消{{.FieldName}}}的赋值标识，设为:false
func (v *{{.Typ}}) UnMark{{.FieldName}}() {
	v.SetPropertyMark(v.Get{{.FieldName}}MarkKey(), false)
	//v.SetFieldMark(v.Get{{.FieldName}}MarkKey(), false) tag的xorm为"-"
}

// Mark{{.FieldName}} 设置{{.FieldName}}}的赋值标识，设为:true
func (v *{{.Typ}}) Mark{{.FieldName}}() {
	v.SetPropertyMark(v.Get{{.FieldName}}MarkKey(), true)
	//v.SetFieldMark(v.Get{{.FieldName}}MarkKey(), true) tag的xorm为"-"
}

// Get{{.FieldName}}FieldMarkKey 获取MarkKey
func (v *{{.Typ}}) Get{{.FieldName}}MarkKey() string {
	return "{{.MarkKey}}"
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

	dbFields := make([]string, 0)
	dbFieldJsonTags := make([]string, 0)
	for _, f := range fs {
		// jsonName := g.fieldNamer.GetJSONFieldName(t, f)

		// 2019-10-16 注释
		// tags := parseFieldTags(f)
		// if tags.omit {
		// 	continue
		// }

		var fieldStr string

		xormField, _ := parseXormFieldName(f)
		if xormField == "-" { // xorm:"-" 忽略fileMarks操作标识

			fieldStr = strings.Replace(fieldIgnoreTpl, "{{.FieldName}}", f.Name, -1)
			fieldStr = strings.Replace(fieldStr, "{{.MarkKey}}", f.Name, -1)
			fieldStr = strings.Replace(fieldStr, "{{.FieldTyp}}", g.getType(f.Type), -1)

		} else {
			if xormField == "" {
				xormField = f.Name // 没有设置xorm字段名，取结构体属性名
			}

			dbFields = append(dbFields, xormField)
			dbFieldJsonTags = append(dbFieldJsonTags, g.fieldNamer.GetJSONFieldName(t, f))
			fieldMarkKey := xormField

			fieldStr = strings.Replace(fieldTpl, "{{.FieldName}}", f.Name, -1)     // {{.FieldName}}: 对外显示的方法(属性)名
			fieldStr = strings.Replace(fieldStr, "{{.MarkKey}}", fieldMarkKey, -1) // {{.MarkKey}}: fieldMark实际的keyName
			fieldStr = strings.Replace(fieldStr, "{{.FieldTyp}}", g.getType(f.Type), -1)
		}

		modelStr += fieldStr
	}

	modelStr = strings.Replace(modelStr, "{{.Typ}}", typ, -1)
	modelStr = strings.Replace(modelStr, "{{.DBFields}}", strings.Join(dbFields, ","), -1)
	modelStr = strings.Replace(modelStr, "{{.DBFieldJSONTags}}", strings.Join(dbFieldJsonTags, ","), -1)

	fmt.Fprintln(g.out, modelStr)
}
