package cluster

import (
	"context"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	apiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pkgclient "github.com/giantswarm/awscnfm/v15/pkg/client"
	"github.com/giantswarm/awscnfm/v15/pkg/env"
	"github.com/giantswarm/awscnfm/v15/pkg/key"
	"github.com/giantswarm/awscnfm/v15/pkg/label"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	var cpClients k8sclient.Interface
	{
		c := pkgclient.ControlPlaneConfig{
			Logger: r.logger,

			KubeConfig: env.ControlPlaneKubeConfig(),
		}

		cpClients, err = pkgclient.NewControlPlane(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		err := cpClients.CtrlClient().DeleteAllOf(
			ctx,
			&apiv1alpha3.Cluster{},
			client.InNamespace(key.OrganizationNamespaceFromName(key.Organization)),
			client.MatchingLabels{label.Cluster: r.flag.TenantCluster},
		)
		if errors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
