package action

import (
	"dbaker/pkg/adapter"
	"dbaker/pkg/config"
	"dbaker/pkg/model"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type introspect struct {
	config  config.Config
	adapter adapter.PostgreSQLAdapter
}

func NewIntrospect(config config.Config, adapter adapter.PostgreSQLAdapter) *introspect {
	return &introspect{
		config,
		adapter,
	}
}

func (i *introspect) Execute() error {
	err := i.adapter.Init()
	if err != nil {
		return err
	}
	defer i.adapter.Close()

	var tables []*model.Table
	for _, tbl := range i.config.Tables {
		tableName, schema := splitTableName(tbl)
		if tableName == "" || schema == "" {
			return fmt.Errorf("provided invalid table name: %s", tbl)
		}

		fmt.Printf("Introspecting table: %s ...\n", tbl)

		table, err := i.adapter.IntrospectTable(tableName, schema)
		if err != nil {
			return err
		}

		tables = append(tables, table)
		fmt.Println("done.")
	}

	recipeFilePath := fmt.Sprintf("./%s.recipe.json", i.config.Database)
	if err := writeJson(recipeFilePath, tables); err != nil {
		return err
	}

	fmt.Printf("Baking recepi written to %s\n", recipeFilePath)

	return nil
}

func writeJson(filePath string, tables []*model.Table) error {
	json, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func splitTableName(table string) (name string, schema string) {
	splits := strings.Split(table, ".")
	if len(splits) != 2 {
		return "", ""
	}

	name = splits[1]
	schema = splits[0]

	return
}
