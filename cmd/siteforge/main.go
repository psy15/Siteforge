package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"siteforge/internal/kube"
	"siteforge/internal/platform"
)

func main() {
	var (
		tenant  = flag.String("tenant", "", "Tenant (namespace) name")
		runtime = flag.String("runtime", "wordpress", "Runtime to deploy (wordpress)")
	)
	flag.Parse()

	if *tenant == "" {
		fmt.Fprintln(os.Stderr, "error: --tenant is required")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := kube.NewClient()
	if err != nil {
		fail(err)
	}

	if err := platform.EnsureTenant(ctx, client, *tenant); err != nil {
		fail(err)
	}

	if err := platform.EnsureRuntime(ctx, client, *tenant, *runtime); err != nil {
		fail(err)
	}

	fmt.Println("Tenant and runtime ready:", *tenant)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
