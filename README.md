# Water API

## Local Development

1. Type `docker compose up` to bring `db`,`api`,`flyway`,`pgadmin` containers online.

## Deployments to Cloud Environment

### Deploying to Develop

Commits into branch `main` cause new automated container image builds, tagged with `:latest`. Upon successful build and tagging, the new image is automatically pushed to a container registry.

Pushing a new image with tag `:latest` to the container registry triggers a new deployment of the `develop` service. New updated containers are launched and attached to the load balancer. Upon successful deployment and passing health checks, the prior containers are automatically detatched from the load balancer and decomissioned.

Making commits in branch `main` is done via a pull request workflow. This means changes should be made in a separate, focused branch, created from branch main. Creating a feature branch is accomplished (starting from branch `main`) with the command `git checkout -b feature/<my-feature-branch>`. After changes are made and committed in the feature branch, push the feature branch to github using the command `git push origin feature/<my-feature-branch>`. Open a pull request targeting branch `main`. Opening this pull request will cause automated tests to run. Once all tests are passing and mandatory code review is complete, the branch can be merged. Merging feature branches to `main` should be done with a Squash and Merge. This ensures that all new commits in the feature branch since branching from `main` are squashed into a single commit in the commit history in `main`.

### Deploying to Stable

Deployments to stable are made by creating and pushing a new tag that has the naming pattern corresponding to semantic versioning `vX.X.X`, where `X` is any number (e.g. `v0.1.0`, `v2.1.2`, `v15.2.1`). This can be done using these commands:

- `git checkout main` (to checkout main branch)
- `git pull origin main` (to ensure you have the latest changes from github)
- (switch to correct commit; If deploying the current state of `main`, skip this step)
- `git tag` (to list existing tags)
- `git tag v0.2.0` (create a tag pointing to current checked out commit in `main`. Note: `v0.2.0` is just an example, change as needed)
- `git push origin v0.2.0` (push the tag to github. This will kick-off a deployment to stable)