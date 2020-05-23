from quart import current_app, g

import aioredis
from aioredis import RedisConnection
from trio_asyncio import aio_as_trio


async def init_redis():
    config = current_app.config
    current_app.redis_pool = await aio_as_trio(aioredis.create_redis_pool)(
        config['REDIS_ENDPOINT']
    )


async def get_redis() -> RedisConnection:
    if 'redis' not in g:
        g.redis = await aio_as_trio(current_app.redis_pool.acquire)()  # type: ignore[attr-defined]  # noqa: E501
    return g.redis


@current_app.teardown_appcontext
async def teardown_redis():
    redis: RedisConnection = g.pop('redis', None)

    if redis:
        redis.release()
