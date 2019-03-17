// +build ignore

package main

import (
	"html/template"
	"os"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/iancoleman/strcase"
)

func main() {
	f, err := os.Create("models/models_versions_fields.generated.go")
	handleErr(err)
	defer f.Close()

	err = headerTemplate.Execute(f, struct {
		Timestamp time.Time
	}{
		Timestamp: time.Now(),
	})
	handleErr(err)

	generatedTablesNames := []string{}

	for _, table := range models.Tables {
		versionsFields := []reflect.StructField{}

		v := reflect.TypeOf(table).Elem()
		typeName := v.Name()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			_, found := field.Tag.Lookup("gen_versions")
			if found {
				versionsFields = append(versionsFields, field)
			}
		}

		if len(versionsFields) > 0 {
			for _, table := range getVersionsTables(typeName, versionsFields) {
				generatedTablesNames = append(generatedTablesNames, table.tableName)
				_, err := f.WriteString(table.tableCode)
				handleErr(err)
			}
		}
	}

	fieldsVersionTablesMapCode := `
var FieldsVersionTablesMap = map[string]interface{}{
`
	for _, tableName := range generatedTablesNames {
		fieldsVersionTablesMapCode += `"` + tableName + `": &` + tableName + "{},\n"
	}

	fieldsVersionTablesMapCode += `}
`

	_, err = f.WriteString(fieldsVersionTablesMapCode)
	handleErr(err)

	fieldsVersionAPITablesMapCode := `
var FieldsVersionAPITablesMap = map[string]interface{}{
`
	for _, tableName := range generatedTablesNames {
		fieldsVersionAPITablesMapCode += `"` + tableName + `": []` + tableName + "{},\n"
	}

	fieldsVersionAPITablesMapCode += `}
`

	_, err = f.WriteString(fieldsVersionAPITablesMapCode)
	handleErr(err)

	// generate init
	_, err = f.WriteString(getInit(generatedTablesNames))
	handleErr(err)
}

func getVersionsTables(typeName string, versionsFields []reflect.StructField) []struct {
	tableName string
	tableCode string
} {
	result := []struct {
		tableName string
		tableCode string
	}{}

	for _, field := range versionsFields {
		fieldName := field.Name

		var fieldTypeName string

		switch field.Type.Kind() {
		case reflect.Slice, reflect.Array:
			fieldTypeName = field.Type.String()
		default:
			fieldTypeName = field.Type.Name()
		}

		if strings.HasPrefix(fieldTypeName, "[]models.") {
			fieldTypeName = "[]" + fieldTypeName[9:]
		}

		name := typeName + fieldName + "Version"
		result = append(result, struct {
			tableName string
			tableCode string
		}{
			tableName: name,
			tableCode: getModelCode(name, fieldTypeName) + "\n",
		})
	}

	return result
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getStructMemberCode(memberName string, memberType string, isPk bool) string {
	result := "\t" + memberName + " " + memberType + " `sql:\""
	if isPk {
		result += ",pk"
	}
	apiName := strcase.ToSnake(memberName)
	result += ",notnull\" json:\"" + apiName + "\" api:\"" + "order,filter" + "\"`\n"
	return result
}

func getModelCode(tableName string, fieldType string) string {
	result := "type " + tableName + " struct {\n"
	result += getStructMemberCode("ItemId", "uint64", true)
	result += getStructMemberCode("Timestamp", "TimestampType", true)
	result += getStructMemberCode("Value", fieldType, false)
	result += "}\n"
	return result
}

var headerTemplate = template.Must(template.New("").Parse(`package models
// generated code, do not touch!
// generated at timestamp {{ .Timestamp }}

`))

func getInit(tableNames []string) string {
	result := `func init() {
	for _, item := range []interface{}{
`

	for _, tableName := range tableNames {
		result += "&" + tableName + "{},\n"
	}

	result += `} {
		Tables = append(Tables, item)
	}
}`
	return result
}
