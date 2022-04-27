package qovery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ClusterDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccClusterDataSourceConfig(
					getTestClusterID(),
					getTestOrganizationID(),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "id", getTestClusterID()),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "organization_id", getTestOrganizationID()),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "credentials_id", getTestAWSCredentialsID()),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "name", "test_terraform_provider"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "description", "cluster to test terraform"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "region", "eu-west-3"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "cpu", "2000"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "memory", "4096"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "min_running_nodes", "3"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "max_running_nodes", "10"),
					resource.TestCheckResourceAttr("data.qovery_cluster.test", "state", "RUNNING"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig(credentialsID string, organizationID string) string {
	return fmt.Sprintf(`
data "qovery_cluster" "test" {
  id = "%s"
  organization_id = "%s"
}
`, credentialsID, organizationID,
	)
}