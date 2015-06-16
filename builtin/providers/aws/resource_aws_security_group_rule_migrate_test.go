package aws

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAWSSecurityGroupRuleMigrateState(t *testing.T) {
	//   "id":"sg-4235098228", "from_port":"0", "source_security_group_id":"sg-11877275"}

	// 2015/06/16 16:04:21 terraform-provider-aws: 2015/06/16 16:04:21 [DEBUG] Attributes after migration:

	// map[string]string{"from_port":"0", "source_security_group_id":"sg-11877275", "id":"sg-3766347571", "security_group_id":"sg-13877277", "cidr_blocks.#":"0", "type":"ingress", "protocol":"-1", "self":"false", "to_port":"0"}, new id: sg-3766347571
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     string
		Meta         interface{}
	}{
		"v0_1": {
			StateVersion: 0,
			ID:           "sg-4235098228",
			Attributes: map[string]string{
				"self":                     "false",
				"to_port":                  "0",
				"security_group_id":        "sg-13877277",
				"cidr_blocks.#":            "0",
				"type":                     "ingress",
				"protocol":                 "-1",
				"from_port":                "0",
				"source_security_group_id": "sg-11877275",
			},
			Expected: "sg-3766347571",
		},
		"v0_2": {
			StateVersion: 0,
			ID:           "sg-1021609891",
			Attributes: map[string]string{
				"security_group_id": "sg-0981746d",
				"from_port":         "0",
				"to_port":           "0",
				"type":              "ingress",
				"self":              "false",
				"protocol":          "-1",
				"cidr_blocks.0":     "172.16.1.0/24",
				"cidr_blocks.1":     "172.16.2.0/24",
				"cidr_blocks.2":     "172.16.3.0/24",
				"cidr_blocks.3":     "172.16.4.0/24",
				"cidr_blocks.#":     "4"},
			Expected: "sg-4100229787",
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceAwsSecurityGroupRuleMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if is.ID != tc.Expected {
			t.Fatalf("bad sg rule id: %s\n\n expected: %s", is.ID, tc.Expected)
		}
	}
}
