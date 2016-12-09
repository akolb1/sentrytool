package sentryapi

import (
	"flag"
	"os"
	"os/user"
	"strconv"
	"testing"
	"fmt"
)

const (
	defaultHost = "localhost"
	defaultPort = "8038"
)

var client ClientAPI

func TestMain(m *testing.M) {
	flag.Parse()
	host := os.Getenv("SENTRY_HOST")
	port := os.Getenv("SENTRY_PORT")
	userName := os.Getenv("SENTRY_USER")

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
	client, err = GetClient(PolicyProtocol, host, portVal, "", userName)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func ExampleSentryClient_CreateRole() {
	roleName := "exampleRole"
	err := client.CreateRole(roleName)
	if err != nil {
		fmt.Println("failed to create role: ", err)
		os.Exit(1);
	}
	fmt.Println("Created role", roleName)
	err = client.RemoveRole(roleName)
	if err != nil {
		fmt.Println("failed to create role: ", err)
		os.Exit(1);
	}
	fmt.Println("Removed role", roleName)
	// Output:
	// Created role exampleRole
	// Removed role exampleRole
}

func TestSentryClient_CreateRole(t *testing.T) {
	roleName := "sentryTestRole"
	err := client.CreateRole(roleName)
	if err != nil {
		t.Error(err)
	}
	err = client.RemoveRole(roleName)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkSentryClient_CreateRole(b *testing.B) {
	roleName := "testRole"
	for i := 0; i < b.N; i++ {
		err := client.CreateRole(roleName)
		if err != nil {
			b.Fail()
		}
		err = client.RemoveRole(roleName)
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkSentryClient_CreateRoles(b *testing.B) {
	roleName := "testRole"
	for i := 0; i < b.N; i++ {
		name := fmt.Sprint("%s_%d", roleName, i)
		err := client.CreateRole(name)
		if err != nil {
			b.Fail()
		}
	}
	for i := 0; i < b.N; i++ {
		name := fmt.Sprint("%s_%d", roleName, i)
		err := client.RemoveRole(name)
		if err != nil {
			continue
			b.Errorf("can't remove %s: %v", name, err)
		}
	}
}