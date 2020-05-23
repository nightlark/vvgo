from typing import Optional

import environ


@environ.config(prefix='VVGO_AUTH')
class Config:
    HOST: str = 'localhost:5000'
    SECRET_KEY: str = environ.var()

    CERTFILE: Optional[str] = environ.var(None)
    KEYFILE: Optional[str] = environ.var(None)

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

    REDIS_ENDPOINT: str = environ.var('redis://localhost')
