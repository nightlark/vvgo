# isort:skip_file pylint: disable=redefined-outer-name, unused-argument

from __future__ import annotations

from vvgo import (  # pylint: disable=unused-import, ungrouped-imports
    _exec_impure_hacks,
)

import attr
from httpx import AsyncClient, Response
from pytest import fixture
from quart import Quart, g
from quart_trio import QuartTrio
from trio_asyncio import aio_as_trio

from vvgo.config import Config
from vvgo.fetch import init_fetch, get_fetch


@fixture
async def fetch_app(base_app: QuartTrio, asyncio_loop):
    app = base_app

    init_fetch(app)

    await app.startup()
    yield app
    await app.shutdown()


async def test_fetch(fetch_app: Quart, asyncio_loop):
    async with fetch_app.app_context():
        client: AsyncClient = get_fetch()
        assert client

        result: Response

        result = await client.get(
            'https://connectivitycheck.gstatic.com/generate_204'
        )

        assert result.status_code == 204

        result = await client.get('http://detectportal.firefox.com')

        assert result.status_code == 200
        assert result.text.strip() == 'success'
