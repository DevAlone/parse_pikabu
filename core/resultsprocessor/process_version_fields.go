package resultsprocessor

import (
	"reflect"

	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// GetModelsChangedFields -
func GetModelsChangedFields(
	oldModelPtr interface{},
	newModelPtr interface{},
) ([]string, error) {
	result := []string{}

	if reflect.TypeOf(oldModelPtr) != reflect.TypeOf(newModelPtr) {
		return nil, errors.New("types should be equal")
	}

	oldModel := reflect.ValueOf(oldModelPtr).Elem()
	newModel := reflect.ValueOf(newModelPtr).Elem()

	oldID := oldModel.FieldByName("PikabuID").Uint()
	newID := newModel.FieldByName("PikabuID").Uint()

	if oldID != newID {
		return nil, errors.New("ids should be equal")
	}

	oldModelType := reflect.TypeOf(oldModel.Interface())

	for i := 0; i < oldModelType.NumField(); i++ {
		fieldType := oldModelType.Field(i)

		_, isVersionedField := fieldType.Tag.Lookup("gen_versions")
		if !isVersionedField {
			continue
		}

		oldField := oldModel.FieldByName(fieldType.Name)
		newField := newModel.FieldByName(fieldType.Name)

		if reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			continue
		}

		result = append(result, fieldType.Name)
	}

	return result, nil
}

func processModelFieldsVersions(
	transaction *pg.Tx,
	oldModelPtr interface{},
	newModelPtr interface{},
	parsingTimestamp models.TimestampType,
) (bool, error) {
	// TODO: consider adding lock here
	wasDataChanged := false

	oldModel := reflect.ValueOf(oldModelPtr).Elem()
	newModel := reflect.ValueOf(newModelPtr).Elem()

	oldID := oldModel.FieldByName("PikabuID").Uint()
	newID := newModel.FieldByName("PikabuID").Uint()

	if oldID != newID {
		return false, errors.New("ids should be equal")
	}

	addedTimestamp := models.TimestampType(oldModel.FieldByName("AddedTimestamp").Int())
	lastUpdateTimestamp := models.TimestampType(oldModel.FieldByName("LastUpdateTimestamp").Int())

	if parsingTimestamp <= lastUpdateTimestamp {
		return false, OldParserResultError{}
	}

	oldModelType := reflect.TypeOf(oldModel.Interface())

	changedFields, err := GetModelsChangedFields(oldModelPtr, newModelPtr)
	if err != nil {
		return false, err
	}

	for _, fieldName := range changedFields {
		// fieldType := oldModelType.Field(i)

		oldField := oldModel.FieldByName(fieldName)
		newField := newModel.FieldByName(fieldName)

		if reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			continue
		}

		wasDataChanged = true

		// generate versions
		versionTable := models.FieldsVersionTablesMap[oldModelType.Name()+fieldName+"Version"]

		insertVersion := func(
			timestamp models.TimestampType,
			value reflect.Value,
			ignoreIfExists bool,
		) error {
			e := reflect.ValueOf(versionTable).Elem()
			e.FieldByName("ItemId").SetUint(oldID)
			e.FieldByName("Timestamp").Set(reflect.ValueOf(timestamp))
			e.FieldByName("Value").Set(value)

			var err error
			if ignoreIfExists {
				var q *orm.Query
				if transaction == nil {
					q = models.Db.Model(versionTable)
				} else {
					q = transaction.Model(versionTable)
				}
				_, err = q.
					OnConflict("DO NOTHING").
					Insert(versionTable)
			} else {
				if transaction == nil {
					err = models.Db.Insert(versionTable)
				} else {
					err = transaction.Insert(versionTable)
				}
			}
			if err != nil {
				return errors.New(err)
			}

			return nil
		}

		var q *orm.Query
		if transaction == nil {
			q = models.Db.Model(versionTable)
		} else {
			q = transaction.Model(versionTable)
		}
		count, err := q.Where("item_id = ?", oldID).Count()
		if err != nil {
			return false, errors.New(err)
		}

		if count == 0 {
			err := insertVersion(
				addedTimestamp,
				oldField,
				false)
			if err != nil {
				return false, err
			}
		}

		err = insertVersion(
			lastUpdateTimestamp,
			oldField,
			true)
		if err != nil {
			return false, err
		}

		err = insertVersion(
			parsingTimestamp,
			newField,
			false)
		if err != nil {
			return false, err
		}

		// set the field
		if !oldField.CanSet() {
			panic(errors.Errorf("field %v from model %v cannot be set", oldField, oldModel))
		}
		oldField.Set(newField)
	}

	return wasDataChanged, nil
}
