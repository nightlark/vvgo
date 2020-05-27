import os
from typing import Mapping, Optional, Type

import environ


@environ.config(prefix='VVGO_AUTH')
class Config:
    HOST: str = environ.var('localhost')
    PORT: int = environ.var(5000)
    SECRET_KEY: str = environ.var()

    TLS_CERTFILE: Optional[str] = environ.var(None)
    TLS_KEYFILE: Optional[str] = environ.var(None)

    DISCORD_CLIENT_ID: str = environ.var()
    DISCORD_CLIENT_SECRET: str = environ.var()

    DISCORD_USER_INFO_ENDPOINT: str = environ.var(
        'https://discord.com/api/users/@me'
    )
    DISCORD_AUTHORIZATION_ENDPOINT: str = environ.var(
        'https://discord.com/api/oauth2/authorize'
    )
    DISCORD_TOKEN_ENDPOINT: str = environ.var(
        'https://discord.com/api/oauth2/token'
    )

    PRESESSION_COOKIE: str = environ.var('vvgo-presession')
    PRESESSION_KEY_PREFIX: str = environ.var('vvgo-presession')
    PRESESSION_EXPIRY: int = environ.var(300)

    REDIS_ENDPOINT: str = environ.var('redis://localhost:6379')


Config: Type[Config] = Config


def load_config(env: Optional[Mapping[str, str]] = None) -> Config:
    return Config.from_environ(env or os.environ)  # type: ignore[attr-defined]
