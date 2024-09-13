package controllers_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	validationErrors "k8s.io/kube-openapi/pkg/validation/errors"

	Recipe "github.com/ramendr/recipe/api/v1alpha1"
)

var (
	testNamespace         *corev1.Namespace
	testNamespaceNameBase = "recipe-test-ns"
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
					Hooks: []*Recipe.Hook{},
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
					Hooks: []*Recipe.Hook{},
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
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
		})

		echoCommandTemplate := "/bin/sh/ -c echo"
		emptyCommandTemplate := ""

		commandGenerator := func(commandTemplate, suffix string) string {
			return commandTemplate + suffix
		}

		appendOp := func(ops []*Recipe.Operation, opName, commandTemplate, suffix string) []*Recipe.Operation {
			op := &Recipe.Operation{Name: opName, Command: commandGenerator(commandTemplate, suffix)}
			return append(ops, op)
		}

		hookOpsRecipe := func(ops []*Recipe.Operation) *Recipe.Recipe {
			return &Recipe.Recipe{
				TypeMeta:   metav1.TypeMeta{Kind: "Recipe", APIVersion: "ramendr.openshift.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-recipe", Namespace: testNamespace.Name},
				Spec: Recipe.RecipeSpec{
					Groups: []*Recipe.Group{},
					Hooks: []*Recipe.Hook{
						{
							Name: "hook-1",
							Type: "exec",
							Ops:  ops,
						},
					},
				},
			}
		}
		It("error on empty command", func() {
			ops := []*Recipe.Operation{}

			ops = appendOp(ops, "op-1", emptyCommandTemplate, "")

			recipe := hookOpsRecipe(ops)
			Expect(k8sClient.Create(context.TODO(), recipe)).To(MatchError(func() *errors.StatusError {
				path := field.NewPath("spec", "hooks[0]", "ops[0]", "command")
				value := ""

				return errors.NewInvalid(
					schema.GroupKind{Group: Recipe.GroupVersion.Group, Kind: "Recipe"},
					recipe.Name,
					field.ErrorList{
						field.Invalid(
							path, value, validationErrors.TooShort(
								path.String(),
								"body",
								1,
								value,
							).Error(),
						),
					},
				)
			}()))
		})
		It("allow unique Ops names", func() {
			ops := []*Recipe.Operation{}

			ops = appendOp(ops, "op-1", echoCommandTemplate, "aaa")
			ops = appendOp(ops, "op-2", echoCommandTemplate, "bbb")

			err := k8sClient.Create(context.TODO(), hookOpsRecipe(ops))

			Expect(err).To(BeNil())
		})
		It("error on duplicate Ops names", func() {
			ops := []*Recipe.Operation{}

			ops = appendOp(ops, "op-1", echoCommandTemplate, "aaa")
			ops = appendOp(ops, "op-1", echoCommandTemplate, "bbb")

			err := k8sClient.Create(context.TODO(), hookOpsRecipe(ops))

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
				},
			}

			err := k8sClient.Create(context.TODO(), recipe)

			Expect(err).ToNot(BeNil())
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
