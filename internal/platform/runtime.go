package platform

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func EnsureRuntime(
	ctx context.Context,
	c *kubernetes.Clientset,
	tenant string,
	runtime string,
) error {

	switch runtime {
	case "wordpress":
		return ensureWordPress(ctx, c, tenant)
	default:
		return fmt.Errorf("unknown runtime: %s", runtime)
	}
}

func ensureWordPress(ctx context.Context, c *kubernetes.Clientset, ns string) error {
	const name = "wordpress"
	labels := map[string]string{"app": name}

	// --- Deployment ---
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: "wordpress:php8.2-apache",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80},
							},
							ReadinessProbe: tcpProbe(),
							Resources:      resourceLimits(),
						},
					},
				},
			},
		},
	}

	existing, err := c.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		// Update existing deployment to match desired spec
		existing.Spec = deploy.Spec
		if _, err := c.AppsV1().Deployments(ns).Update(ctx, existing, metav1.UpdateOptions{}); err != nil {
			return fmt.Errorf("update deployment: %w", err)
		}
	} else {
		if _, err := c.AppsV1().Deployments(ns).Create(ctx, deploy, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create deployment: %w", err)
		}
	}

	// --- Service ---
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	_, err = c.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	if _, err := c.CoreV1().Services(ns).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	return nil
}

// --- helpers ---

func tcpProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(80),
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       5,
	}
}

func resourceLimits() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    quantity("100m"),
			corev1.ResourceMemory: quantity("128Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    quantity("500m"),
			corev1.ResourceMemory: quantity("256Mi"),
		},
	}
}

func quantity(v string) resource.Quantity {
	q, err := resource.ParseQuantity(v)
	if err != nil {
		panic(err) // safe here because values are constants
	}
	return q
}

func int32Ptr(v int32) *int32 { return &v }
