package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admission "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type configmap struct {
	Name      string
	Namespace string
}

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

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

	//// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		loggerErr.Printf("Invalid content type %s", contentType)
		return
	}

	//bodyBytes = []byte(`{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"08db498d-788a-4d1f-9f7e-893356de483f","kind":{"group":"apps","version":"v1","kind":"ReplicaSet"},"resource":{"group":"apps","version":"v1","resource":"replicasets"},"requestKind":{"group":"apps","version":"v1","kind":"ReplicaSet"},"requestResource":{"group":"apps","version":"v1","resource":"replicasets"},"name":"blue-green-fe-599f64489f","namespace":"nginx","operation":"UPDATE","userInfo":{"username":"system:serviceaccount:argo-rollouts:argo-rollouts","uid":"ca0b9aad-489a-45de-a3e4-8ee6717f50f4","groups":["system:serviceaccounts","system:serviceaccounts:argo-rollouts","system:authenticated"],"extra":{"authentication.kubernetes.io/pod-name":["argo-rollouts-64595b7c59-5xpt4"],"authentication.kubernetes.io/pod-uid":["109a1486-cae9-400c-8ba1-7b3cc66a8b80"]}},"object":{"kind":"ReplicaSet","apiVersion":"apps/v1","metadata":{"name":"blue-green-fe-599f64489f","namespace":"nginx","selfLink":"/apis/apps/v1/namespaces/nginx/replicasets/blue-green-fe-599f64489f","uid":"450bc3a2-2cda-414b-ae33-c82ba2cb3ef0","resourceVersion":"340313545","generation":4,"creationTimestamp":"2023-11-06T09:36:32Z","labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"599f64489f"},"annotations":{"rollout.argoproj.io/desired-replicas":"2","rollout.argoproj.io/revision":"61","rollout.argoproj.io/workload-generation":"64"},"ownerReferences":[{"apiVersion":"argoproj.io/v1alpha1","kind":"Rollout","name":"blue-green-fe","uid":"c61c483d-8640-49ea-942c-0d2bfbe9dae6","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:37:13Z","fieldsType":"FieldsV1","fieldsV1":{"f:status":{"f:availableReplicas":{},"f:fullyLabeledReplicas":{},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{}}},"subresource":"status"},{"manager":"rollouts-controller","operation":"Update","apiVersion":"apps/v1","time":"2023-11-07T07:01:54Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:rollout.argoproj.io/desired-replicas":{},"f:rollout.argoproj.io/revision":{},"f:rollout.argoproj.io/workload-generation":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"c61c483d-8640-49ea-942c-0d2bfbe9dae6\"}":{}}},"f:spec":{"f:replicas":{},"f:selector":{},"f:template":{"f:metadata":{"f:annotations":{".":{},"f:reloader.stakater.com/last-reloaded-from":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:name":{}},"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"STAKATER_DYNAMIC_CM_5E0BADD93B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}}},"f:envFrom":{},"f:image":{},"f:imagePullPolicy":{},"f:lifecycle":{".":{},"f:preStop":{".":{},"f:exec":{".":{},"f:command":{}}}},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8080,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:name":{},"f:protocol":{}}},"f:resources":{".":{},"f:limits":{".":{},"f:memory":{}},"f:requests":{".":{},"f:cpu":{},"f:memory":{}}},"f:securityContext":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:serviceAccount":{},"f:serviceAccountName":{},"f:terminationGracePeriodSeconds":{}}}}}}]},"spec":{"replicas":0,"selector":{"matchLabels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","rollouts-pod-template-hash":"599f64489f"}},"template":{"metadata":{"name":"fe","creationTimestamp":null,"labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"599f64489f"},"annotations":{"reloader.stakater.com/last-reloaded-from":"{\"type\":\"CONFIGMAP\",\"name\":\"dynamic-cm-43b5151c5d\",\"namespace\":\"nginx\",\"hash\":\"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385\",\"containerRefs\":[\"fe\"],\"observedAt\":1698994902}"}},"spec":{"containers":[{"name":"fe","image":"docker.io/kostiscodefresh/loan:latest","ports":[{"name":"http","containerPort":8080,"protocol":"TCP"}],"envFrom":[{"configMapRef":{"name":"dynamic-cm-5e0badd93b","optional":true}}],"env":[{"name":"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_5E0BADD93B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"}],"resources":{"limits":{"memory":"1Gi"},"requests":{"cpu":"100m","memory":"200Mi"}},"lifecycle":{"preStop":{"exec":{"command":["/bin/bash","-c","sleep 30"]}}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler"}}},"status":{"replicas":2,"fullyLabeledReplicas":2,"readyReplicas":2,"availableReplicas":2,"observedGeneration":3}},"oldObject":{"kind":"ReplicaSet","apiVersion":"apps/v1","metadata":{"name":"blue-green-fe-599f64489f","namespace":"nginx","uid":"450bc3a2-2cda-414b-ae33-c82ba2cb3ef0","resourceVersion":"340313545","generation":3,"creationTimestamp":"2023-11-06T09:36:32Z","labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"599f64489f"},"annotations":{"rollout.argoproj.io/desired-replicas":"2","rollout.argoproj.io/revision":"61","rollout.argoproj.io/workload-generation":"64","scale-down-deadline":"2023-11-07T07:01:55Z"},"ownerReferences":[{"apiVersion":"argoproj.io/v1alpha1","kind":"Rollout","name":"blue-green-fe","uid":"c61c483d-8640-49ea-942c-0d2bfbe9dae6","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:37:13Z","fieldsType":"FieldsV1","fieldsV1":{"f:status":{"f:availableReplicas":{},"f:fullyLabeledReplicas":{},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{}}},"subresource":"status"},{"manager":"rollouts-controller","operation":"Update","apiVersion":"apps/v1","time":"2023-11-07T07:01:15Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:rollout.argoproj.io/desired-replicas":{},"f:rollout.argoproj.io/revision":{},"f:rollout.argoproj.io/workload-generation":{},"f:scale-down-deadline":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"c61c483d-8640-49ea-942c-0d2bfbe9dae6\"}":{}}},"f:spec":{"f:replicas":{},"f:selector":{},"f:template":{"f:metadata":{"f:annotations":{".":{},"f:reloader.stakater.com/last-reloaded-from":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:name":{}},"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"STAKATER_DYNAMIC_CM_5E0BADD93B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}}},"f:envFrom":{},"f:image":{},"f:imagePullPolicy":{},"f:lifecycle":{".":{},"f:preStop":{".":{},"f:exec":{".":{},"f:command":{}}}},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8080,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:name":{},"f:protocol":{}}},"f:resources":{".":{},"f:limits":{".":{},"f:memory":{}},"f:requests":{".":{},"f:cpu":{},"f:memory":{}}},"f:securityContext":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:serviceAccount":{},"f:serviceAccountName":{},"f:terminationGracePeriodSeconds":{}}}}}}]},"spec":{"replicas":2,"selector":{"matchLabels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","rollouts-pod-template-hash":"599f64489f"}},"template":{"metadata":{"name":"fe","creationTimestamp":null,"labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"599f64489f"},"annotations":{"reloader.stakater.com/last-reloaded-from":"{\"type\":\"CONFIGMAP\",\"name\":\"dynamic-cm-43b5151c5d\",\"namespace\":\"nginx\",\"hash\":\"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385\",\"containerRefs\":[\"fe\"],\"observedAt\":1698994902}"}},"spec":{"containers":[{"name":"fe","image":"docker.io/kostiscodefresh/loan:latest","ports":[{"name":"http","containerPort":8080,"protocol":"TCP"}],"envFrom":[{"configMapRef":{"name":"dynamic-cm-5e0badd93b","optional":true}}],"env":[{"name":"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_5E0BADD93B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"}],"resources":{"limits":{"memory":"1Gi"},"requests":{"cpu":"100m","memory":"200Mi"}},"lifecycle":{"preStop":{"exec":{"command":["/bin/bash","-c","sleep 30"]}}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler"}}},"status":{"replicas":2,"fullyLabeledReplicas":2,"readyReplicas":2,"availableReplicas":2,"observedGeneration":3}},"dryRun":false,"options":{"kind":"UpdateOptions","apiVersion":"meta.k8s.io/v1"}}}`)

	logger.Printf("Handling request: %s", string(bodyBytes))
	var responseObj runtime.Object

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

	replicaSet := appsv1.ReplicaSet{}
	oldReplicaSet := appsv1.ReplicaSet{}

	if _, _, err := deserializer.Decode(rawObject, nil, &replicaSet); err != nil {
		loggerErr.Print(err)
		return &admission.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	if _, _, err := deserializer.Decode(rawOldObject, nil, &oldReplicaSet); err != nil {
		loggerErr.Print(err)
		return &admission.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	if oldReplicaSet.Labels["delete-configmap"] != "true" {
		return &admission.AdmissionResponse{Allowed: true}
	}

	var configmapsList []configmap

	if *oldReplicaSet.Spec.Replicas > 0 && *replicaSet.Spec.Replicas == 0 {
		for _, v := range oldReplicaSet.Spec.Template.Spec.Containers {
			for _, env := range v.EnvFrom {
				configmapsList = append(configmapsList, configmap{Name: env.ConfigMapRef.Name, Namespace: oldReplicaSet.Namespace})
			}
		}
	}

	fmt.Println(configmapsList)

	if len(configmapsList) > 0 {
		deleteConfigmap(configmapsList)
	}

	return &admission.AdmissionResponse{Allowed: true}
}

func deleteConfigmap(configmapsList []configmap) {
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

	for _, cm := range configmapsList {
		cmDetail, cmDetailErr := clientset.CoreV1().ConfigMaps(cm.Namespace).Get(context.TODO(), cm.Name, metav1.GetOptions{})

		if cmDetailErr != nil {
			loggerErr.Println("Get configmap info error: ", cmDetailErr.Error())
			continue
		}

		if cmDetail.Annotations["delete-on-pod-termination"] != "true" {
			continue
		}

		if cmDetail.Annotations["deleted"] != "true" {
			jobDetail, jobDetailErr := clientset.BatchV1().Jobs(cm.Namespace).Get(context.TODO(), cm.Name, metav1.GetOptions{})

			if jobDetailErr != nil {
				loggerErr.Println("Get job create configmap info error: ", jobDetailErr.Error())
				continue
			}

			// annotations and ownerRef to Job
			payload := fmt.Sprintf(`{"metadata": {"annotations": {"deleted": "true"}, "ownerReferences": [{"apiVersion": "batch/v1", "blockOwnerDeletion": "true", "controller": "true", "kind": "Job", "name": "%s", "uid": "%s}]}}`, cm.Name, jobDetail.UID)
			clientset.CoreV1().ConfigMaps(cm.Namespace).Patch(context.TODO(), cm.Name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
			continue
		}

		deleteErr := clientset.CoreV1().ConfigMaps(cm.Namespace).Delete(context.Background(), cm.Name, metav1.DeleteOptions{})
		if deleteErr != nil {
			loggerErr.Printf("Error when delete configmap: %s , %s", cm.Name, deleteErr.Error())
		}

		deleteJobErr := clientset.BatchV1().Jobs(cm.Namespace).Delete(context.Background(), cm.Name, metav1.DeleteOptions{})
		if deleteJobErr != nil {
			loggerErr.Printf("Error when delete job: %s , %s", cm.Name, deleteJobErr.Error())
		}
	}
}
