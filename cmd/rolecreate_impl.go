package cmd

import (
	"github.com/akolb1/sentrytool/sentryapi"
	"fmt"
)

func roleCreate(host string, port int, user string, verbose bool, names []string) error  {
	client, err := sentryapi.GetClient(host, port, user)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, roleName := range names {
		err = client.CreateRole(roleName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if verbose {
			fmt.Println("created role ", roleName)
		}
	}
	return nil
}
