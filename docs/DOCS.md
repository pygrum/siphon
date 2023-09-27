# Configuration

Here is a full example configuration for Siphon:

```yaml
RefreshRate: 5 # Refresh sample list every 5 minutes
Sources:
- name: VirusTotal
  endpoint:
  APIKey:
- name: MalwareBazaar
  endpoint:
  APIKey:
- name: HybridAnalysis
  endpoint:
  APIKey:
```

Each entry in `source` correlates to a supported threat intelligence feed integration.
All that is needed is an endpoint and API key.