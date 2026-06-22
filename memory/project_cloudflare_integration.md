---
name: project-cloudflare-integration
description: Cloudflare Workers AI added as a third provider alongside Bedrock and Bedrock Mantle, for both chat and image generation
metadata:
  type: project
---

Cloudflare Workers AI was integrated as a third AI provider (`ProviderCloudflare`) for both text and image generation. This replaced deprecated/legacy AWS Bedrock image models with Flux Schnell from Cloudflare.

Models added:
- `CHAT_MODEL_LLAMA` ("llama") → `@cf/meta/llama-3.3-70b-instruct-fp8-fast`
- `IMAGE_MODEL_FLUX` ("flux") → `@cf/black-forest-labs/flux-1-schnell`

**Why:** AWS Bedrock image generation models (Nova Canvas, Titan) are legacy/deprecated. Cloudflare Workers AI provides actively maintained image models (Flux) and additional chat models.

**How to apply:** Cloudflare is optional — activated only when `CLOUDFLARE_ACCOUNT_ID` and `CLOUDFLARE_API_KEY` env vars are set. Without them, `CloudflareAPI` is nil and cloudflare models are excluded from the registry. The `CLOUDFLARE_ACCOUNT_ID` and `CLOUDFLARE_API_KEY` GHA secrets are already wired into `deploy.yaml` as `CloudFlareAccountId` and `CloudFlareAPIKey` SAM parameters.
