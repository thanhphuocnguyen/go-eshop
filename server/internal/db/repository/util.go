package repository

import (
	"context"
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

func GetVariantSKUWithAttributeNames(q *Queries, c context.Context, productSku string, attributeValueIDs []int32) (string, error) {
	attrs, err := q.GetAttributeValuesByIDs(c, attributeValueIDs)
	if err != nil {
		return "", err
	}
	attributeNames := make([]string, 0)
	for _, attr := range attrs {
		if attr.DisplayValue.Valid {
			attributeNames = append(attributeNames, attr.DisplayValue.String)
		} else {
			attributeNames = append(attributeNames, attr.Value)
		}
	}
	return GetVariantSKU(productSku, attributeNames), nil
}
