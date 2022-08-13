package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"fn-kube-state/pkg/models"
	"fn-kube-state/pkg/util"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type KubeQuery interface {
	GetPods(ctx context.Context) (models.Pods, error)
	GetDeploymentByGroup(ctx context.Context, namespace, appGroup string) (models.Deployments, error)
	Watch(ctx context.Context, client *models.Client)
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

func (m *kubeQuery) Watch(ctx context.Context, client *models.Client) {

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(m.client, time.Second*30)
	svcInformer := kubeInformerFactory.Core().V1().Pods().Informer()

	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

			pod := obj.(*v1.Pod)
			message := fmt.Sprintf("Service added: %s \n", pod.Name)
			db := &models.SseMessage{
				Message: message,
			}
			client.Events <- db
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			message := fmt.Sprintf("Service deleted: %s \n", pod.Name)
			db := &models.SseMessage{
				Message: message,
			}
			client.Events <- db
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)
			message := fmt.Sprintf("Service changes: %s -> %s\n", oldPod.Name, newPod.Name)
			db := &models.SseMessage{
				Message: message,
			}
			client.Events <- db
		},
	})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
	for {
		time.Sleep(5 * time.Second)
	}
}
