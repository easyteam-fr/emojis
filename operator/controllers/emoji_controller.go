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
