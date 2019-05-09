package models

// generated code, do not touch!
// generated at timestamp 2019-05-09 21:16:06.578661585 &#43;0000 UTC m=&#43;0.002796116

type PikabuUserSignupTimestampDistribution_86400 struct {
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int64         `sql:",notnull" json:"value"`
}

type PikabuUserLastUpdateTimestampDistribution_86400 struct {
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int64         `sql:",notnull" json:"value"`
}

type PikabuUserNextUpdateTimestampDistribution_86400 struct {
	Timestamp TimestampType `sql:",pk,notnull" json:"timestamp" api:"order,filter"`
	Value     int64         `sql:",notnull" json:"value"`
}

// distribution table name: base table name, base column name, distribution table
var GeneratedDistributionFields = map[string]struct {
	BaseTableName, BaseColumnName string
	DistributionTable             interface{}
	BucketSize                    int
}{
	"PikabuUserSignupTimestampDistribution_86400":     {"PikabuUser", "SignupTimestamp", &PikabuUserSignupTimestampDistribution_86400{}, 86400},
	"PikabuUserLastUpdateTimestampDistribution_86400": {"PikabuUser", "LastUpdateTimestamp", &PikabuUserLastUpdateTimestampDistribution_86400{}, 86400},
	"PikabuUserNextUpdateTimestampDistribution_86400": {"PikabuUser", "NextUpdateTimestamp", &PikabuUserNextUpdateTimestampDistribution_86400{}, 86400},
}
var GeneratedDistributionFieldsAPI = map[string]interface{}{
	"PikabuUserSignupTimestampDistribution_86400":     []PikabuUserSignupTimestampDistribution_86400{},
	"PikabuUserLastUpdateTimestampDistribution_86400": []PikabuUserLastUpdateTimestampDistribution_86400{},
	"PikabuUserNextUpdateTimestampDistribution_86400": []PikabuUserNextUpdateTimestampDistribution_86400{},
}

func init() {
	for _, item := range []interface{}{
		&PikabuUserSignupTimestampDistribution_86400{},
		&PikabuUserLastUpdateTimestampDistribution_86400{},
		&PikabuUserNextUpdateTimestampDistribution_86400{},
	} {
		Tables = append(Tables, item)
	}
}
