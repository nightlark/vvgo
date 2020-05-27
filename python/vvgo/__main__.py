import trio_asyncio

from .core import create_app

app = create_app()
config = app.config
url_prefix = config['URL_PREFIX']
print(f'Running on {url_prefix} (CTRL-C to quit)')
task = app.run_task(
    host=config['HOST'],
    port=config['PORT'],
    debug=True,
    use_reloader=True,
    certfile=config['TLS_CERTFILE'],
    keyfile=config['TLS_KEYFILE'],
)
trio_asyncio.run(task)
