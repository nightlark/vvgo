# isort:skip_file

from __future__ import annotations

from vvgo import _exec_impure_hacks

import asyncio

import attr
from aioredis import ConnectionsPool, Redis
from pytest import fixture
from quart import Quart, current_app, g
from trio_asyncio import aio_as_trio

from vvgo.config import Config
from vvgo.redis import get_redis, init_redis

from .fixtures import app_config, asyncio_loop, quart_trio_app


@fixture
async def redis_app(quart_trio_app, app_config, asyncio_loop):
    app: Quart = quart_trio_app
    config = app.config
    config.from_mapping(attr.asdict(app_config))

    init_redis(app)

    await app.startup()

    yield app

    await app.shutdown()


async def test_redis(redis_app: Quart, asyncio_loop):
    print(asyncio_loop)
    async with redis_app.app_context():
        pool: ConnectionsPool = current_app.redis_pool
        redis: Redis = await get_redis()
        assert redis

        result = await aio_as_trio(redis.echo)(b'Hello, world!')

        assert result == b'Hello, world!'
