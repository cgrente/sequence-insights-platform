# Python client

```bash
pip install requests
```

Example:

```py
from sequence_insights_client import SequenceInsightsClient

c = SequenceInsightsClient("http://localhost:8080")
print(c.health())
print(c.ingest([1, -2, 3, 0]))
```
