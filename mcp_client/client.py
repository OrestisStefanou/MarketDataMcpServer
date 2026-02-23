# THIS IS AN MCP CLIENT JUST TO TEST THE TOOLS RESPONSES

import uuid

import asyncio
from fastmcp import Client, FastMCP
from datetime import datetime, timedelta


# HTTP server
client = Client("http://localhost:8080/mcp")


async def main():
    async with client:
        # Basic server interaction
        await client.ping()

        # Current datetime
        now = datetime.now()

        # Datetime 5 days ago
        five_days_ago = now - timedelta(days=5)
        
        result = await client.call_tool(
            name="getInvestingIdeaStocks", 
            arguments={
                'idea_id': "5d87f6d2-6fa8-4686-a043-c5ef0c5985ef",
            }               
        )
        print(result.structured_content)

        

asyncio.run(main())