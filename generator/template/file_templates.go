package template

var autoGenWarningTemplate = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior  
// and will be lost if the code is regenerated
//

`

var tableSQLBuilderTemplate = ` 
{{define "column-list" -}}
	{{- range $i, $c := . }}
{{- $field := columnField $c}}
		{{- if gt $i 0 }}, {{end}}{{$field.Name}}Column
	{{- end}}
{{- end}}

package {{package}}

import (
	"github.com/go-jet/jet/v2/{{dialect.PackageName}}"
)

var {{tableTemplate.InstanceName}} = new{{tableTemplate.TypeName}}("{{schemaName}}", "{{.Name}}", "")

type {{structImplName}} struct {
	{{dialect.PackageName}}.Table
	
	// Columns
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
	{{$field.Name}} {{dialect.PackageName}}.Column{{$field.Type}} {{- if $c.Comment }} // {{$c.Comment}} {{end}}
{{- end}}

	AllColumns     {{dialect.PackageName}}.ColumnList
	MutableColumns {{dialect.PackageName}}.ColumnList
}

type {{tableTemplate.TypeName}} struct {
	{{structImplName}}

	{{toUpper insertedRowAlias}} {{structImplName}}
}

// AS creates new {{tableTemplate.TypeName}} with assigned alias
func (a {{tableTemplate.TypeName}}) AS(alias string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new {{tableTemplate.TypeName}} with assigned schema name
func (a {{tableTemplate.TypeName}}) FromSchema(schemaName string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new {{tableTemplate.TypeName}} with assigned table prefix
func (a {{tableTemplate.TypeName}}) WithPrefix(prefix string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new {{tableTemplate.TypeName}} with assigned table suffix
func (a {{tableTemplate.TypeName}}) WithSuffix(suffix string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func new{{tableTemplate.TypeName}}(schemaName, tableName, alias string) *{{tableTemplate.TypeName}} {
	return &{{tableTemplate.TypeName}}{
		{{structImplName}}: new{{tableTemplate.TypeName}}Impl(schemaName, tableName, alias),
		{{toUpper insertedRowAlias}}:  new{{tableTemplate.TypeName}}Impl("", "{{insertedRowAlias}}", ""),
	}
}

func new{{tableTemplate.TypeName}}Impl(schemaName, tableName, alias string) {{structImplName}} {
	var (
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
		{{$field.Name}}Column = {{dialect.PackageName}}.{{$field.Type}}Column("{{$c.Name}}")
{{- end}}
		allColumns     = {{dialect.PackageName}}.ColumnList{ {{template "column-list" .Columns}} }
		mutableColumns = {{dialect.PackageName}}.ColumnList{ {{template "column-list" .MutableColumns}} }
	)

	return {{structImplName}}{
		Table: {{dialect.PackageName}}.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
		{{$field.Name}}: {{$field.Name}}Column,
{{- end}}

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
`

var tableSqlBuilderSetSchemaTemplate = `package {{package}}

// UseSchema changes all global tables/views with the value returned
// returned by calling FromSchema on them. Passing an empty string to this function
// will cause queries to be generated without any table/view alias.
func UseSchema(schema string) {
{{- range .}}
	{{ .InstanceName }} = {{ .InstanceName }}.FromSchema(schema)
{{- end}}
}
`

var tableModelFileTemplate = `package {{package}}

{{ with modelImports }}
import (
{{- range .}}
	"{{.}}"
{{- end}}
)
{{end}}

{{$modelTableTemplate := tableTemplate}}
type {{$modelTableTemplate.TypeName}} struct {
{{- range .Columns}}
{{- $field := structField .}}
	{{$field.Name}} {{$field.Type.Name}} ` + "{{$field.TagsString}}" + ` {{- if .Comment }} // {{.Comment}} {{end}}
{{- end}}
}

`

var enumSQLBuilderTemplate = `package {{package}}

import "github.com/go-jet/jet/v2/{{dialect.PackageName}}"

var {{enumTemplate.InstanceName}} = &struct {
{{- range $index, $value := .Values}}
	{{enumValueName $value}} {{dialect.PackageName}}.StringExpression
{{- end}}
} {
{{- range $index, $value := .Values}}
	{{enumValueName $value}}: {{dialect.PackageName}}.NewEnumValue("{{$value}}"),
{{- end}}
}
`

var enumModelTemplate = `package {{package}}
{{- $enumTemplate := enumTemplate}}

import "errors"

type {{$enumTemplate.TypeName}} string

const (
{{- range $_, $value := .Values}}
	{{valueName $value}} {{$enumTemplate.TypeName}} = "{{$value}}"
{{- end}}
)

func (e *{{$enumTemplate.TypeName}}) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
{{- range $_, $value := .Values}}
	case "{{$value}}":
		*e = {{valueName $value}}
{{- end}}
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for {{$enumTemplate.TypeName}} enum")
	}

	return nil
}

func (e {{$enumTemplate.TypeName}}) String() string {
	return string(e)
}

`
