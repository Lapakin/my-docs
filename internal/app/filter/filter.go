package filter

type Filters map[string]string

// Operator defines logical operators for combining conditions
type Operator string

const (
	And Operator = "AND"
	Or  Operator = "OR"
)

// ConditionOperator defines comparison operators
type ConditionOperator string

const (
	Equals             ConditionOperator = "="
	NotEquals          ConditionOperator = "!="
	GreaterThan        ConditionOperator = ">"
	GreaterThanOrEqual ConditionOperator = ">="
	LessThan           ConditionOperator = "<"
	LessThanOrEqual    ConditionOperator = "<="
	Like               ConditionOperator = "LIKE"
	ILike              ConditionOperator = "ILIKE"
	In                 ConditionOperator = "IN"
	NotIn              ConditionOperator = "NOT IN"
	IsNull             ConditionOperator = "IS NULL"
	IsNotNull          ConditionOperator = "IS NOT NULL"
)

// Condition represents a single filter condition
type Condition struct {
	Name     string
	Column   string
	Operator ConditionOperator
}

type Conditions []*Condition

// Where represents a group of conditions combined by an operator
type Where struct {
	Operator   Operator
	Conditions Conditions
}

type Wheres []*Where
