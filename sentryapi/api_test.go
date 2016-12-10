package sentryapi

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"testing"
)

const (
	defaultHost = "localhost"
	defaultPort = "8038"
	component   = "solr"
	hostEnv     = "SENTRY_HOST"
	portEnv     = "SENTRY_PORT"
	userEnv     = "SENTRY_USER"
)

var (
	client        ClientAPI
	genericClient ClientAPI
)

func TestMain(m *testing.M) {
	// Read config from environment vars and create Hive and
	// generic clients.
	flag.Parse()
	host := os.Getenv(hostEnv)
	port := os.Getenv(portEnv)
	userName := os.Getenv(userEnv)

	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	if userName == "" {
		currentUser, _ := user.Current()
		userName = currentUser.Username
	}

	portVal, err := strconv.Atoi(port)
	if err != nil {
		panic("invalid port" + port)
	}

	// Create Hive protocol client
	client, err = GetClient(PolicyProtocol, host, portVal, "", userName)
	if err != nil {
		panic(err)
	}
	// Create generic protocol client
	genericClient, err = GetClient(PolicyProtocol, host, portVal, component, userName)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func createRemoveRole(client ClientAPI, roleName string) error {
	if err := client.CreateRole(roleName); err != nil {
		return fmt.Errorf("failed to create role %s: %v", roleName, err)
	}
	if err := client.RemoveRole(roleName); err != nil {
		return fmt.Errorf("failed to delete role %s: %v", roleName, err)

	}
	return nil
}

func ExampleClientAPI_CreateRole() {
	roleName := "exampleRole"
	err := client.CreateRole(roleName)
	if err != nil {
		fmt.Println("failed to create role: ", err)
		os.Exit(1)
	}
	fmt.Println("Created role", roleName)
	err = client.RemoveRole(roleName)
	if err != nil {
		fmt.Println("failed to create role: ", err)
		os.Exit(1)
	}
	fmt.Println("Removed role", roleName)
	// Output:
	// Created role exampleRole
	// Removed role exampleRole
}

func TestSentryClient_CreateRole(t *testing.T) {
	if err := createRemoveRole(client, "sentryTestRole"); err != nil {
		t.Error(err)
	}
}

// Repeatedly create and remove a role
func BenchmarkSentryClient_CreateRole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := createRemoveRole(client, "testRole"); err != nil {
			b.Fail()
		}
	}
}

// Repeatedly create roles then repeatedly remove them
func BenchmarkSentryClient_CreateRoles(b *testing.B) {
	roleName := "testRole"
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("%s_%d", roleName, i)
		err := client.CreateRole(name)
		if err != nil {
			b.Fail()
		}
	}
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("%s_%d", roleName, i)
		err := client.RemoveRole(name)
		if err != nil {
			b.Errorf("can't remove %s: %v", name, err)
			continue
		}
	}
}

func TestGenericSentryClient_CreateRole(t *testing.T) {
	if err := createRemoveRole(genericClient,
		"sentryGenTestRole"); err != nil {
		t.Error(err)
	}
}

// Repeatedly create and remove a role
func BenchmarkGenericSentryClient_CreateRole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := createRemoveRole(genericClient, "testRole"); err != nil {
			b.Fail()
		}
	}
}

// Repeatedly create roles then repeatedly remove them
func BenchmarkGenericSentryClient_CreateRoles(b *testing.B) {
	roleName := "testRole"
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("%s_%d", roleName, i)
		err := genericClient.CreateRole(name)
		if err != nil {
			b.Fail()
		}
	}
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("%s_%d", roleName, i)
		err := genericClient.RemoveRole(name)
		if err != nil {
			b.Errorf("can't remove %s: %v", name, err)
			continue
		}
	}
}

func BenchmarkSentryClient_ListPrivilegesByRole(b *testing.B) {
	roleName := "testRole"
	err := client.CreateRole(roleName)
	if err != nil {
		b.Fail()
	}

	for i := 0; i < b.N; i++ {
		_, err := client.ListPrivilegesByRole(roleName, nil)
		if err != nil {
			b.Error("can't list privileges")
			break
		}
	}

	err = client.RemoveRole(roleName)
	if err != nil {
		b.Fail()
	}
}

func BenchmarkSentryClient_GrantAndRevokePrivilege(b *testing.B) {
	roleName := "testRole"
	err := client.CreateRole(roleName)
	if err != nil {
		b.Fail()
	}
	priv := Privilege{
		Server:      "some_server",
		Database:    "some_database",
		Table:       "some_table",
		Column:      "some_column",
		URI:         "some_uri",
		Action:      "all",
		GrantOption: true,
	}

	for i := 0; i < b.N; i++ {
		err := client.GrantPrivilege(roleName, &priv)
		if err != nil {
			b.Errorf("can't grant privilege: %v", err)
			break
		}
		err = client.RevokePrivilege(roleName, &priv)
		if err != nil {
			b.Errorf("can't revoke privilege: %v", err)
			break
		}
	}

	err = client.RemoveRole(roleName)
	if err != nil {
		b.Fail()
	}
}
