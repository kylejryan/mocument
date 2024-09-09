package utils

type TypeGeneratorFunc func(FieldSchema) interface{}

var typeGenerators = map[string]TypeGeneratorFunc{
	"string":  generateRandomString,
	"float":   generateRandomFloat,
	"integer": generateRandomInteger,
}

func RegisterTypeGenerator(typeName string, generatorFunc TypeGeneratorFunc) {
	typeGenerators[typeName] = generatorFunc
}

func GenerateDataWithPlugins(field FieldSchema) interface{} {
	if generator, ok := typeGenerators[field.Type]; ok {
		return generator(field)
	}
	return nil
}
