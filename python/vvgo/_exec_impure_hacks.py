# isort:skip_file

import trio_asyncio
import asyncio

asyncio.get_running_loop = asyncio.events.get_running_loop
