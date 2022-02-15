# mygithubapp

A simple Go web application that will connect to Github, connect to a repo, create a branch and create a Pull request on the newly created branch for the repo by modifying some of the files.

It is using OAuth 2.0 based implementation to authenticate the user with Github (user authenticates with Github, not with the application). This WebApp requests Github for access
on the repositories, User authenticates with GitHub, and App operates the repository on behalf of the user.

I have added files like Dockerfile, Makefile and docker-compose.yml to host this WebApp on a container.
