//+build ignore

package main

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/DevAlone/parse_pikabu/models"
)

type FieldT struct {
	Field      reflect.StructField
	BucketSize int
}

var headerTemplate = template.Must(template.New("").Parse(`package models
// generated code, do not touch!
// generated at timestamp {{ .Timestamp }}

`))

func main() {
	f, err := os.Create("models/models_distributions_fields.generated.go")
	handleErr(err)
	defer func() { _ = f.Close() }()

	err = headerTemplate.Execute(f, struct {
		Timestamp time.Time
	}{
		Timestamp: time.Now(),
	})
	handleErr(err)

	generatedTables := []struct {
		BaseTableName      string
		BaseColumnName     string
		GeneratedTableName string
		BucketSize         int
	}{}

	for _, table := range models.Tables {
		fields := []FieldT{}

		v := reflect.TypeOf(table).Elem()
		typeName := v.Name()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag, found := field.Tag.Lookup("gen_distributions")
			if !found {
				continue
			}
			if len(tag) == 0 {
				fmt.Printf("%v %v\n", table, field.Name)
				panic(0)
			}
			bucketSize, err := strconv.Atoi(tag)
			handleErr(err)
			fields = append(fields, FieldT{
				Field:      field,
				BucketSize: bucketSize,
			})
		}

		if len(fields) > 0 {
			for _, field := range fields {
				baseTableNameCamelCase := typeName
				// baseTableNameSnakeCase := strcase.ToSnake(baseTableNameCamelCase) + "s"
				baseColumnName := field.Field.Name
				generatedTableNameCamelCase := baseTableNameCamelCase + baseColumnName + "Distribution_" + fmt.Sprint(field.BucketSize)
				err := renderModelCode(f, generatedTableNameCamelCase, "int64")
				handleErr(err)
				generatedTables = append(generatedTables, struct {
					BaseTableName      string
					BaseColumnName     string
					GeneratedTableName string
					BucketSize         int
				}{
					BaseTableName:      baseTableNameCamelCase,
					BaseColumnName:     baseColumnName,
					GeneratedTableName: generatedTableNameCamelCase,
					BucketSize:         field.BucketSize,
				})
			}
		}
	}

	generatedDistributionFieldsCode := `
// distribution table name: base table name, base column name, distribution table
var GeneratedDistributionFields = map[string]struct{BaseTableName, BaseColumnName string; DistributionTable interface{}; BucketSize int}{
`
	for _, table := range generatedTables {
		generatedDistributionFieldsCode += `"` + table.GeneratedTableName + `": {"` +
			table.BaseTableName + `", "` + table.BaseColumnName + `", &` + table.GeneratedTableName + "{}, " + fmt.Sprint(table.BucketSize) + "},\n"
	}
	generatedDistributionFieldsCode += "}"

	_, err = f.WriteString(generatedDistributionFieldsCode)
	handleErr(err)

	// for api
	fieldsDistributionAPITablesMapCode := `
var GeneratedDistributionFieldsAPI = map[string]interface{}{
`
	for _, table := range generatedTables {
		fieldsDistributionAPITablesMapCode += `"` + table.GeneratedTableName + `": []` + table.GeneratedTableName + "{},\n"
	}

	fieldsDistributionAPITablesMapCode += `}
`

	_, err = f.WriteString(fieldsDistributionAPITablesMapCode)
	handleErr(err)

	// generate init
	_, err = f.WriteString(getInit(generatedTables))
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func renderModelCode(file *os.File, modelName string, valueFieldType string) error {
	return template.Must(template.New("").Parse(`
{{ $backtick := "`+"`"+`" }}

type {{.ModelName}} struct {
	Timestamp TimestampType {{ $backtick }}sql:",pk,notnull" json:"timestamp" api:"order,filter"{{ $backtick }}
	Value {{.ValueType}} {{ $backtick }}sql:",notnull" json:"value"{{ $backtick }}
}
`)).Execute(file, struct {
		ModelName string
		ValueType string
	}{
		ModelName: modelName,
		ValueType: valueFieldType,
	})
}

func getInit(tableNames []struct {
	BaseTableName      string
	BaseColumnName     string
	GeneratedTableName string
	BucketSize         int
}) string {
	result := `func init() {
	for _, item := range []interface{}{
`

	for _, table := range tableNames {
		result += "&" + table.GeneratedTableName + "{},\n"
	}

	result += `} {
		Tables = append(Tables, item)
	}
}`
	return result
}
