/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	emojivoto "github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
	appv1alpha1 "github.com/easyteam-fr/emojis/operator/api/v1alpha1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
)

// EmojiReconciler reconciles a Emoji object
type EmojiReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.natives.easyteam.fr,resources=emojis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.natives.easyteam.fr,resources=emojis/status,verbs=get;update;patch

func (r *EmojiReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("emoji", req.NamespacedName)

	// Get the Emoji from the cache, check if its supported and update its status
	var emoji appv1alpha1.Emoji
	if err := r.Get(ctx, req.NamespacedName, &emoji); err != nil {
		if ignoreNotFound(err) != nil {
			log.Error(err, "unable to fetch Emoji")
		}
		log.Info("Emoji not synchronized yet")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	// name the emoji finalizer
	emojiFinalizerName := "emoji.finalizers.app.natives.easyteam.fr"

	// examine DeletionTimestamp to determine if object is under deletion
	if !emoji.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if containsString(emoji.ObjectMeta.Finalizers, emojiFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if emoji.Status.Supported != nil && *emoji.Status.Supported == true {
				if err := emojiDelete(&emoji); err != nil {
					return ctrl.Result{}, err
				}
			}
			// Notify the client when the deletion is done by removing
			// the Finalizer for the resource
			emoji.ObjectMeta.Finalizers = removeString(
				emoji.ObjectMeta.Finalizers,
				emojiFinalizerName,
			)
			if err := r.Update(context.Background(), &emoji); err != nil {
				return ctrl.Result{}, err
			}
			log.Info("finalizer removed")
			return ctrl.Result{}, nil
		}
	}
	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !containsString(emoji.ObjectMeta.Finalizers, emojiFinalizerName) {
		emoji.ObjectMeta.Finalizers = append(
			emoji.ObjectMeta.Finalizers,
			emojiFinalizerName,
		)
		if err := r.Update(context.Background(), &emoji); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("finalizer registered")
		return ctrl.Result{}, nil
	}

	if emoji.Status.Supported == nil {
		supported := false
		e := emojivoto.NewAllEmoji().WithShortcode(fmt.Sprintf(":%s:", emoji.Name))
		if e != nil {
			supported = true
		}
		emoji.Status.Supported = &supported
		if err := r.Status().Update(ctx, &emoji); err != nil {
			log.Error(err, "unable to update emoji status")
			return ctrl.Result{}, err
		}
		log.Info("Emoji status.supported updated")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *EmojiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Emoji{}).
		Complete(r)
}

func ignoreNotFound(err error) error {
	if apierror.IsNotFound(err) {
		return nil
	}
	return err
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// Simulate the removal of an emoji
func emojiDelete(emoji *appv1alpha1.Emoji) error {
	return nil
}
