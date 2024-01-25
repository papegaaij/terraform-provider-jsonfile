// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jsonfile_data.data", "value", "test"),
					resource.TestCheckResourceAttr("jsonfile_data.data", "nested.fixed", "fixed"),
				),
			},
		},
	})
}

func testAccExampleResourceConfig() string {
	return `
resource "jsonfile_data" "data" {
	value = "test"
}`
}
