package cmd

import (
	"github.com/akolb1/sentrytool/sentryapi"
	"fmt"
	"sort"
)

func roleList(host string, port int, user string) error  {
	client, err := sentryapi.GetClient(host, port, user)
	if err != nil {
		return err
	}
	defer client.Close()

	roles, err := client.ListRoleByGroup("")
	if err != nil {
		return err
	}
	sort.Strings(roles)
	for _, r := range roles {
		fmt.Println(r)
	}
	return nil
}
