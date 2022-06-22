package main

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podWebhook struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *podWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	if err := a.decoder.Decode(req, pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for i := range pod.Spec.Volumes {
		volume := &pod.Spec.Volumes[i]
		if volume.PersistentVolumeClaim == nil {
			// only update volumes that are of type pvc
			continue
		}

		// TODO (optional)
		// query given PVC and check for a given storageClassName

		for ii := range pod.Spec.Containers {
			container := &pod.Spec.Containers[ii]
			for iii := range container.VolumeMounts {
				volumeMount := &container.VolumeMounts[iii]

				if volumeMount.Name != volume.Name {
					continue
				}

				volumeMount.MountPropagation = (*corev1.MountPropagationMode)(pointer.String("HostToContainer"))
			}
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (a *podWebhook) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
