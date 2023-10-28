import asyncio
from maelstrom import Node, Body, Request

node = Node()
counter = 0
cond = asyncio.Condition()


async def increase_counter(value: int):
    counter += value
    async with cond:
        cond.notify_all()


@node.handler()
async def add(req: Request):
    delta = req.body["delta"]
    return {"type": "add_ok"}


@node.handler
async def read(req: Request) -> Body:
    return {"type": "read_ok", "value": counter}


async def gossip_task(neighbor: str) -> None:
    sent = set()
    while True:
        # assert sent <= messages
        if len(sent) == len(messages):
            async with cond:
                # Wait for the next update to our message set.
                await cond.wait_for(lambda: len(sent) != len(messages))

        to_send = messages - sent
        body = {"type": "broadcast_many", "messages": list(to_send)}
        resp = await node.rpc(neighbor, body)
        if resp["type"] == "broadcast_many_ok":
            sent.update(to_send)
