package old_models

type FieldVersionBase struct {
	Timestamp int32  `sql:",pk,notnull"`
	ItemId    uint64 `sql:",pk,notnull"`
}

type Int32FieldVersion struct {
	FieldVersionBase
	Value int32 `sql:",notnull"`
}
type Int64FieldVersion struct {
	FieldVersionBase
	Value int64 `sql:",notnull"`
}
type Uint64FieldVersion struct {
	FieldVersionBase
	Value uint64 `sql:",notnull"`
}
type TextFieldVersion struct {
	FieldVersionBase
	Value string `sql:",notnull"`
}
type BoolFieldVersion struct {
	FieldVersionBase
	Value bool `sql:",notnull"`
}
