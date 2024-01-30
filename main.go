package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blang/semver"
	random "github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// to destroy our program, we can run `go run main.go destroy`
	destroy := false
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "destroy" {
			destroy = true
		}
	}

	deployFunc := func(ctx *pulumi.Context) error {
		pet, err := random.NewRandomPet(ctx, "Fluffy", &random.RandomPetArgs{})
		if err != nil {
			return err
		}
		ctx.Export("pet.id", pet.ID())
		return nil
	}

	ctx := context.Background()

	projectName := "auto-install-example"
	stackName := "dev"

	// Note that if Pulumi is already installed in `Root` this checks that the
	// version is compatible and skips the download.
	pulumiCommand, err := auto.InstallPulumiCommand(ctx, &auto.PulumiCommandOptions{
		// Version defaults to the version matching the current pulumi/sdk.
		// Since this example is using a pre-release version of the SDK, we
		// have to specify the version explicitly.
		Version: semver.MustParse("3.102.0"),
		// Where to install the Pulumi CLI, defaults to $HOME/.pulumi/versions/$VERSION
		Root: ".pulumi",
	})
	if err != nil {
		fmt.Printf("failed to install pulumi: %s\n", err)
		os.Exit(1)
	}

	// Pass the PulumiCommand instance that we just installed as a workspace option.
	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployFunc, auto.Pulumi(pulumiCommand))
	if err != nil {
		fmt.Printf("Failed to upsert stack: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created/Selected stack %q\n", stackName)

	if destroy {
		fmt.Println("Starting stack destroy")
		stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
		_, err := s.Destroy(ctx, stdoutStreamer)
		if err != nil {
			fmt.Printf("Failed to destroy stack: %v", err)
		}
		fmt.Println("Stack successfully destroyed")
		os.Exit(0)
	}

	fmt.Println("Starting update")

	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Println("Update succeeded!")

	petID, ok := res.Outputs["pet.id"].Value.(string)
	if !ok {
		fmt.Println("Failed to unmarshall")
		os.Exit(1)
	}

	fmt.Printf("Pet ID: %s\n", petID)
}
