package controllers_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	Recipe "github.com/ramendr/recipe/api/v1alpha1"
)

var (
	testNamespace         *corev1.Namespace
	testNamespaceNameBase string = "recipe-test-ns"
)

var _ = Describe("RecipeController", func() {
	var testCtx context.Context
	var cancel context.CancelFunc

	BeforeEach(func() {
		testCtx, cancel = context.WithCancel(context.TODO())
		Expect(k8sClient).NotTo(BeNil())
		testNamespace = createUniqueNamespace(testCtx)
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(testCtx, testNamespace)).To(Succeed())

		cancel()
	})

	Context("Groups", func() {
		It("allow unique names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{
						{
							Name: "group-1",
							Type: "resource",
						},
						{
							Name: "group-2",
							Type: "resource",
						},
					},
					Hooks:     []*Recipe.Hook{},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
		It("error on duplicate names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{
						{
							Name: "group-1",
							Type: "resource",
						},
						{
							Name: "group-1",
							Type: "resource",
						},
					},
					Hooks:     []*Recipe.Hook{},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("Hooks", func() {
		It("allow unique names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
						},
						{
							Name: "hook-2",
							Type: "exec",
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
		It("error on duplicate names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
						},
						{
							Name: "hook-1",
							Type: "exec",
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
		})
		It("allow unique Ops names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
							Ops: []*Recipe.Operation{
								{
									Name: "op-1",
								},
								{
									Name: "op-2",
								},
							},
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
		It("error on duplicate Ops names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
							Ops: []*Recipe.Operation{
								{
									Name: "op-1",
								},
								{
									Name: "op-1",
								},
							},
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
		})
		It("allow unique Check names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
							Chks: []*Recipe.Check{
								{
									Name: "check-1",
								},
								{
									Name: "check-2",
								},
							},
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
		It("error on duplicate Check names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
							Chks: []*Recipe.Check{
								{
									Name: "check-1",
								},
								{
									Name: "check-1",
								},
							},
						},
					},
					Workflows: []*Recipe.Workflow{},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("Workflows", func() {
		It("allow unique names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks:  []*Recipe.Hook{},
					Workflows: []*Recipe.Workflow{
						{
							Name:     "workflow-1",
							Sequence: []map[string]string{},
						},
						{
							Name:     "workflow-2",
							Sequence: []map[string]string{},
						},
					},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
		It("error on duplicate names", func() {
			recipe := &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks:  []*Recipe.Hook{},
					Workflows: []*Recipe.Workflow{
						{
							Name:     "workflow-1",
							Sequence: []map[string]string{},
						},
						{
							Name:     "workflow-2",
							Sequence: []map[string]string{},
						},
					},
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).To(BeNil())
		})
	})
})

// since object names are reused, use unique namespaces
func createUniqueNamespace(testCtx context.Context) *corev1.Namespace {
	testNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: testNamespaceNameBase,
		},
	}
	Expect(k8sClient.Create(testCtx, testNamespace)).To(Succeed())
	Expect(testNamespace.GetName()).NotTo(BeEmpty())

	return testNamespace
}
