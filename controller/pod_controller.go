package controller

import (
	"fmt"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type PodInfo struct {
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}

type PodController struct {
	clientset   *kubernetes.Clientset
	podInformer cache.SharedIndexInformer
	podStore    map[string]PodInfo
	lock        sync.RWMutex
}

func NewPodController() *PodController {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	tweek := func(option *metav1.ListOptions) {
		option.Kind = "pod" // Example label selector
	}
	factory := informers.NewFilteredSharedInformerFactory(clientset, 30*time.Second, metav1.NamespaceAll, tweek)
	podinformer := factory.Core().V1().Pods().Informer()

	podController := &PodController{
		clientset:   clientset,
		podInformer: podinformer,
		podStore:    make(map[string]PodInfo),
	}
	podinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    podController.handllerAdd,
		UpdateFunc: podController.handllerUpdate,
		DeleteFunc: podController.handllerDelete,
	})
	return podController

}

func (pc *PodController) handllerAdd(obj interface{}) {
	pc.updatePod(obj)

}
func (pc *PodController) handllerUpdate(oldobj, newObject interface{}) {
	pc.updatePod(newObject)

}

func (pc *PodController) handllerDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
	pc.lock.Lock()
	delete(pc.podStore, key)
	pc.lock.Unlock()
	fmt.Printf("Pod deleted: %s/%s\n", pod.Namespace, pod.Name)
}

func (pc *PodController) updatePod(obj interface{}) {
	pod := obj.(*v1.Pod)
	key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
	info := PodInfo{
		PodName:   pod.Name,
		Namespace: pod.Namespace,
		Status:    string(pod.Status.Phase),
	}
	pc.lock.Lock()
	pc.podStore[key] = info
	pc.lock.Unlock()
	fmt.Printf("Pod updated: %s/%s, Status: %s\n", pod.Namespace)

}

func (pc *PodController) Run() {
	stopChan := make(chan struct{})
	go pc.podInformer.Run(stopChan)
	cache.WaitForCacheSync(stopChan, pc.podInformer.HasSynced)
	<-stopChan
}
func (pc *PodController) GetPods() []PodInfo {
	pc.lock.RLock()
	result := []PodInfo{}
	for _, pod := range pc.podStore {
		result = append(result, pod)
	}
	pc.lock.RUnlock()
	return result
}

func (pc *PodController) GetSummary() map[string]int {
	pc.lock.RLock()
	summary := make(map[string]int)
	for _, pod := range pc.podStore {
		summary[pod.Status]++
	}
	pc.lock.RUnlock()
	return summary
}
