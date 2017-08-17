package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/bosh-load-tests/action"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
)

var _ = Describe("GetInstances", func() {
	var (
		cliRunner      clirunnerfakes.FakeRunner
		directorInfo   DirectorInfo
		deploymentName string
		instancesInfo  action.InstancesInfo
	)

	BeforeEach(func() {
		cliRunner = clirunnerfakes.FakeRunner{}
		directorInfo = action.DirectorInfo{
			UUID: "director-uuid",
			URL:  "https://example.com",
			Name: "My Little Director",
		}
		deploymentName = "my-deployment"

		instancesInfo = NewInstances(directorInfo, deploymentName, cliRunner)
	})

	It("runs instances command for given deployment and parses out IDs of instances", func() {
		// cliRunner.RunWithOutputReturns("")

		foundInstances := instancesInfo.GetInstances()
		Expect(cliRunner.RunWithOutputCallCount()).To(Equal(1))
		Expect(cliRunner.RunWithOutputArgsForCall(0)).To(Equal([]string{
			"-d",
			"my-deployment",
			"instances",
			"--json",
		}))

		Expect
	})
})
