/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lms

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lmsv1alpha1 "github.com/krestomatio/lms-moodle-operator/api/lms/v1alpha1"
)

var _ = Describe("LMSMoodle Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"
		const templateName = "test-lmsmoodle-template"
		// Namespace and child resource name derived from resourceName
		const childName = LMSMoodleNamePrefix + resourceName

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		lmsmoodle := &lmsv1alpha1.LMSMoodle{}

		// Same GVKs as cmd/main.go; testdata/crd defines their CRDs
		postgresGVK := schema.GroupVersionKind{Group: "postgres.krestomat.io", Version: "v1alpha1", Kind: "Postgres"}
		keydbGVK := schema.GroupVersionKind{Group: "keydb.krestomat.io", Version: "v1alpha1", Kind: "Keydb"}
		nfsGVK := schema.GroupVersionKind{Group: "nfs.krestomat.io", Version: "v1alpha1", Kind: "Ganesha"}
		moodleGVK := schema.GroupVersionKind{Group: "m4e.krestomat.io", Version: "v1alpha1", Kind: "Moodle"}

		BeforeEach(func() {
			By("creating the LMSMoodleTemplate used by the LMSMoodle")
			template := &lmsv1alpha1.LMSMoodleTemplate{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: templateName}, template)
			if err != nil && errors.IsNotFound(err) {
				template = &lmsv1alpha1.LMSMoodleTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Name: templateName,
					},
					Spec: lmsv1alpha1.LMSMoodleTemplateSpec{
						MoodleSpec: newTestMoodleSpec(),
					},
				}
				Expect(k8sClient.Create(ctx, template)).To(Succeed())
			}

			By("creating the custom resource for the Kind LMSMoodle")
			err = k8sClient.Get(ctx, typeNamespacedName, lmsmoodle)
			if err != nil && errors.IsNotFound(err) {
				resource := &lmsv1alpha1.LMSMoodle{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: lmsv1alpha1.LMSMoodleSpec{
						LMSMoodleTemplateName: templateName,
						LMSMoodleTemplateSpec: lmsv1alpha1.LMSMoodleTemplateSpec{
							MoodleSpec: newTestMoodleSpec(),
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &lmsv1alpha1.LMSMoodle{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance LMSMoodle")
			// No manager runs the finalizer in this suite, so remove it to let the delete finish.
			resource.SetFinalizers(nil)
			Expect(k8sClient.Update(ctx, resource)).To(Succeed())
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())

			By("Cleanup the LMSMoodleTemplate")
			template := &lmsv1alpha1.LMSMoodleTemplate{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: templateName}, template)).To(Succeed())
			Expect(k8sClient.Delete(ctx, template)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &LMSMoodleReconciler{
				Client:      k8sClient,
				Scheme:      k8sClient.Scheme(),
				MoodleGVK:   moodleGVK,
				NfsGVK:      nfsGVK,
				KeydbGVK:    keydbGVK,
				PostgresGVK: postgresGVK,
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking the namespace and the child resources applied before Moodle")
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: childName}, &corev1.Namespace{})).To(Succeed())
			for _, gvk := range []schema.GroupVersionKind{postgresGVK, keydbGVK, nfsGVK} {
				child := newUnstructuredObject(gvk)
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: childName, Namespace: childName}, child)).To(Succeed())
			}
		})
	})
})
