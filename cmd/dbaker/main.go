package main

import "github.com/spf13/cobra"

func main() {
	dbakerCommand := cobra.Command{
		Use:   "dbaker",
		Short: "Fake data generator (DB + Faker = DBaker)",
		Long:  "Introspect live database instance, generate & write fake data right back into the instance",
	}
	dbakerCommand.AddCommand(newIntrospecCommand(), newGenerateCommand())

	dbakerCommand.Execute()
}
