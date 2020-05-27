from httpx import AsyncClient
from quart import Quart, g


def init_fetch(app: Quart):
    app.teardown_appcontext(_teardown_fetch)


def get_fetch() -> AsyncClient:
    if 'fetch' not in g:
        g.fetch = AsyncClient(  # type: ignore
            headers={'user-agent': 'vvgo-py/0.0.0'}
        )
    return g.fetch


async def _teardown_fetch(_exception=None):
    fetch: AsyncClient = g.pop('fetch', None)

    if fetch:
        await fetch.aclose()
