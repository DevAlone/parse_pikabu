package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg/orm"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/parser"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func processFilter(req *orm.Query, resultType reflect.Type, filter string) (*orm.Query, error) {
	filter = strings.TrimSpace(filter)
	if len(filter) == 0 {
		return req, nil
	}

	src := common.NewTextSource(filter)

	expression, errs := parser.Parse(src)
	if len(errs.GetErrors()) != 0 {
		return req, errors.Errorf(errs.ToDisplayString())
	}

	typeProvider := types.NewRegistry()
	env := checker.NewStandardEnv(packages.DefaultPackage, typeProvider)
	// declare fields to filter on
	err := env.Add(getFieldsToFilter(resultType)...)
	if err != nil {
		return req, err
	}

	err = env.Add(getFilterFunctions()...)
	if err != nil {
		return req, err
	}

	c, errs := checker.Check(expression, src, env)
	if len(errs.GetErrors()) != 0 {
		return req, errors.Errorf(errs.ToDisplayString())
	}

	e := c.GetExpr()

	params := []interface{}{}
	sqlExpression, err := makeSQLExpression(e, &params)
	if err != nil {
		return req, err
	}

	return req.Where(sqlExpression, params...), nil
}

func makeSQLExpression(e *expr.Expr, params *[]interface{}) (string, error) {
	switch v := e.ExprKind.(type) {
	case *expr.Expr_CallExpr:
		if len(v.CallExpr.Args) != 2 {
			logger.Log.Errorf("wrong number of arguments: %v", v.CallExpr.Args)
			return "", errors.Errorf("some very bad shit happened")
		}

		left, err := makeSQLExpression(v.CallExpr.Args[0], params)
		if err != nil {
			return "", err
		}
		right, err := makeSQLExpression(v.CallExpr.Args[1], params)
		if err != nil {
			return "", err
		}

		function, err := celFunctionToSQL(v.CallExpr.Function)
		if err != nil {
			return "", err
		}

		return " " + left + function + right + " ", nil
	case *expr.Expr_IdentExpr:
		return ` "` + v.IdentExpr.Name + `" `, nil
	case *expr.Expr_ConstExpr:
		switch c := v.ConstExpr.ConstantKind.(type) {
		case *expr.Constant_Int64Value:
			// return " " + fmt.Sprint(c.Int64Value) + " ", nil
			*params = append(*params, c.Int64Value)
			return " ?" + fmt.Sprint(len(*params)-1) + " ", nil
		case *expr.Constant_Uint64Value:
			*params = append(*params, c.Uint64Value)
			return " ?" + fmt.Sprint(len(*params)-1) + " ", nil
		case *expr.Constant_StringValue:
			*params = append(*params, c.StringValue)
			return " ?" + fmt.Sprint(len(*params)-1) + " ", nil
		case *expr.Constant_BoolValue:
			*params = append(*params, c.BoolValue)
			return " ?" + fmt.Sprint(len(*params)-1) + " ", nil
		default:
			bytes, _ := json.Marshal(c)
			logger.Log.Debugf("unknown kind of constant: %v", string(bytes))
			return "", errors.Errorf("unknown(or not allowed) kind of constant")
		}
	default:
		bytes, _ := json.Marshal(v)
		logger.Log.Debugf("unknown kind of expression: %v", string(bytes))
		return "", errors.Errorf("unknown(or not allowed) kind of expression")
	}
}

func celFunctionToSQL(function string) (string, error) {
	if len(function) < 3 {
		logger.Log.Errorf("unknown cel function %v", function)
		return "", errors.Errorf("some very bad shit happened")
	}

	function = strings.TrimPrefix(function, "_")
	function = strings.TrimSuffix(function, "_")
	function = strings.ToLower(function)

	switch function {
	case "||":
		return " OR ", nil
	case "&&":
		return " AND ", nil
	case "==":
		return "=", nil
	case ">", "<", "!=", ">=", "<=":
		return function, nil
	case "ilike":
		return " ILIKE ", nil
	default:
		logger.Log.Debugf("unknown function: %v", function)
		return "", errors.Errorf("you're not allowed to use function \"%v\"", function)
	}
}

func getFieldsToFilter(modelType reflect.Type) []*exprpb.Decl {
	result := []*exprpb.Decl{}

	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		if tag, found := fieldType.Tag.Lookup("api"); found {
			for _, item := range strings.Split(tag, ",") {
				item = strings.TrimSpace(item)
				if item == "filter" {
					decl, err := fieldToDecl(fieldType)
					if err != nil {
						logger.Log.Errorf("Forgot to create decl for type %v", fieldType.Type)
					} else {
						result = append(result, decl)
					}
				}
			}
		}
	}

	return result
}

func getFilterFunctions() []*exprpb.Decl {
	result := []*exprpb.Decl{}

	result = append(result, decls.NewFunction("ilike", decls.NewOverload(
		"field, string",
		[]*exprpb.Type{
			decls.String,
			decls.String,
		},
		decls.Bool,
	)))

	return result
}

func fieldToDecl(fieldType reflect.StructField) (*exprpb.Decl, error) {
	fieldAPIName := ""
	if jsonTag, found := fieldType.Tag.Lookup("json"); found {
		jsonName := strings.Split(jsonTag, ",")[0]
		jsonName = strings.TrimSpace(jsonName)
		if len(jsonName) > 0 {
			fieldAPIName = jsonName
		}
	}
	if len(fieldAPIName) == 0 {
		fieldAPIName = fieldType.Name
	}

	switch fieldType.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decls.NewIdent(fieldAPIName, decls.Int, nil), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// return decls.NewIdent(fieldAPIName, decls.Uint, nil), nil
		return decls.NewIdent(fieldAPIName, decls.Int, nil), nil
	case reflect.String:
		return decls.NewIdent(fieldAPIName, decls.String, nil), nil
	case reflect.Bool:
		return decls.NewIdent(fieldAPIName, decls.Bool, nil), nil
	}

	return nil, errors.Errorf("forgot to create decl for field %v", fieldType)
}
