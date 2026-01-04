package platform

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func EnsureTenant(ctx context.Context, c *kubernetes.Clientset, name string) error {
	_, err := c.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		return nil // already exists → success
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"siteforge/tenant": "true",
			},
		},
	}

	_, err = c.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create tenant %s: %w", name, err)
	}

	return nil
}
