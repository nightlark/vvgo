# isort:skip_file

from __future__ import annotations

from . import _exec_impure_hacks

from typing import cast

import aioredis
from aioredis import ConnectionsPool, Redis, RedisConnection
from quart import Quart, current_app, g
from trio import MultiError
from trio_asyncio import aio_as_trio


def init_redis(app: Quart):
    async def _init_redis_pool():
        config = app.config

        app.redis_pool: ConnectionsPool = await aio_as_trio(
            aioredis.create_pool
        )(config['REDIS_ENDPOINT'])
        print('_init_redis', app.redis_pool)

    async def _teardown_redis_pool():
        pool: ConnectionsPool = getattr(app, 'redis_pool')
        if pool:
            pool.close()

    app.before_serving(_init_redis_pool)
    app.teardown_appcontext(_teardown_redis_connection)
    app.after_serving(_teardown_redis_pool)


async def get_redis() -> Redis:
    if 'redis' not in g:
        conn: RedisConnection = await aio_as_trio(
            current_app.redis_pool.acquire
        )()
        g.redis = Redis(conn)  # type: ignore[attr-defined]  # noqa: E501
    return g.redis


async def _teardown_redis_connection(_exception=None):
    redis: Redis = g.pop('redis', None)

    if redis:
        current_app.redis_pool.release(redis.connection)
