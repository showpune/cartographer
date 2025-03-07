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

package cmd

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"github.com/vmware-tanzu/cartographer/pkg/controllers"
	"github.com/vmware-tanzu/cartographer/pkg/utils"

	"os"
)

type Command struct {
	Port                       int
	CertDir                    string
	Logger                     logr.Logger
	WorkNamespace              string
	workedForNamespaceResource bool
	workedForClusterResource   bool
}

const (
	kappctrlWorkNamespaceEnvKey = "KAPPCTRL_WORK_NAMESPACE"
	cluster                     = "cluster"
)

func (cmd *Command) Execute(ctx context.Context) error {
	log.SetLogger(cmd.Logger)
	l := log.Log.WithName("cartographer")

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	scheme := runtime.NewScheme()
	if err := utils.AddToScheme(scheme); err != nil {
		return fmt.Errorf("add to scheme: %w", err)
	}

	cmd.workedForNamespaceResource = true
	cmd.workedForClusterResource = true
	var namespaceWorkedFor = ""
	if namespaceWorkedFor, ok := os.LookupEnv(kappctrlWorkNamespaceEnvKey); ok {
		if namespaceWorkedFor == cluster {
			cmd.workedForNamespaceResource = false
			cmd.workedForClusterResource = true
		} else if namespaceWorkedFor != "" {
			cmd.workedForNamespaceResource = true
			cmd.workedForClusterResource = false
		}
		l.Info("Start for work namespace", "namespace", namespaceWorkedFor)
	} else {
		l.Info("Start for all namespace", "namespace", namespaceWorkedFor)
	}

	mgr, err := manager.New(cfg, manager.Options{
		Port:               cmd.Port,
		CertDir:            cmd.CertDir,
		Scheme:             scheme,
		MetricsBindAddress: "0",
		Namespace:          namespaceWorkedFor,
	})
	if err != nil {
		return fmt.Errorf("failed to create new manager: %w", err)
	}

	if err := registerControllers(mgr, cmd); err != nil {
		return fmt.Errorf("failed to register controllers: %w", err)
	}

	if cmd.CertDir != "" {
		if err := registerWebhooks(mgr, cmd); err != nil {
			return fmt.Errorf("failed to register webhooks: %w", err)
		}
	} else {
		l.Info("Not registering the webhook server. Must pass a directory containing tls.crt and tls.key to --cert-dir")
	}

	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("manager start: %w", err)
	}

	return nil
}

func registerControllers(mgr manager.Manager, cmd Command) error {
	if cmd.workedForNamespaceResource {
		if err := (&controllers.WorkloadReconciler{}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("failed to register workload controller: %w", err)
		}

		if err := (&controllers.DeliverableReconciler{}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("failed to register deliverable controller: %w", err)
		}

		if err := (&controllers.RunnableReconciler{}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("failed to register runnable controller: %w", err)
		}
	}
	if cmd.workedForClusterResource {
		if err := (&controllers.SupplyChainReconciler{}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("failed to register supply chain controller: %w", err)
		}

		if err := (&controllers.DeliveryReconiler{}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("failed to register delivery controller: %w", err)
		}
	}

	return nil
}

func registerWebhooks(mgr manager.Manager, cmd Command) error {
	if cmd.workedForClusterResource {
		if err := (&v1alpha1.ClusterSupplyChain{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster supply chain webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterDelivery{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster delivery webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterConfigTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster config template webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterDeploymentTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster deployment template webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterImageTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster image template webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterRunTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster run template webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterSourceTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster source template webhook: %w", err)
		}

		if err := (&v1alpha1.ClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("failed to setup cluster template webhook: %w", err)
		}
	}

	return nil
}
