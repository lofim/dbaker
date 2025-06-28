package main

import (
	"dbaker/pkg/action"
	"dbaker/pkg/adapter"
	"dbaker/pkg/config"

	"github.com/spf13/cobra"
)

func newGenerateCommand() *cobra.Command {
	var config config.Config
	introspectCmd := cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate fake data",
		Long:    "Generate fake data and write them directly into the live database instance",
		RunE: func(_ *cobra.Command, _ []string) error {
			pgAdapter := adapter.NewPostgreSQLAdapter(config)
			action := action.NewGenerate(config, pgAdapter)

			return action.Execute()
		},
	}

	introspectCmd.Flags().StringVarP(&config.Host, "host", "H", "", "host of the db to introspect")
	introspectCmd.Flags().UintVarP(&config.Port, "port", "P", 5432, "port of the db to introspect")
	introspectCmd.Flags().StringVarP(&config.Database, "database", "d", "", "database (pg) to connect to")
	introspectCmd.Flags().StringVarP(&config.Username, "username", "u", "", "database user")
	introspectCmd.Flags().StringVarP(&config.Password, "password", "p", "", "database user")
	introspectCmd.Flags().Uint32VarP(&config.DataSize, "size", "s", 0, "dataset size, number of rows to generate")
	introspectCmd.Flags().Uint32VarP(&config.IterFrom, "iterFrom", "i", 0, "iteration index from which to start generating unique values")

	introspectCmd.MarkFlagRequired("host")
	introspectCmd.MarkFlagRequired("database")
	introspectCmd.MarkFlagRequired("username")
	introspectCmd.MarkFlagRequired("password")
	introspectCmd.MarkFlagRequired("size")

	return &introspectCmd
}
