package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	var opsfilegenPath string
	BeforeSuite(func() {
		var err error
		opsfilegenPath, err = gexec.Build("github.com/crawsible/opsfilegen")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("generates an opsfile from a source target manifest", func() {
		wd, _ := os.Getwd()

		sourceManifestPath := filepath.Join(wd, "fixtures/source.yml")
		targetManifestPath := filepath.Join(wd, "fixtures/target.yml")
		expectedOpsFilePath := filepath.Join(wd, "fixtures/expected_opsfile.yml")
		expectedOutput, _ := ioutil.ReadFile(expectedOpsFilePath)

		command := exec.Command(
			opsfilegenPath,
			sourceManifestPath,
			targetManifestPath,
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		actualOutput := session.Wait(5 * time.Second).Out.Contents()

		Expect(actualOutput).To(MatchYAML(expectedOutput))
	})
})
