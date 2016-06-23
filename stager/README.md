# Staging apps on Docker instead of Diego

## Overview

This experiment attempts to stage and run CF applications on Docker. It piggy-backs on the existing Diego cc-bridge while the version targetting k8s is still under development.

## Changes to the default bosh-lite setup ##

This experiment depends on a working bosh-lite installation of cf-release, diego-release, and [kubernetes-release](https://github.com/cf-furnace/kubernetes-release). The latter in turn requires a consul-release to be uploaded before it can be deployed.

In addition the following modifications should be made. Note that the cc-bridge components are now (CFv237) part of `capi-release`, which is included in `cf-release`. The deployment manifest is however still included in `diego-release`. So after making the changes given below, create and upload a new `cf-release`, but also deploy the `diego-release` again to get the cc-bridge components recompiled.

### cf-release deployment manifest

For convenience the `cc.default_to_diego_backend` property in `bosh-lite/deployments/cf.yml` should be set to `true` to avoid having to deal with the [Diego-Enabler](https://github.com/cloudfoundry-incubator/Diego-Enabler) on each app push:

```yaml
    default_to_diego_backend: true
```

### Stager component

This small change to the `stager` adds the build requests to the log file:

```diff
--- a/backend/buildpack_backend.go
+++ b/backend/buildpack_backend.go
@@ -214,7 +214,7 @@ func (backend *traditionalBackend) BuildRecipe(stagingGuid string, request cc_me
                TrustedSystemCertificatesPath: TrustedSystemCertificatesPath,
        }

-       logger.Debug("staging-task-request")
+       logger.Info("staging-task-request", lager.Data{"taskDefinition": taskDefinition})

        return taskDefinition, stagingGuid, backend.config.TaskDomain, nil
 }
```

### NSync component

The corresponding change to `nsync` records the launch request as well:

```diff
--- a/recipebuilder/buildpack_recipe_builder.go
+++ b/recipebuilder/buildpack_recipe_builder.go
@@ -234,7 +234,7 @@ func (b *BuildpackRecipeBuilder) Build(desiredApp *cc_messages.DesireAppRequestF
        setupAction := models.Serial(setup...)
        actionAction := models.Codependent(actions...)

-       return &models.DesiredLRP{
+       desiredLRP := &models.DesiredLRP{
                Privileged: true,

                Domain: cc_messages.AppLRPDomain,
@@ -272,7 +272,9 @@ func (b *BuildpackRecipeBuilder) Build(desiredApp *cc_messages.DesireAppRequestF

                TrustedSystemCertificatesPath: TrustedSystemCertificatesPath,
                VolumeMounts:                  desiredApp.VolumeMounts,
-       }, nil
+       }
+       buildLogger.Info("desired-app-request", lager.Data{"desiredLRP": desiredLRP})
+       return desiredLRP, nil
 }

 func (b BuildpackRecipeBuilder) ExtractExposedPorts(desiredApp *cc_messages.DesireAppRequestFromCC) ([]uint32, error) {
```

## Testing everything works

All the commands needed here can be sourced from the `helpers` file. That file also includes additional documentation for each command.

```bash
. helpers
```

Everything is being run from the local host via bosh.

First push an application using the normal `cf` command. Once finished, run `build_req` and `launch_req` to make sure the corresponding staging task and launch LRP can be extracted from the stager and nsync logs. If this doesn't work, then the patches above have not been deployed correctly, or the backend default hasn't be switched to Diego.

To build the stager run the `rebuild` command. It will create both a Docker registry and a stager docker image inside the `kube-node/0` node.

The `restage` command will now re-run the last staging command with the newly generated stager image. It will push the generated droplet back to the cc0uploader, and will also create a docker image for that app on `kube-node/0`. This app will both be stored locally on the node, and also pushed to the registry. For convenience it uses the app name and not a GUID for the image right now.

The generated app image can be launched with the `launch` command inside `kube-node/0`, or with `launch_local` on the local machine, assuming it has Docker for Mac installed. The former command will re-use the image already stored inside the node while the latter will fetch the image from the docker registry to the local host. Both commands will then launch a browser to view the app.

The app will always use port 8080, and only one app can be running at a time (actually 2, if you run one locally and the other on `kube-node/0`.

Running the `cache_update` command after pushing an app using auto-detect will fetch all buildpacks and put them onto the stager image.  This simulates the local caching done by Diego on the host.
