import asyncio
import logging
import time
from typing import Any
from urllib.parse import urljoin

import aiojobs
import httpx
from aiojobs import Scheduler
from config import keycloak_base_url
from httpx import AsyncClient
from httpx import Request
from pydantic import BaseModel
from pydantic import Field


class OIDCData(BaseModel):
    access_token: str
    expires_in: int
    refresh_expires_in: int
    token_type: str
    not_before_policy: int = Field(alias="not-before-policy")
    scope: str


class KeycloakAuth(httpx.Auth):  # type: ignore
    scheduler: Scheduler
    oidc_data: OIDCData

    def __init__(self, client_id: str, client_secret: str, minimum_token_rest_of_life: int) -> None:
        self.token_url = urljoin(keycloak_base_url, "/auth/realms/myrealm/protocol/openid-connect/token")

        self.client_id = client_id
        self.client_secret = client_secret
        self.refresh_seconds_before_expire = minimum_token_rest_of_life

        self._autoupdate_started = False
        self._refresh_until: float = time.time() - 1  # Unix timestamp

    async def start_auto_update(self) -> None:
        if self._autoupdate_started:
            logging.error("Auto update already started")
            return

        self._autoupdate_started = True

        self.scheduler = await aiojobs.create_scheduler()

        async def task() -> None:
            while True:
                await self._refresh_token_if_needed(time.time())
                await asyncio.sleep(1)

        await self.scheduler.spawn(task())

    async def _refresh_token_if_needed(self, now: float) -> None:
        if now < self._refresh_until:
            return

        response = httpx.post(
            url=self.token_url,
            data={"grant_type": "client_credentials"},
            auth=(self.client_id, self.client_secret),
        )

        if response.status_code != 200:
            msg = f"Failed to refresh token: {response.status_code}, response: {response.text}"
            raise Exception(msg)

        self.oidc_data = OIDCData(**response.json())
        self._refresh_until = time.time() + self.oidc_data.expires_in - self.refresh_seconds_before_expire

        logging.error("Token refreshed")

    async def async_auth_flow(self, request: Request) -> Request:
        request.headers["Authorization"] = f"BEARER {self.oidc_data.access_token}"
        yield request


class ServiceClient:
    def __init__(self, async_client: AsyncClient) -> None:
        self.client = async_client

    async def call_server(self, url: str) -> dict[str, Any]:
        response = await self.client.get(url)
        return dict(status_code=response.status_code, response=response.text)


if __name__ == "__main__":
    keycloak_auth = KeycloakAuth("my-client", "my-secret", 100)
    async_client = AsyncClient(auth=keycloak_auth)
    client = ServiceClient(async_client)
