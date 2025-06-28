package main

import (
	"dbaker/pkg/action"
	"dbaker/pkg/adapter"
	"dbaker/pkg/config"

	"github.com/spf13/cobra"
)

func newIntrospecCommand() *cobra.Command {
	var config config.Config
	introspectCmd := cobra.Command{
		Use:     "introspect",
		Aliases: []string{"i"},
		Short:   "Introspect database for data gen.",
		Long:    "Introspect a live database instance and create intermediate representation for data gen.",
		RunE: func(_ *cobra.Command, _ []string) error {
			pgAdapter := adapter.NewPostgreSQLAdapter(config)
			action := action.NewIntrospect(config, pgAdapter)

			return action.Execute()
		},
	}

	introspectCmd.Flags().StringVarP(&config.Host, "host", "H", "", "host of the db to introspect")
	introspectCmd.Flags().UintVarP(&config.Port, "port", "P", 5432, "port of the db to introspect")
	introspectCmd.Flags().StringVarP(&config.Database, "database", "d", "", "database (pg) to connect to")
	introspectCmd.Flags().StringVarP(&config.Username, "username", "u", "", "database user")
	introspectCmd.Flags().StringVarP(&config.Password, "password", "p", "", "database user")
	introspectCmd.Flags().StringArrayVarP(&config.Tables, "tables", "t", []string{}, "tables to include in the introspection")

	introspectCmd.MarkFlagRequired("host")
	introspectCmd.MarkFlagRequired("database")
	introspectCmd.MarkFlagRequired("username")
	introspectCmd.MarkFlagRequired("password")
	introspectCmd.MarkFlagRequired("tables")

	return &introspectCmd
}
