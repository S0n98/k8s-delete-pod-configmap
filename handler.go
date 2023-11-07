package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admission "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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

	//bodyBytes = []byte(`{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"358527d5-564d-4917-b962-37a914851df8","kind":{"group":"apps","version":"v1","kind":"ReplicaSet"},"resource":{"group":"apps","version":"v1","resource":"replicasets"},"requestKind":{"group":"apps","version":"v1","kind":"ReplicaSet"},"requestResource":{"group":"apps","version":"v1","resource":"replicasets"},"name":"blue-green-fe-df769f5c4","namespace":"nginx","operation":"UPDATE","userInfo":{"username":"system:serviceaccount:argocd:argocd-server","uid":"70e895ff-69dd-4f30-bb33-6c1b989b4a0e","groups":["system:serviceaccounts","system:serviceaccounts:argocd","system:authenticated"],"extra":{"authentication.kubernetes.io/pod-name":["argocd-local-server-77cb99cb7c-d9xvc"],"authentication.kubernetes.io/pod-uid":["5a1ba978-c60b-44e9-af57-be6788011256"]}},"object":{"kind":"ReplicaSet","apiVersion":"apps/v1","metadata":{"name":"blue-green-fe-df769f5c4","namespace":"nginx","uid":"f68a38db-96bc-4c3d-83b1-b9a6021636f3","resourceVersion":"339831282","generation":6,"creationTimestamp":"2023-11-06T09:10:55Z","labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"df769f5c4"},"annotations":{"rollout.argoproj.io/desired-replicas":"2","rollout.argoproj.io/revision":"59","rollout.argoproj.io/workload-generation":"62"},"ownerReferences":[{"apiVersion":"argoproj.io/v1alpha1","kind":"Rollout","name":"blue-green-fe","uid":"c61c483d-8640-49ea-942c-0d2bfbe9dae6","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"rollouts-controller","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:12:26Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:rollout.argoproj.io/desired-replicas":{},"f:rollout.argoproj.io/revision":{},"f:rollout.argoproj.io/workload-generation":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"c61c483d-8640-49ea-942c-0d2bfbe9dae6\"}":{}}},"f:spec":{"f:replicas":{},"f:selector":{},"f:template":{"f:metadata":{"f:annotations":{".":{},"f:reloader.stakater.com/last-reloaded-from":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:name":{}},"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}}},"f:envFrom":{},"f:image":{},"f:imagePullPolicy":{},"f:lifecycle":{".":{},"f:preStop":{".":{},"f:exec":{".":{},"f:command":{}}}},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8080,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:name":{},"f:protocol":{}}},"f:resources":{".":{},"f:limits":{".":{},"f:memory":{}},"f:requests":{".":{},"f:cpu":{}}},"f:securityContext":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:serviceAccount":{},"f:serviceAccountName":{},"f:terminationGracePeriodSeconds":{}}}}}},{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:12:27Z","fieldsType":"FieldsV1","fieldsV1":{"f:status":{"f:observedGeneration":{},"f:replicas":{}}},"subresource":"status"},{"manager":"argocd-server","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T10:29:02Z","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{"f:template":{"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{"f:resources":{"f:requests":{"f:memory":{}}}}}}}}}}]},"spec":{"replicas":0,"selector":{"matchLabels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","rollouts-pod-template-hash":"df769f5c4"}},"template":{"metadata":{"name":"fe","creationTimestamp":null,"labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"df769f5c4"},"annotations":{"reloader.stakater.com/last-reloaded-from":"{\"type\":\"CONFIGMAP\",\"name\":\"dynamic-cm-43b5151c5d\",\"namespace\":\"nginx\",\"hash\":\"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385\",\"containerRefs\":[\"fe\"],\"observedAt\":1698994902}"}},"spec":{"containers":[{"name":"fe","image":"docker.io/kostiscodefresh/loan:latest","ports":[{"name":"http","containerPort":8080,"protocol":"TCP"}],"envFrom":[{"configMapRef":{"name":"dynamic-cm-5e0badd93b","optional":true}}],"env":[{"name":"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"}],"resources":{"limits":{"memory":"1Gi"},"requests":{"cpu":"100m","memory":"100Mi"}},"lifecycle":{"preStop":{"exec":{"command":["/bin/bash","-c","sleep 30"]}}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler"}}},"status":{"replicas":0,"observedGeneration":5}},"oldObject":{"kind":"ReplicaSet","apiVersion":"apps/v1","metadata":{"name":"blue-green-fe-df769f5c4","namespace":"nginx","uid":"f68a38db-96bc-4c3d-83b1-b9a6021636f3","resourceVersion":"339831282","generation":5,"creationTimestamp":"2023-11-06T09:10:55Z","labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"df769f5c4"},"annotations":{"rollout.argoproj.io/desired-replicas":"2","rollout.argoproj.io/revision":"59","rollout.argoproj.io/workload-generation":"62"},"ownerReferences":[{"apiVersion":"argoproj.io/v1alpha1","kind":"Rollout","name":"blue-green-fe","uid":"c61c483d-8640-49ea-942c-0d2bfbe9dae6","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"rollouts-controller","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:12:26Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:rollout.argoproj.io/desired-replicas":{},"f:rollout.argoproj.io/revision":{},"f:rollout.argoproj.io/workload-generation":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"c61c483d-8640-49ea-942c-0d2bfbe9dae6\"}":{}}},"f:spec":{"f:replicas":{},"f:selector":{},"f:template":{"f:metadata":{"f:annotations":{".":{},"f:reloader.stakater.com/last-reloaded-from":{}},"f:labels":{".":{},"f:app.kubernetes.io/instance":{},"f:app.kubernetes.io/name":{},"f:delete-configmap":{},"f:rollouts-pod-template-hash":{}},"f:name":{}},"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP\"}":{".":{},"f:name":{},"f:value":{}}},"f:envFrom":{},"f:image":{},"f:imagePullPolicy":{},"f:lifecycle":{".":{},"f:preStop":{".":{},"f:exec":{".":{},"f:command":{}}}},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8080,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:name":{},"f:protocol":{}}},"f:resources":{".":{},"f:limits":{".":{},"f:memory":{}},"f:requests":{".":{},"f:cpu":{}}},"f:securityContext":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:serviceAccount":{},"f:serviceAccountName":{},"f:terminationGracePeriodSeconds":{}}}}}},{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T09:12:27Z","fieldsType":"FieldsV1","fieldsV1":{"f:status":{"f:observedGeneration":{},"f:replicas":{}}},"subresource":"status"},{"manager":"argocd-server","operation":"Update","apiVersion":"apps/v1","time":"2023-11-06T10:29:02Z","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{"f:template":{"f:spec":{"f:containers":{"k:{\"name\":\"fe\"}":{"f:resources":{"f:requests":{"f:memory":{}}}}}}}}}}]},"spec":{"replicas":0,"selector":{"matchLabels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","rollouts-pod-template-hash":"df769f5c4"}},"template":{"metadata":{"name":"fe","creationTimestamp":null,"labels":{"app.kubernetes.io/instance":"blue-green","app.kubernetes.io/name":"fe","delete-configmap":"true","rollouts-pod-template-hash":"df769f5c4"},"annotations":{"reloader.stakater.com/last-reloaded-from":"{\"type\":\"CONFIGMAP\",\"name\":\"dynamic-cm-43b5151c5d\",\"namespace\":\"nginx\",\"hash\":\"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385\",\"containerRefs\":[\"fe\"],\"observedAt\":1698994902}"}},"spec":{"containers":[{"name":"fe","image":"docker.io/kostiscodefresh/loan:latest","ports":[{"name":"http","containerPort":8080,"protocol":"TCP"}],"envFrom":[{"configMapRef":{"name":"dynamic-cm-5e0badd93b","optional":true}}],"env":[{"name":"STAKATER_DYNAMIC_CM_A53455E90A_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EE376EB0F8_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_EAB382CDD0_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_D8F0D3B85B_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"},{"name":"STAKATER_DYNAMIC_CM_6ADBCFAEE7_CONFIGMAP","value":"01a7f3b76d0bdffc6c0e891b4becd3fa1f019385"}],"resources":{"limits":{"memory":"1Gi"},"requests":{"cpu":"100m","memory":"200Mi"}},"lifecycle":{"preStop":{"exec":{"command":["/bin/bash","-c","sleep 30"]}}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler"}}},"status":{"replicas":0,"observedGeneration":5}},"dryRun":false,"options":{"kind":"UpdateOptions","apiVersion":"meta.k8s.io/v1"}}}`)

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

	var configmapsList []string

	if oldReplicaSet.Status.Replicas > 0 && replicaSet.Status.Replicas == 0 {
		for _, v := range oldReplicaSet.Spec.Template.Spec.Containers {
			for _, env := range v.EnvFrom {
				configmapsList = append(configmapsList, env.ConfigMapRef.Name)
			}
		}
	}

	fmt.Println(configmapsList)

	return &admission.AdmissionResponse{Allowed: true}
}

func deleteConfigmap(name string) {

}
