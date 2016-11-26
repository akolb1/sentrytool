package cmd

import (
	"github.com/akolb1/sentrytool/sentryapi"
	"fmt"
)

func roleDelete(host string, port int, user string, verbose bool, names []string) error  {
	client, err := sentryapi.GetClient(host, port, user)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, roleName := range names {
		err = client.RemoveRole(roleName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if verbose {
			fmt.Println("removed role ", roleName)
		}
	}
	return nil
}
