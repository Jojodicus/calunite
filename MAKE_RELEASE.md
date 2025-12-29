# Making a Release

These are the steps on how to make a release and publish them to Docker Hub.

## Check if everything works

Build container locally [see README.md](/README.md#️-development).

Test old and new features:
1. Webserver starts correctly
2. Merging works
3. Importing calendar works (import on phone from local IP/domain)
4. YAML options work (title, event formatting, private calendars, ...)
5. ...

## Make the Release

Go to the [releases page](https://github.com/Jojodicus/calunite/releases) and draft a new release.

Create a new tag with the version number, i.e. `1.6`.  
Release semantics: `major.minor[.hotfix]`

Title: `CalUnite <tag>`  
Content:
```md
`docker pull jojodicus/calunite:<tag>`

New:
- List features
- #2 Reference issues or pull requests like this

Optional footer for further information. Referencing users with @user will make them show up as collaborators
```
Click "Generate release notes" to generate the full changelog diff.

After publishing the release, Github Actions will automatically build the container and push it to Docker Hub.
You can check the status in the [Actions overview](https://github.com/Jojodicus/calunite/actions).

## Check Docker Hub

If the README.md has changed, update it on [Docker Hub](https://hub.docker.com/r/jojodicus/calunite) as well.
As some links change, use the `genDockerReadme.sh` helper script:

```sh
./genDockerReadme.sh | xclip # for x11
./genDockerReadme.sh | wl-copy # for wayland
```

Delete the old description and paste the new one.

After the push has been completed, check Docker Scout for findings.
