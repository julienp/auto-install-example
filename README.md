# Pulumi Auto-Install Example

This example shows how to use the [Automation API](https://www.pulumi.com/docs/using-pulumi/automation-api/) to install the Pulumi CLI.

```go
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
```
