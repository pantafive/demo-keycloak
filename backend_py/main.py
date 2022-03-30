from __future__ import annotations

import logging
from typing import Any
from urllib.parse import urljoin

import uvicorn
from auth import RequireRoles
from auth import auth
from client import KeycloakAuth
from client import ServiceClient
from config import backend_go_url
from config import keycloak_client_id
from config import keycloak_client_secret
from fastapi import Depends
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from httpx import AsyncClient

logger = logging.getLogger(__name__)
logger.setLevel("INFO")


keycloak_auth = KeycloakAuth(keycloak_client_id, keycloak_client_secret, 290)
async_client = AsyncClient(auth=keycloak_auth)
client = ServiceClient(async_client)


app = FastAPI()

app.add_middleware(
    CORSMiddleware, allow_origins=["*"], allow_credentials=True, allow_methods=["*"], allow_headers=["*"]
)


@app.on_event("startup")
async def startup_event() -> None:
    await auth.update_keycloak_public_key()
    await keycloak_auth.start_auto_update()


@app.on_event("shutdown")
async def shutdown_event() -> None:
    await async_client.aclose()


@app.get("/user")
async def user(token: dict[str, Any] = Depends(RequireRoles(auth, "user"))) -> dict[str, str]:
    return token


@app.get("/admin")
async def admin(token: dict[str, Any] = Depends(RequireRoles(auth, "admin"))) -> dict[str, str]:
    return token


@app.get("/service")
async def service() -> Any:
    return await client.call_server(urljoin(backend_go_url, "service"))


if __name__ == "__main__":
    uvicorn.run(app, port=3000, loop="asyncio")
