package method

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func GenerateRandomSourcesTypeAndName(sourceType string) (string, string, string) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource := fmt.Sprintf("%s.%s", sourceType, name)
	dataSource := fmt.Sprintf("data.%s.%s", sourceType, name)
	return resource, dataSource, name
}
