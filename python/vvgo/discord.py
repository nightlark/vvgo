import secrets

import httpx
from aioredis import RedisConnection
from oauthlib.oauth2 import WebApplicationClient
from quart import Blueprint, current_app, g, redirect, request
from trio_asyncio import aio_as_trio

from .fetch import get_fetch
from .redis import get_redis
from .util import BearerAuth, abort_with_plain_response

discord = Blueprint('discord', __name__)


def get_discord_oauth_client() -> WebApplicationClient:
    config = current_app.config
    if 'discord_oauth_client' not in g:
        g.discord_oauth_client = WebApplicationClient(  # type: ignore
            client_id=config['DISCORD_CLIENT_ID'],
            scope=config['DISCORD_SCOPES'],
        )
    return g.discord_oauth_client


@discord.route('/login/discord')
async def login():
    config = current_app.config
    client = get_discord_oauth_client()
    redis: RedisConnection = await get_redis()

    pre_session = secrets.token_urlsafe(256)
    state = secrets.token_urlsafe(256)

    await aio_as_trio(redis.set)(
        f'{config["PRESESSION_KEY_PREFIX"]}:{pre_session}:csrf-state',
        state,
        expire=config['PRESESSION_EXPIRY'],
    )

    uri, _headers, _body = client.prepare_authorization_request(
        config['DISCORD_AUTHORIZATION_ENDPOINT'],
        state=state,
        redirect_url=config['DISCORD_REDIRECT_URI'],
    )

    resp = redirect(uri)
    resp.set_cookie(
        config['PRESESSION_COOKIE'],
        pre_session,
        max_age=config['PRESESSION_EXPIRY'],
    )
    return resp


@discord.route('/callback/discord')
async def callback_discord():
    # TODO: Refactor into multiple functions
    config = current_app.config
    if (
        config['PRESESSION_COOKIE'] not in request.cookies
        or 'state' not in request.args
        or 'code' not in request.args
    ):
        return abort_with_plain_response('Invalid callback', 400)

    client = get_discord_oauth_client()
    fetch = get_fetch()
    redis = await get_redis()
    pre_session = request.cookies[config['PRESESSION_COOKIE']]

    state = await aio_as_trio(redis.get)(
        config['PRESESSION_CSRF_STATE_FORMAT'].format(pre_session)
    )
    if request.args['state'] != state:
        return abort_with_plain_response('Invalid callback', 400)

    code = request.args['code']
    body = client.prepare_request_body(
        code=code,
        redirect_uri=config['DISCORD_REDIRECT_URI'],
        client_secret=config['DISCORD_CLIENT_SECRET'],
    )

    try:
        resp = await fetch.post(config['DISCORD_TOKEN_ENDPOINT'], data=body)
        resp.raise_for_status()
        resp_json = resp.json()
        token = resp['access_token']
    except (httpx.HTTPError, ValueError, KeyError):
        return abort_with_plain_response('Could not retrieve token', 500)

    try:
        resp = await fetch.get(
            config['DISCORD_USER_INFO_ENDPOINT'], auth=BearerAuth(token)
        )
        resp.raise_for_status()
        resp_json = resp.json()
        user_id = resp_json['id']
    except (httpx.HTTPError, ValueError, KeyError):
        return abort_with_plain_response('Could not retrieve user id', 500)

    pipe = redis.pipeline()
    pipe.delete(config['PRESESSION_CSRF_STATE_FORMAT'].format(pre_session))
    pipe.set(
        config['PRESESSION_DISCORD_USER_ID_FORMAT'].format(pre_session),
        user_id,
        expire=config['PRESESSION_EXPIRY'],
    )
    await aio_as_trio(pipe.execute)()

    return redirect('/')
