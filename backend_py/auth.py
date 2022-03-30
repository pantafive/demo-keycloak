from __future__ import annotations

import asyncio
import logging
import os
from typing import Any
from urllib.parse import urljoin

import httpx
from fastapi import Depends
from fastapi import HTTPException
from fastapi.security import OAuth2AuthorizationCodeBearer
from httpx import AsyncClient
from jose import JWTError
from jose import jwt
from starlette import status

APP_BASE_URL = "http://localhost:8000/"
KEYCLOAK_BASE_URL = os.getenv("KEYCLOAK_BASE_URL", "http://localhost:8080")
AUTH_URL = f"{KEYCLOAK_BASE_URL}/auth/realms/Clients/protocol/openid-connect/auth?client_id=app&response_type=code"
TOKEN_URL = f"{KEYCLOAK_BASE_URL}/auth/realms/Clients/protocol/openid-connect/token"

oauth2_scheme = OAuth2AuthorizationCodeBearer(authorizationUrl=AUTH_URL, tokenUrl=TOKEN_URL)

max_retries = 30


class Auth:
    def __init__(self, keycloak_base_url: str) -> None:
        self.keycloak_public_key_url = urljoin(keycloak_base_url, "/auth/realms/myrealm/")
        self._keycloak_public_key = ""

    async def __call__(self, *args: Any, **kwargs: Any) -> Auth:
        return self

    @property
    def keycloak_public_key(self) -> str:
        if not self._keycloak_public_key:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Keycloak public key not set."
            )
        return self._keycloak_public_key

    async def update_keycloak_public_key(self) -> None:
        client: AsyncClient
        for i in range(1, max_retries):
            async with httpx.AsyncClient() as client:
                try:
                    response = await client.get(self.keycloak_public_key_url)
                except httpx.RequestError as e:
                    logging.warning(
                        f"{e.__class__.__name__}: {i}/{max_retries} Keycloak not ready yet. Sleep for one second before retrying."
                    )
                    await asyncio.sleep(1)
                    continue

                public_key = response.json()["public_key"]
                self._keycloak_public_key = f"-----BEGIN PUBLIC KEY-----\n{public_key}\n-----END PUBLIC KEY-----"


auth = Auth(KEYCLOAK_BASE_URL)


class RequireRoles:
    def __init__(self, kc_auth: Auth = Depends(auth), *roles: str) -> None:
        self.kc_auth = kc_auth
        self.roles = roles

    def __call__(self, token: str = Depends(oauth2_scheme)) -> dict[str, Any]:
        try:
            decoded: dict[str, Any] = jwt.decode(
                token, self.kc_auth.keycloak_public_key, audience="account", algorithms=["RS256"]
            )
        except JWTError as e:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail=str(e),
                headers={"WWW-Authenticate": "Bearer"},
            )

        username = decoded["preferred_username"]
        roles = decoded["realm_access"]["roles"]

        for required_role in self.roles:
            if required_role not in roles:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED,
                    detail=f"'{username}' doesn't have '{required_role}' role",
                    headers={"WWW-Authenticate": "Bearer"},
                )

        return decoded
