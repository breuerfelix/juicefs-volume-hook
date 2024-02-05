package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const VolumeAnnotation = "juicefs.volume.hook/mount-propagation"

type podWebhook struct {
	Client         client.Client
	decoder        *admission.Decoder
	Annotation     bool
	StorageClasses []string
}

func (a *podWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &core.Pod{}
	if err := a.decoder.Decode(req, pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// check for the existence of a pod annotation if enabled
	if a.Annotation {
		value, ok := pod.Annotations[VolumeAnnotation]
		if !ok {
			return admission.Allowed("Got no pod annotation.")
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if !parsed {
			return admission.Allowed("Pod annotation says false.")
		}
	}

	for i := range pod.Spec.Volumes {
		volume := &pod.Spec.Volumes[i]
		if volume.PersistentVolumeClaim == nil {
			// only update volumes that are of type pvc
			continue
		}

		// check for a given storageClassName if enabled
		if len(a.StorageClasses) > 0 {
			pvc := core.PersistentVolumeClaim{}
			if err := a.Client.Get(ctx, types.NamespacedName{
				Name:      volume.PersistentVolumeClaim.ClaimName,
				Namespace: pod.Namespace,
			}, &pvc); err != nil {
				return admission.Errored(http.StatusBadRequest, err)
			}

			if pvc.Spec.StorageClassName == nil {
				// skip this mount
				continue
			}

			found := false
			for _, name := range a.StorageClasses {
				if *pvc.Spec.StorageClassName == name {
					found = true
					break
				}
			}

			if !found {
				// pvc does not have the given StorageClassName, skip it
				continue
			}
		}

		for ii := range pod.Spec.Containers {
			container := &pod.Spec.Containers[ii]
			for iii := range container.VolumeMounts {
				volumeMount := &container.VolumeMounts[iii]

				if volumeMount.Name != volume.Name {
					continue
				}

				volumeMount.MountPropagation = (*core.MountPropagationMode)(ptr.To[string]("HostToContainer"))
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
