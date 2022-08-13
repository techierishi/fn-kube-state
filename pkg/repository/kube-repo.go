package repository

import (
	"context"
	"log"

	"fn-kube-state/pkg/models"
	"fn-kube-state/pkg/util"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type KubeQuery interface {
	GetPods(ctx context.Context) (models.Pods, error)
	GetDeploymentByGroup(ctx context.Context, namespace, appGroup string) (models.Deployments, error)
}

type kubeQuery struct {
	client kubernetes.Interface
}

func NewKubeQuery() KubeQuery {
	return &kubeQuery{}
}

func (m *kubeQuery) GetPods(ctx context.Context) (models.Pods, error) {
	pods, err := m.client.CoreV1().Pods("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	util.PrintJSON("pods.Items", pods.Items)

	podList := make(models.Pods, 0)
	for _, pod := range pods.Items {

		podStatus := string(pod.Status.Phase)
		podObj := &models.Pod{
			Name:   pod.ObjectMeta.Name,
			Labels: pod.ObjectMeta.Labels,
			Status: podStatus,
		}

		podList = append(podList, podObj)
	}

	return podList, nil
}

func (m *kubeQuery) GetDeploymentByGroup(ctx context.Context, namespace, appGroup string) (models.Deployments, error) {

	listOptions := metaV1.ListOptions{
		Limit: 100,
	}

	if len(appGroup) > 0 {
		labelSelector := metaV1.LabelSelector{MatchLabels: map[string]string{"applicationGroup": appGroup}}
		listOptions = metaV1.ListOptions{
			LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
			Limit:         100,
		}
	}

	deps, err := m.client.AppsV1().Deployments(namespace).List(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	util.PrintJSON("deps.Items", deps.Items)

	depList := make(models.Deployments, 0)
	for _, dep := range deps.Items {

		depObj := &models.Deployment{
			Name:             dep.Name,
			ApplicationGroup: dep.Labels["applicationGroup"],
			RunningPodsCount: dep.Status.AvailableReplicas,
		}

		depList = append(depList, depObj)

	}

	return depList, nil
}
