package repository

import (
	"strings"
)

func GetVariantSKU(productSku string, attributeNames []string) string {
	sku := productSku
	for _, name := range attributeNames {
		if len(name) < 3 {
			sku += "-" + name
		} else {
			sku += "-" + strings.ToUpper(name[0:2]) + strings.ToUpper(name[len(name)-1:])
		}
	}
	return sku
}

func GetVariantSKUWithAttributeNames(productSku string, attrs []AttributeValue) string {

	attributeNames := make([]string, 0)
	for _, attr := range attrs {
		attributeNames = append(attributeNames, attr.Value)

	}
	return GetVariantSKU(productSku, attributeNames)
}
