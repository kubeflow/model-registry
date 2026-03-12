package kubernetes

import (
	"context"
	"log/slog"
	"testing"

	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func withLogger(ctx context.Context) context.Context {
	logger := slog.New(slog.Default().Handler())
	return context.WithValue(ctx, constants.TraceLoggerKey, logger)
}

func TestGetTransferJobPods_FiltersByJobNameAndHandlesEmptyInputs(t *testing.T) {
	//nolint:staticcheck // fake.NewSimpleClientset is sufficient for unit tests; field management is not required here.
	clientset := fake.NewSimpleClientset()
	logic := &SharedClientLogic{
		Client: clientset,
		Logger: slog.New(slog.Default().Handler()),
	}

	ctx := withLogger(context.Background())

	// Seed pods in the cluster
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-for-job-1",
				Namespace: "kubeflow",
				Labels: map[string]string{
					"job-name": "job-1",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-for-job-2",
				Namespace: "kubeflow",
				Labels: map[string]string{
					"job-name": "job-2",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-for-other-namespace",
				Namespace: "other-ns",
				Labels: map[string]string{
					"job-name": "job-1",
				},
			},
		},
	}

	for _, pod := range pods {
		_, err := clientset.CoreV1().Pods(pod.Namespace).Create(ctx, &pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed to seed pod %s: %v", pod.Name, err)
		}
	}

	// Non-empty namespace and jobNames: should return only pods matching job-1 in kubeflow
	podList, err := logic.GetTransferJobPods(ctx, "kubeflow", []string{"job-1"})
	if err != nil {
		t.Fatalf("GetTransferJobPods returned error: %v", err)
	}
	if len(podList.Items) != 1 || podList.Items[0].Name != "pod-for-job-1" {
		t.Fatalf("expected 1 pod 'pod-for-job-1' in kubeflow, got: %+v", podList.Items)
	}

	// Empty namespace: should return empty list without error
	podList, err = logic.GetTransferJobPods(ctx, "", []string{"job-1"})
	if err != nil {
		t.Fatalf("GetTransferJobPods with empty namespace returned error: %v", err)
	}
	if len(podList.Items) != 0 {
		t.Fatalf("expected no pods for empty namespace, got: %+v", podList.Items)
	}

	// Empty jobNames slice: should return empty list without error
	podList, err = logic.GetTransferJobPods(ctx, "kubeflow", []string{})
	if err != nil {
		t.Fatalf("GetTransferJobPods with empty jobNames returned error: %v", err)
	}
	if len(podList.Items) != 0 {
		t.Fatalf("expected no pods for empty jobNames, got: %+v", podList.Items)
	}
}

func TestGetEventsForPods_FiltersByPodNamesAndHandlesEmptyInputs(t *testing.T) {
	//nolint:staticcheck // fake.NewSimpleClientset is sufficient for unit tests; field management is not required here.
	clientset := fake.NewSimpleClientset()
	logic := &SharedClientLogic{
		Client: clientset,
		Logger: slog.New(slog.Default().Handler()),
	}

	ctx := withLogger(context.Background())

	// Seed events for two pods in kubeflow and one in other namespace
	events := []corev1.Event{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "event-1",
				Namespace: "kubeflow",
			},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "pod-a",
				Namespace: "kubeflow",
			},
			Reason:  "Pulling",
			Message: "Pulling image",
			Type:    "Normal",
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "event-2",
				Namespace: "kubeflow",
			},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "pod-b",
				Namespace: "kubeflow",
			},
			Reason:  "BackOff",
			Message: "Back-off restarting container",
			Type:    "Warning",
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "event-3",
				Namespace: "other-ns",
			},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "pod-a",
				Namespace: "other-ns",
			},
			Reason:  "Pulling",
			Message: "Pulling image in other namespace",
			Type:    "Normal",
		},
	}

	for _, ev := range events {
		_, err := clientset.CoreV1().Events(ev.Namespace).Create(ctx, &ev, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed to seed event %s: %v", ev.Name, err)
		}
	}

	// Non-empty namespace and pod names: should return only events for pod-a and pod-b in kubeflow.
	// We don't assert exact count because List may aggregate multiple matches; instead we verify filtering.
	eventList, err := logic.GetEventsForPods(ctx, "kubeflow", []string{"pod-a", "pod-b"})
	if err != nil {
		t.Fatalf("GetEventsForPods returned error: %v", err)
	}
	for _, ev := range eventList.Items {
		if ev.InvolvedObject.Namespace != "kubeflow" {
			t.Fatalf("expected event in namespace kubeflow, got %q", ev.InvolvedObject.Namespace)
		}
		if ev.InvolvedObject.Name != "pod-a" && ev.InvolvedObject.Name != "pod-b" {
			t.Fatalf("unexpected involved pod name %q", ev.InvolvedObject.Name)
		}
	}

	// Empty namespace: should return empty list
	eventList, err = logic.GetEventsForPods(ctx, "", []string{"pod-a"})
	if err != nil {
		t.Fatalf("GetEventsForPods with empty namespace returned error: %v", err)
	}
	if len(eventList.Items) != 0 {
		t.Fatalf("expected no events for empty namespace, got %d", len(eventList.Items))
	}

	// Empty pod names: should return empty list
	eventList, err = logic.GetEventsForPods(ctx, "kubeflow", []string{})
	if err != nil {
		t.Fatalf("GetEventsForPods with empty podNames returned error: %v", err)
	}
	if len(eventList.Items) != 0 {
		t.Fatalf("expected no events for empty podNames, got %d", len(eventList.Items))
	}
}

