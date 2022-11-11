package object

type ObjectType int32

const (
	ClosureType ObjectType = iota
	DoubleType
	FunPtrType
	StringType
	FloatType
	BoolType
	LongType
	IntType
)

type Object interface {
	String() string
	Type() ObjectType
}
