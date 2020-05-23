import os

import attr
from quart_trio import QuartTrio as Quart

from .config import Config
from .discord import discord


def create_app(env=None):
    _cfg: Config = Config.from_environ(env or os.environ)

    app = Quart(__name__)
    config = app.config
    config.from_mapping(attr.asdict(_cfg))

    app.register_blueprint(discord)

    scheme = 'https' if config['CERTFILE'] and config['KEYFILE'] else 'http'
    config['URL_PREFIX'] = f'{scheme}://{config["HOST"]}'
    config['DISCORD_REDIRECT_URI'] = f'{config["URL_PREFIX"]}/callback/discord'
    config[
        'PRESESSION_CSRF_STATE_FORMAT'
    ] = f'{config["PRESESSION_KEY_PREFIX"]}:{{}}:csrf-state'
    config[
        'PRESESSION_DISCORD_USER_ID_FORMAT'
    ] = f'{config["PRESESSION_KEY_PREFIX"]}:{{}}:discord-user-id'

    config['DISCORD_SCOPES'] = ['identify']

    return app


if __name__ == '__main__':
    app = create_app()
    app.run(
        use_reloader=True,
        certfile=app.config['CERTFILE'],
        keyfile=app.config['KEYFILE'],
    )
