/*
Copyright 2023.

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

//
// Generated by:
//
// operator-sdk create webhook --group redis --version v1beta1 --kind Redis --plural redises --programmatic-validation --defaulting
//

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// RedisDefaults -
type RedisDefaults struct {
	ContainerImageURL string
}

var redisDefaults RedisDefaults

// log is for logging in this package.
var redislog = logf.Log.WithName("redis-resource")

// SetupRedisDefaults - initialize Redis spec defaults for use with either internal or external webhooks
func SetupRedisDefaults(defaults RedisDefaults) {
	redisDefaults = defaults
	redislog.Info("Redis defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Redis) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-redis-openstack-org-v1beta1-redis,mutating=true,failurePolicy=fail,sideEffects=None,groups=redis.openstack.org,resources=redises,verbs=create;update,versions=v1beta1,name=mredis.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Redis{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Redis) Default() {
	redislog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this Redis spec
func (spec *RedisSpec) Default() {
	if spec.ContainerImage == "" {
		spec.ContainerImage = redisDefaults.ContainerImageURL
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-redis-openstack-org-v1beta1-redis,mutating=false,failurePolicy=fail,sideEffects=None,groups=redis.openstack.org,resources=redises,verbs=create;update,versions=v1beta1,name=vredis.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Redis{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateCreate() (admission.Warnings, error) {
	redislog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	redislog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Redis) ValidateDelete() (admission.Warnings, error) {
	redislog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
