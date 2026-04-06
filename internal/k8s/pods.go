package k8s

import (
	"context"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo contains basic pod information
type PodInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	Ready     string            `json:"ready"`
	Age       string            `json:"age"`
	Labels    map[string]string `json:"labels"`
}

// ListPods returns all pods in the namespace with optional label selector
func (c *Client) ListPods(ctx context.Context, labelSelector string) ([]PodInfo, error) {
	opts := metav1.ListOptions{}
	if labelSelector != "" {
		opts.LabelSelector = labelSelector
	}

	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var result []PodInfo
	for _, pod := range pods.Items {
		result = append(result, PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Ready:     getReadyContainers(pod),
			Age:       age(pod.CreationTimestamp.Time),
			Labels:    pod.Labels,
		})
	}

	return result, nil
}

// GetPodLogs returns logs for a specific pod
func (c *Client) GetPodLogs(ctx context.Context, podName string, tailLines int64) (string, error) {
	opts := &corev1.PodLogOptions{}
	if tailLines > 0 {
		opts.TailLines = &tailLines
	}

	req := c.Clientset.CoreV1().Pods(c.Namespace).GetLogs(podName, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get pod logs: %w", err)
	}
	defer stream.Close()

	logs, err := io.ReadAll(stream)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// KillPod deletes a pod by name
func (c *Client) KillPod(ctx context.Context, podName string) error {
	err := c.Clientset.CoreV1().Pods(c.Namespace).Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}
	return nil
}

// getReadyContainers returns ready/total container count
func getReadyContainers(pod corev1.Pod) string {
	ready := 0
	total := len(pod.Status.ContainerStatuses)
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, total)
}

// age returns human-readable duration
func age(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
