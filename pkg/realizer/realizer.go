// Copyright 2021 VMware
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package realizer

//go:generate go run -modfile ../../hack/tools/go.mod github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/cartographer/pkg/conditions"
	"github.com/vmware-tanzu/cartographer/pkg/selector"
	"github.com/vmware-tanzu/cartographer/pkg/utils"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"github.com/vmware-tanzu/cartographer/pkg/logger"
	"github.com/vmware-tanzu/cartographer/pkg/realizer/statuses"
	"github.com/vmware-tanzu/cartographer/pkg/templates"
)

func MakeSupplychainOwnerResources(supplyChain *v1alpha1.ClusterSupplyChain) []OwnerResource {
	var resources []OwnerResource
	for _, resource := range supplyChain.Spec.Resources {
		resources = append(resources, OwnerResource{
			Name: resource.Name,
			TemplateRef: v1alpha1.TemplateReference{
				Kind: resource.TemplateRef.Kind,
				Name: resource.TemplateRef.Name,
			},
			TemplateOptions: resource.TemplateRef.Options,
			Params:          resource.Params,
			Sources:         resource.Sources,
			Images:          resource.Images,
			Configs:         resource.Configs,
		})
	}
	return resources
}

func MakeDeliveryOwnerResources(delivery *v1alpha1.ClusterDelivery) []OwnerResource {
	var resources []OwnerResource
	for _, resource := range delivery.Spec.Resources {
		resources = append(resources, OwnerResource{
			Name: resource.Name,
			TemplateRef: v1alpha1.TemplateReference{
				Kind: resource.TemplateRef.Kind,
				Name: resource.TemplateRef.Name,
			},
			TemplateOptions: resource.TemplateRef.Options,
			Params:          resource.Params,
			Sources:         resource.Sources,
			Configs:         resource.Configs,
			Deployment:      resource.Deployment,
		})
	}
	return resources
}

//counterfeiter:generate . Realizer
type Realizer interface {
	Realize(ctx context.Context, resourceRealizer ResourceRealizer, blueprintName string, ownerResources []OwnerResource, resourceStatuses statuses.ResourceStatuses) error
}

type ResourceLabeler func(resource OwnerResource) templates.Labels

type realizer struct{}

func NewRealizer() Realizer {
	return &realizer{}
}

func (r *realizer) Realize(ctx context.Context, resourceRealizer ResourceRealizer, blueprintName string, ownerResources []OwnerResource, resourceStatuses statuses.ResourceStatuses) error {
	log := logr.FromContextOrDiscard(ctx)
	log.V(logger.DEBUG).Info("Realize")

	outs := NewOutputs()
	var firstError error

	for _, resource := range ownerResources {
		log = log.WithValues("resource", resource.Name)
		ctx = logr.NewContext(ctx, log)
		template, stampedObject, out, err := resourceRealizer.Do(ctx, resource, blueprintName, outs)

		if stampedObject != nil {
			log.V(logger.DEBUG).Info("realized resource as object",
				"object", stampedObject)
		}

		if err != nil {
			log.Error(err, "failed to realize resource")

			if firstError == nil {
				firstError = err
			}
		}

		outs.AddOutput(resource.Name, out)

		previousRealizedResource := resourceStatuses.GetPreviousRealizedResource(resource.Name)

		var realizedResource *v1alpha1.RealizedResource

		if (stampedObject == nil || template == nil) && previousRealizedResource != nil {
			realizedResource = previousRealizedResource
		} else {
			realizedResource = generateRealizedResource(resource, template, stampedObject, out, previousRealizedResource)
		}

		healthCondition := metav1.Condition{}
		if template != nil { //TODO: honestly why should this be nil?!
			healthCondition = determineHealth(template.GetHealthRule(), realizedResource, stampedObject)
		}
		resourceStatuses.Add(realizedResource, err, healthCondition)
	}

	return firstError
}

func determineHealth(rule *v1alpha1.HealthRule, realizedResource *v1alpha1.RealizedResource, stampedObject *unstructured.Unstructured) metav1.Condition {
	//TODO: logging...? esp see below...
	if rule == nil {
		if len(realizedResource.Outputs) > 0 {
			return conditions.OutputAvailableResourcesHealthyCondition()
		} else {
			return conditions.OutputNotAvailableResourcesHealthyCondition()
		}
	}
	if rule.AlwaysHealthy != nil {
		return conditions.AlwaysHealthyResourcesHealthyCondition()
	}
	if rule.SingleConditionType != "" {
		result, err := utils.SinglePathEvaluate(fmt.Sprintf(".status.conditions[?(@.type==\"%s\")].status", rule.SingleConditionType), stampedObject)
		if err != nil { //TODO: logging?
			return conditions.SingleConditionTypeEvaluationErrorCondition(err)
		}
		if len(result) != 0 { //TODO: logging?
			return conditions.SingleConditionTypeNoResultResourcesCondition()
		}
		if conditionStatus, ok := result[0].(metav1.ConditionStatus); ok {
			return conditions.SingleConditionMatchCondition(conditionStatus, rule.SingleConditionType)
		}
	}
	if rule.MultiMatch != nil {
		if anyMultiMatch(rule.MultiMatch.Unhealthy, stampedObject) {
			return conditions.MultiMatchResourcesHealthyCondition(metav1.ConditionFalse)
		}
		if allMultiMatch(rule.MultiMatch.Healthy, stampedObject) {
			return conditions.MultiMatchResourcesHealthyCondition(metav1.ConditionTrue)
		}
	}
	return conditions.UnknownResourcesHealthyCondition()
}

