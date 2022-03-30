import asyncio
import logging
import time

import pytest
from client import KeycloakAuth

logger = logging.getLogger(__name__)
logger.setLevel("DEBUG")

client_id = "backend"
client_secret = "00000000–0000–0000–0000–000000000000"


@pytest.mark.asyncio  # type: ignore
async def test_refresh_token() -> None:
    client = KeycloakAuth(client_id, client_secret, 295)
    await client._refresh_token_if_needed(time.time())

    assert client.oidc_data.access_token


@pytest.mark.asyncio  # type: ignore
async def test_token_autoupdate_cicle() -> None:
    client = KeycloakAuth(client_id, client_secret, 299)
    await client.start_auto_update()

    await asyncio.sleep(5)

    assert client.oidc_data.access_token

    await client.scheduler.close()
