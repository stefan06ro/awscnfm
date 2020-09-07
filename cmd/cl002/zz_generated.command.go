package cl002

import (
	"io"
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/awscnfm/v12/cmd/cl002/ac000"
	"github.com/giantswarm/awscnfm/v12/cmd/cl002/ac001"
	"github.com/giantswarm/awscnfm/v12/cmd/cl002/ac002"
)

const (
	name        = "cl002"
	description = "Conformance tests for cluster scope cl002."
)

type Config struct {
	Logger micrologger.Logger
	Stderr io.Writer
	Stdout io.Writer
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}

	var err error

	var ac000Cmd *cobra.Command
	{
		c := ac000.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		ac000Cmd, err = ac000.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var ac001Cmd *cobra.Command
	{
		c := ac001.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		ac001Cmd, err = ac001.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var ac002Cmd *cobra.Command
	{
		c := ac002.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		ac002Cmd, err = ac002.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  description,
		RunE:  r.Run,
	}

	f.Init(c)

	c.AddCommand(ac000Cmd)
	c.AddCommand(ac001Cmd)
	c.AddCommand(ac002Cmd)

	return c, nil
}
