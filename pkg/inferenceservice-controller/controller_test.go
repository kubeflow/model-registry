package inferenceservicecontroller_test

import (
	"fmt"
	"time"

	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("InferenceService Controller", func() {
	When("Creating a new InferenceService with Model Registry labels", func() {
		It("If a label with inference service id is missing, it should add it after creating the required resources on model registry", func() {
			const CorrectInferenceServicePath = "./testdata/inferenceservices/inference-service-correct.yaml"

			inferenceService := &kservev1beta1.InferenceService{}
			Expect(ConvertFileToStructuredResource(CorrectInferenceServicePath, inferenceService)).To(Succeed())

			if err := cli.Create(ctx, inferenceService); err != nil && !errors.IsAlreadyExists(err) {
				Fail(err.Error())
			}

			Eventually(func() error {
				isvc := &kservev1beta1.InferenceService{}
				err := cli.Get(ctx, types.NamespacedName{
					Name:      inferenceService.Name,
					Namespace: inferenceService.Namespace,
				}, isvc)
				if err != nil {
					return err
				}

				if isvc.Labels[inferenceServiceIDLabel] != "1" {
					return fmt.Errorf("Label for InferenceServiceID is not set, got %s", isvc.Labels[inferenceServiceIDLabel])
				}

				return nil
			}, 10*time.Second, 1*time.Second).Should(Succeed())
		})
	})
})
