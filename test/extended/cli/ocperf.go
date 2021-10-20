package cli

import (
	"fmt"
	"log"
	"time"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	exutil "github.com/openshift/openshift-tests/test/extended/util"
)

var _ = g.Describe("[cli] oc cli perf", func() {
	defer g.GinkgoRecover()

	oc := exutil.NewCLIWithoutNamespace("default")

	g.It("Create 200 projects and time various oc commands durations", func() {
		deploymentConfigFixture := exutil.FixturePath("testdata", "cli", "oc-perf.yaml")

		start := time.Now()
		g.By("Try to create project and DC")
		for i := 0; i < 25; i++ {
			namespace := fmt.Sprintf("e2e-oc-cli-perf%d", i)
			err := oc.Run("new-project").Args(namespace).Execute()
			defer oc.Run("delete").Args("project", namespace).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
			err = oc.Run("create").Args("-n", namespace, "-f", deploymentConfigFixture).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
			defer oc.Run("delete").Args("-n", namespace, "-f", deploymentConfigFixture).Execute()
		}
		duration := time.Since(start)
		log.Printf("Duration for creating 200 projects and 1 deploymentConfig in each of those is %.2f seconds", duration.Seconds())

		start = time.Now()
		g.By("Try to get dcs, sa, and secrets")
		for i := 0; i < 25; i++ {
			namespace := fmt.Sprintf("e2e-oc-cli-perf%d", i)
			err := oc.Run("get").Args("dc", "-n", namespace).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
			err = oc.Run("get").Args("sa", "-n", namespace).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
			err = oc.Run("get").Args("secrets", "-n", namespace).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
		}
		duration = time.Since(start)
		log.Printf("Duration for gettings dc, sa, secrets in each of those is %.2f seconds", duration.Seconds())

		start = time.Now()
		g.By("Try to scale the dc replicas to 0")
		for i := 0; i < 25; i++ {
			namespace := fmt.Sprintf("e2e-oc-cli-perf%d", i)
			err := oc.Run("scale").Args("dc", "-n", namespace, "--replicas=0", "--all").Execute()
			o.Expect(err).NotTo(o.HaveOccurred())
		}
		duration = time.Since(start)
		log.Printf("Duration for scale the dc replicas to 0 in each of those is %.2f seconds", duration.Seconds())
	})
})
