# Demo Keycloak project

## Requirements

To run this demo in read-only mode you need only [docker-compose](https://docs.docker.com/compose/).

You will need [npm](https://www.npmjs.com/) and [poetry](https://python-poetry.org/)
to run frontend and backend on your local machine (without docker).

## Description

The goal of this project to be a "hello world" for authentication with keycloak.
The project consists of three apps: fronted, backend and keycloak. You can run
them on your local machine or inside docker. Run `make help` for more information.

**TL;DR: Start the stack `make docker` then open http://localhost:8000/ and
login as alice (password: `alice`)**

The Keycloak service is provided with a predefined configuration and realm.
To simplify the project we use h2 storage. Keycloak data is stored in a
persistent volume, so you can safely restart the container.

Available services:

- **keycloak** - http://localhost:8080: user/password: `admin/admin`

- **frontend** - http://localhost:8000: user/password:
`alice/alice` (_admin_ role), `bob/bob` (_user_ role)

- **backend_py** - http://localhost:3000:
[/admin](http://localhost:3000/admin) (require _admin_ role),
[/user](http://localhost:3000/user) (require _user_ role) ,
[/service](http://localhost:3000/service) (no roles required)

- **backend_go** - http://localhost:3333:
[/service](http://localhost:3333/service) (require _service_ role)

**myrealm** has _**frontend**_ and _**backed_py**_ clients. Available users are:

- _**alice**_ (password `alice`) — user with **admin** and **user** roles
- _**bob**_ (password `bob`) — user with **user** role
- _**service**_ (access throw key) — user with **service** role

**backend_py** has Service Account in the Keycloak and works as a client for **backend_go**.

Feel free to experiment with Keycloak you are able to reset to "factory defaults" with `make clean`.

## Troubleshooting

- `keycloak_1 | User with username 'admin' already added to` - run **`make clean`**.

## Development Notes

### Register Client

```shell
docker-compose exec keycloak bash

cd /opt/jboss/keycloak/bin/

# Login as admin
./kcadm.sh config credentials --server http://localhost:8080/auth --realm master --user admin --password admin

# Create client for backend_py
./kcadm.sh create clients -r myrealm -s clientId=backend_py -s enabled=true -s clientAuthenticatorType=client-secret -s secret=00000000–0000–0000–0000–000000000000

# enable the service account
./kcadm.sh update clients/61719f60-c170-4601-9ff9-6b9fc566e09b -r myrealm -s 'redirectUris=["*"]'  -s serviceAccountsEnabled=true

# add `service` role to Service account
# ./kcadm.sh add-roles --uusername service-account-<CLIENT-ID> --rolename <ROLE-NAME> -r <REALM-NAME>
./kcadm.sh add-roles --uusername service-account-backend_py --rolename service -r myrealm
```

## References

- [Keycloak - Identity and Access Management for Modern Applications](https://github.com/PacktPublishing/Keycloak-Identity-and-Access-Management-for-Modern-Applications) (book)

- [Building an effective identity and access management architecture with Keycloak](https://youtu.be/RupQWmYhrLA) (youtube)

- [FastAPI -Simple OAuth2](https://fastapi.tiangolo.com/tutorial/security/simple-oauth2/)

- [Getting Started with Service Accounts in Keycloak](https://medium.com/@mihirrajdixit/getting-started-with-service-accounts-in-keycloak-c8f6798a0675)

- [Build and Secure an API in Python with FastAPI ](https://developer.okta.com/blog/2020/12/17/build-and-secure-an-api-in-python-with-fastapi)
