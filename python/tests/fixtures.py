# isort:skip_file

from __future__ import annotations

from vvgo import _exec_impure_hacks

import asyncio
import subprocess
from pathlib import Path
from typing import AsyncGenerator, Callable, Coroutine, TypeVar

from pytest import fixture
from quart_trio import QuartTrio
from trio_asyncio import open_loop

from vvgo.config import Config, load_config


@fixture(scope='session')
def app_config() -> Config:
    return load_config()


@fixture
async def quart_trio_app() -> AsyncGenerator[QuartTrio, None]:
    app = QuartTrio(__name__)
    await app.startup()
    yield app
    await app.shutdown()


@fixture
async def asyncio_loop():
    async with open_loop() as loop:
        asyncio.set_event_loop(loop)
        yield loop
