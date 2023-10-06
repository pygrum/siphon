# Configuration

## Siphon
Here is a full example configuration for Siphon:

```yaml
RefreshRate: 5 # Refresh sample list every 5 minutes
cert_file: "/path/to/cert/file.crt"
key_file: "/path/to/key/file.crt"
Sources:
- name: "VirusTotal"
  endpoint:
  APIKey:
- name: "MalwareBazaar"
  endpoint:
  APIKey:
- name: "HybridAnalysis"
  endpoint:
  APIKey:
```

Each entry in `source` correlates to a supported threat intelligence feed integration.
All that is needed is an endpoint and API key.

## Siphon agent
The siphon agent monitors folders located in a honeypot. 
Here is an example configuration for the Siphon agent:

```yaml
cache: true
cert_file: "/path/to/certificate/file"
key_file: "/path/to/key/file"
monitor_folders:
  - path: "/path/to/folder_1"
    recursive: true
  - path: "/path/to/folder_2"
    recursive: false
```

If `cache` is set to true, then files written to disk in the monitored folders
will be saved to a protected folder, in case they get deleted or moved later on.

The `recursive` monitoring option means that that specific folder and all subfolders will
be monitored for file changes, otherwise, only the top level folder will be monitored.

## Agent generator
Here is an example agent generator file. These fields are used to build a new agent.
```yaml
cert_file: "/path/to/siphon/client/certificate"
src_path: "/path/to/siphon/source/folder"
name: "agent"
os: "windows"
arch: "amd64"
host: "0.0.0.0"
port: "8080"
outfile: "agent.exe"
```

### Agent build steps

After generating an agent with `siphon_gen`, run the `setup_agent.sh` script (under /scripts) to
create a folder with the agent, tls key pair and default configuration file, which you should tweak as you wish.

You'll need to update the configuration file with:
1. The path to your key pair when _it's on the target machine_ (this is set to `/root/agent.<key|crt>` by default)
2. Folders to monitor on the host, following the same format as shown under 'Siphon agent' above