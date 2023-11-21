package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecFactory  = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecFactory.UniversalDeserializer()
)

// add kind AdmissionReview in scheme
func init() {
	//_ = corev1.AddToScheme(runtimeScheme)
	_ = admission.AddToScheme(runtimeScheme)
	//_ = v1.AddToScheme(runtimeScheme)
}

func validate(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte

	// read request body
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			bodyBytes = data
		} else {
			loggerErr.Print("Reading body failed")
			return
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		loggerErr.Printf("Invalid content type %s", contentType)
		return
	}

	logger.Printf("Handling request: %s", string(bodyBytes))
	var responseObj runtime.Object

	// decode body parse to admission.AdmissionReview
	if obj, gvk, err := deserializer.Decode(bodyBytes, nil, nil); err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		loggerErr.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return

	} else {
		requestedAdmissionReview, ok := obj.(*admission.AdmissionReview)
		if !ok {
			loggerErr.Printf("Expected v1.AdmissionReview but got: %T", obj)
			return
		}

		// create response body

		responseAdmissionReview := &admission.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(*gvk)
		responseAdmissionReview.Response = extendedTask(*requestedAdmissionReview)
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		responseObj = responseAdmissionReview
	}

	responseBytes, err := json.Marshal(responseObj)

	if err != nil {
		loggerErr.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseBytes)
}

func extendedTask(ar admission.AdmissionReview) *admission.AdmissionResponse {

	rawObject := ar.Request.Object.Raw
	rawOldObject := ar.Request.OldObject.Raw

	// get old, new configmap
	newConfigmap := corev1.ConfigMap{}
	oldConfigmap := corev1.ConfigMap{}

	if _, _, err := deserializer.Decode(rawObject, nil, &newConfigmap); err != nil {
		loggerErr.Print(err)
		return &admission.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	if _, _, err := deserializer.Decode(rawOldObject, nil, &oldConfigmap); err != nil {
		loggerErr.Print(err)
		return &admission.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// only handle when configmap's annotation 'ftech.rollouts.promoted: false'

	if oldConfigmap.Annotations["ftech.rollouts.promoted"] == "true" {
		return &admission.AdmissionResponse{Allowed: true}
	}

	// get rollouts app name from annotations and clean up resources

	if oldConfigmap.Labels["ftech.rollouts.app"] != "" {
		logger.Printf("Processing argocd app: %s", oldConfigmap.Labels["ftech.rollouts.app"])
		cleanUpResource(oldConfigmap.Namespace, oldConfigmap.Labels["ftech.rollouts.app"])
	}

	return &admission.AdmissionResponse{Allowed: true}
}

func cleanUpResource(namespace string, appName string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		loggerErr.Println(err.Error())
		return
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		loggerErr.Println(err.Error())
		return
	}

	// create list option
	labelCondition := fmt.Sprintf("ftech.rollouts.app=%s", appName)

	// get list pod by rollout's app name created by argocd webhook
	podsList, podDetailErr := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelCondition})

	if podDetailErr != nil {
		loggerErr.Println("List pod info error: ", podDetailErr.Error())
		return
	}

	for _, pod := range podsList.Items {
		// delete all pod found
		deletePodErr := clientset.CoreV1().Pods(pod.Namespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})

		logger.Printf("Deleted pod %s", pod.Name)

		if deletePodErr != nil {
			if strings.Contains(deletePodErr.Error(), "not found") {
				continue
			}

			loggerErr.Printf("Error when delete pod: %s , %s", pod.Name, deletePodErr.Error())
		}
	}

	// list all configmap by rollout's app name
	configMapList, cmDetailErr := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelCondition})

	if cmDetailErr != nil {
		loggerErr.Println("List configmap info error: ", cmDetailErr.Error())
		return
	}

	// delete all configmap except 5 latest created.
	configMaps := configMapList.Items[:]
	sort.Slice(configMaps, func(i, j int) bool {
		return configMaps[i].CreationTimestamp.Before(&configMaps[j].CreationTimestamp)
	})

	logger.Print(configMaps)

	if len(configMaps) < 6 {
		return
	}

	backgroundDeletetion := metav1.DeletePropagationBackground

	for _, cm := range configMaps[:len(configMaps)-5] {
		deleteJobErr := clientset.BatchV1().Jobs(cm.Namespace).Delete(context.Background(), cm.Name, metav1.DeleteOptions{PropagationPolicy: &backgroundDeletetion})

		logger.Printf("Deleted configmap, job %s", cm.Name)

		if deleteJobErr != nil {
			if strings.Contains(deleteJobErr.Error(), "not found") {
				continue
			}

			loggerErr.Printf("Error when delete job: %s , %s", cm.Name, deleteJobErr.Error())
		}
	}
}
