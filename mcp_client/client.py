# THIS IS AN MCP CLIENT JUST TO TEST THE TOOLS RESPONSES

import uuid

import asyncio
from fastmcp import Client, FastMCP
from datetime import datetime, timedelta


# HTTP server
#client = Client("https://marketdatamcpserver-production.up.railway.app/mcp")
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
            name="updateUserContext", 
            arguments={
                'user_id': "orestis_user_id",
                "user_profile": {
                    "first_name": "Orestis",
                    "last_name": "Stefanou",
                    "age": "28",
                    "risk_tolerance": "high",
                },
                "user_portfolio":[]
            }
        )
        # result = await client.call_tool(
        #     name="getUserContext", 
        #     arguments={
        #         'user_id': "orestis_user_id",
        #     }               
        # )
        print(result.structured_content)

        

asyncio.run(main())