func anyMultiMatch(rule v1alpha1.HealthMatchRule, stampedObject *unstructured.Unstructured) bool {
	for _, conditionRule := range rule.MatchConditions {
		result, err := utils.SinglePathEvaluate(fmt.Sprintf(".status.conditions[?(@.type==\"%s\")].status", conditionRule.Type), stampedObject)
		if err != nil { //TODO: logging?
			//probably not err like return conditions.SingleConditionTypeEvaluationErrorCondition(err)
		}
		if len(result) != 0 { //TODO: logging?
			//probably not err err like  return conditions.SingleConditionTypeNoResultResourcesCondition()
		}
		if conditionStatus, ok := result[0].(metav1.ConditionStatus); ok {
			if conditionStatus == conditionRule.Status {
				return true
			}
		}

	}
	for _, matchFieldRule := range rule.MatchFields {
		matches, err := selector.Matches(matchFieldRule.FieldSelectorRequirement, stampedObject)
		if err != nil {
			//TODO: logging? prob not err right?
		}
		if matches {
			return true
		}
	}
	return false
}

func allMultiMatch(rule v1alpha1.HealthMatchRule, stampedObject *unstructured.Unstructured) bool {
	for _, conditionRule := range rule.MatchConditions {
		result, err := utils.SinglePathEvaluate(fmt.Sprintf(".status.conditions[?(@.type==\"%s\")].status", conditionRule.Type), stampedObject)
		if err != nil { //TODO: logging?
			//probably not err like return conditions.SingleConditionTypeEvaluationErrorCondition(err)
			return false
		}
		if len(result) != 0 { //TODO: logging?
			//probably not err err like  return conditions.SingleConditionTypeNoResultResourcesCondition()
			return false
		}
		if conditionStatus, ok := result[0].(metav1.ConditionStatus); ok {
			if conditionStatus != conditionRule.Status {
				return false
			}
		} else {
			return false
		}

	}
	for _, matchFieldRule := range rule.MatchFields {
		matches, err := selector.Matches(matchFieldRule.FieldSelectorRequirement, stampedObject)
		if err != nil {
			//TODO: logging? prob not err right?
			return false
		}
		if !matches {
			return false
		}
	}
	return true
}

func generateRealizedResource(resource OwnerResource, template templates.Template, stampedObject *unstructured.Unstructured, output *templates.Output, previousRealizedResource *v1alpha1.RealizedResource) *v1alpha1.RealizedResource {
	if previousRealizedResource == nil {
		previousRealizedResource = &v1alpha1.RealizedResource{}
	}

	var inputs []v1alpha1.Input
	for _, source := range resource.Sources {
		inputs = append(inputs, v1alpha1.Input{Name: source.Resource})
	}

	for _, image := range resource.Images {
		inputs = append(inputs, v1alpha1.Input{Name: image.Resource})
	}

	if resource.Deployment != nil {
		inputs = append(inputs, v1alpha1.Input{Name: resource.Deployment.Resource})
	}

	for _, config := range resource.Configs {
		inputs = append(inputs, v1alpha1.Input{Name: config.Resource})
	}

	var templateRef *corev1.ObjectReference
	var outputs []v1alpha1.Output
	if template != nil {
		templateRef = &corev1.ObjectReference{
			Kind:       template.GetKind(),
			Name:       template.GetName(),
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		}

		outputs = getOutputs(template, previousRealizedResource, output)
	}

	var stampedRef *corev1.ObjectReference
	if stampedObject != nil {
		stampedRef = &corev1.ObjectReference{
			Kind:       stampedObject.GetKind(),
			Namespace:  stampedObject.GetNamespace(),
			Name:       stampedObject.GetName(),
			APIVersion: stampedObject.GetAPIVersion(),
		}
	}

	return &v1alpha1.RealizedResource{
		Name:        resource.Name,
		StampedRef:  stampedRef,
		TemplateRef: templateRef,
		Inputs:      inputs,
		Outputs:     outputs,
	}
}

func getOutputs(template templates.Template, previousRealizedResource *v1alpha1.RealizedResource, output *templates.Output) []v1alpha1.Output {
	outputs, err := template.GenerateResourceOutput(output)
	if err != nil {
		outputs = previousRealizedResource.Outputs
	} else {
		currTime := metav1.NewTime(time.Now())
		for j, out := range outputs {
			outputs[j].LastTransitionTime = currTime
			for _, previousOutput := range previousRealizedResource.Outputs {
				if previousOutput.Name == out.Name {
					if previousOutput.Digest == out.Digest {
						outputs[j].LastTransitionTime = previousOutput.LastTransitionTime
					}
					break
				}
			}
		}
	}

	return outputs
}
