package zpa

/*
func TestAccDataSourceInspectionProfile_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName, variable.InspectionProfileDescription, variable.InspectionProfileParanoia),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "paranoia_level", resourceTypeAndName, "paranoia_level"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "predefined_controls.#", "7"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
*/
