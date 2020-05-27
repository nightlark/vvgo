import os
from typing import Mapping, Union

import attr
from quart_trio import QuartTrio as Quart

from .config import Config, load_config
from .discord import discord
from .fetch import init_fetch
from .redis import init_redis


def create_app(env: Union[None, Config, Mapping[str, str]] = None):
    _cfg: Config
    if isinstance(env, Config):
        _cfg = env
    else:
        _cfg = load_config(env or os.environ)

    app = Quart(__name__)
    config = app.config
    config.from_mapping(attr.asdict(_cfg))

    app.register_blueprint(discord)
    init_fetch(app)
    init_redis(app)

    scheme = (
        'https' if config['TLS_CERTFILE'] and config['TLS_KEYFILE'] else 'http'
    )
    host = config['HOST']
    port = config['PORT']
    url_prefix = f'{scheme}://{host}:{port}'
    presession_key_prefix = config['PRESESSION_KEY_PREFIX']

    config['SCHEME'] = scheme
    config['URL_PREFIX'] = url_prefix

    config['DISCORD_REDIRECT_URI'] = f'{url_prefix}/callback/discord'
    config['DISCORD_SCOPES'] = ['identify']

    config[
        'PRESESSION_CSRF_STATE_FORMAT'
    ] = f'{presession_key_prefix}:{{}}:csrf-state'
    config[
        'PRESESSION_DISCORD_USER_ID_FORMAT'
    ] = f'{presession_key_prefix}:{{}}:discord-user-id'

    return app
