package v3

import (
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/clients"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

func TestPolicyWorkflow(t *testing.T) {
	t.Parallel()

	client, err := clients.NewElbV3Client(t)
	th.AssertNoErr(t, err)

	lbID := createLoadBalancer(t, client)
	defer deleteLoadbalancer(t, client, lbID)

	listenerID := createListener(t, client, lbID)
	defer deleteListener(t, client, listenerID)

	poolID := createPool(t, client, lbID)
	defer deletePool(t, client, poolID)

	createOpts := policies.CreateOpts{
		Action:         policies.ActionRedirectToPool,
		ListenerID:     listenerID,
		RedirectPoolID: poolID,
		Description:    "Go SDK test policy",
		Name:           tools.RandomString("sdk-pol-", 5),
		Position:       37,
	}
	created, err := policies.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Created L7 Policy")
	id := created.ID

	defer func() {
		th.AssertNoErr(t, policies.Delete(client, id).ExtractErr())
		t.Log("Deleted L7 Policy")
	}()

	got, err := policies.Get(client, id).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, created, got)

	listOpts := policies.ListOpts{
		ListenerID: []string{listenerID},
	}
	pages, err := policies.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)
	policySlice, err := policies.ExtractPolicies(pages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(policySlice))
	th.AssertEquals(t, id, policySlice[0].ID)

	nameUpdated := tools.RandomString("updated-", 5)
	updateOpts := policies.UpdateOpts{
		Name: &nameUpdated,
	}
	updated, err := policies.Update(client, id, updateOpts).Extract()
	th.AssertNoErr(t, err)
	t.Log("Updated l7 Policy")
	th.AssertEquals(t, created.Action, updated.Action)
	th.AssertEquals(t, id, updated.ID)

	got2, _ := policies.Get(client, id).Extract()
	th.AssertDeepEquals(t, updated, got2)
}
