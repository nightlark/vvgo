# isort:skip_file

from __future__ import annotations

from vvgo import _exec_impure_hacks

import asyncio
import subprocess
from pathlib import Path
from typing import AsyncGenerator, Callable, Coroutine, TypeVar

import attr
from pytest import fixture
from quart_trio import QuartTrio
from trio_asyncio import open_loop

from vvgo.config import Config, load_config


@fixture(scope='session')
def app_config() -> Config:
    return load_config()


@fixture
def base_app(app_config: Config) -> QuartTrio:
    """
    Create an app with config, but no extensions.
    """
    app = QuartTrio(__name__)
    config = app.config
    config.from_mapping(attr.asdict(app_config))
    return app


@fixture
async def asyncio_loop():
    async with open_loop() as loop:
        asyncio.set_event_loop(loop)
        yield loop
