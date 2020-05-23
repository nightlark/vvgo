from quart import current_app, g

from httpx import AsyncClient


def get_fetch() -> AsyncClient:
    if 'fetch' not in g:
        g.fetch = AsyncClient(  # type: ignore
            headers={'user-agent': 'vvgo-py/0.0.0'}
        )
    return g.fetch


@current_app.teardown_appcontext
async def teardown_fetch():
    fetch: AsyncClient = g.pop('fetch', None)

    if fetch:
        await fetch.aclose()
