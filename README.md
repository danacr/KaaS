# KaaS 🧀

### Ephemeral Kubernetes as a Service using Typhoon K8s on Digital Ocean

The purpose of this repository is to provision development clusters in a really simple and fast way. Result of [this](https://twitter.com/errordeveloper/status/1240262848351211520) twitter thread.

It relies on [Simple Typhoon K8s](https://github.com/danacr/simple-typhoon-k8s) to provision clusters using kubernetes batch jobs.

Current configuration:

```yaml
region: nyc3
minutes: 30
```
