package libdrmaa

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Jobtemplate", func() {

	Context("basic tests", func() {

		It("should convert a JobTemplate back and forth", func() {
			originalTemplate := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
			originalTemplate.ExtensionList = map[string]string{
				"DRMAA1_NATIVE_SPECIFICATION": "-l gpu",
			}
			s, err := drmaa.MakeSession()
			Expect(err).To(BeNil())
			defer s.Exit()
			jt, err := s.AllocateJobTemplate()
			Expect(err).To(BeNil())
			err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(originalTemplate, &jt)
			Expect(err).To(BeNil())
			convertedJobTemplate, err := ConvertDRMAAJobTemplateToDRMAA2JobTemplate(&jt)
			Expect(err).To(BeNil())
			Expect(convertedJobTemplate.RemoteCommand).To(Equal(originalTemplate.RemoteCommand))
			Expect(len(convertedJobTemplate.Args)).To(BeNumerically("==", len(originalTemplate.Args)))
			Expect(convertedJobTemplate.Args[0]).To(Equal(originalTemplate.Args[0]))
			Expect(convertedJobTemplate.ExtensionList["DRMAA1_NATIVE_SPECIFICATION"]).To(Equal(originalTemplate.ExtensionList["DRMAA1_NATIVE_SPECIFICATION"]))
		})

		It("should convert a JobTemplate back and forth", func() {
			originalTemplate := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
			s, err := drmaa.MakeSession()
			Expect(err).To(BeNil())
			defer s.Exit()
			jt, err := s.AllocateJobTemplate()
			Expect(err).To(BeNil())
			err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(originalTemplate, &jt)
			Expect(err).To(BeNil())
			convertedJobTemplate, err := ConvertDRMAAJobTemplateToDRMAA2JobTemplate(&jt)
			Expect(err).To(BeNil())
			Expect(convertedJobTemplate.RemoteCommand).To(Equal(originalTemplate.RemoteCommand))
			Expect(len(convertedJobTemplate.Args)).To(BeNumerically("==", len(originalTemplate.Args)))
			Expect(convertedJobTemplate.Args[0]).To(Equal(originalTemplate.Args[0]))
			Expect(convertedJobTemplate.ExtensionList).To(BeNil())
		})

	})

	Context("Runtime tests", func() {

		It("should set the environment variables", func() {
			jt := drmaa2interface.JobTemplate{
				RemoteCommand:  "/bin/bash",
				Args:           []string{"-c", "exit $EXIT"},
				JobEnvironment: map[string]string{"EXIT": "77"},
			}
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			jobID, err := d.AddJob(jt)
			Expect(err).To(BeNil())

			err = d.Wait(jobID, time.Second*60, drmaa2interface.Failed, drmaa2interface.Done)
			Expect(err).To(BeNil())

			ji, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())
			Expect(ji.ExitStatus).To(BeNumerically("==", 77))
		})

	})

	Context("Regressions", func() {

		It("should not fail when extension is not set", func() {
			originalTemplate := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
			s, err := drmaa.MakeSession()
			Expect(err).To(BeNil())
			jt, err := s.AllocateJobTemplate()
			s.Exit()
			Expect(err).To(BeNil())

			err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(originalTemplate, &jt)
			Expect(err).To(BeNil())
		})

		It("should prefix path with : for SGE to indicate that data comes from local host", func() {
			originalTemplate := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
				InputPath:     "someFile",
				OutputPath:    "someOutputFile",
				ErrorPath:     "someErrorFile",
			}
			s, err := drmaa.MakeSession()
			Expect(err).To(BeNil())
			defer s.Exit()

			jt, err := s.AllocateJobTemplate()
			Expect(err).To(BeNil())

			err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(originalTemplate, &jt)
			Expect(err).To(BeNil())
			inputPath, err := jt.InputPath()
			Expect(err).To(BeNil())
			Expect(inputPath).To(Equal(":someFile"))
			outputPath, err := jt.OutputPath()
			Expect(err).To(BeNil())
			Expect(outputPath).To(Equal(":someOutputFile"))
			errorPath, err := jt.ErrorPath()
			Expect(err).To(BeNil())
			Expect(errorPath).To(Equal(":someErrorFile"))
		})

	})

})
