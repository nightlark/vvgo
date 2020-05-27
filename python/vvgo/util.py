from typing import Generator

import httpx
from quart import Response, abort


class BearerAuth(httpx.Auth):
    def __init__(self, token):
        self.auth_header = f'Bearer {token}'

    def auth_flow(
        self, request: httpx.Request
    ) -> Generator[httpx.Request, httpx.Response, None]:
        request.headers['Authorization'] = self.auth_header
        yield request


def abort_with_plain_response(text, code=500):
    return abort(Response(text, code, mimetype='text/plain'))
