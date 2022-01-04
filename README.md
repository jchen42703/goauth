# goauth

[WIP] Basic username password cookie based authentication with Go Lang

# Overview

- Use a Postgres DB to store Sign-in and Sign-up info
- Redis for caching

# Tutorials

https://www.sohamkamani.com/golang/password-authentication-and-storage/

https://www.sohamkamani.com/golang/session-based-authentication/

# Getting Started

The server relies on environment variables for secrets.

To load the environment variables from a `.env` to your environment:

```
env $(cat .env | xargs) go run main.go
```
