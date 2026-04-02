# AWS Bedrock

## Supported models

## Text Models

### Discovery
```shell
aws bedrock list-inference-profiles --region us-east-1 --query 'inferenceProfileSummaries[*].{id:inferenceProfileId,name:inferenceProfileName}' --output json 
```
### Models
```json
[
    {
        "id": "us.anthropic.claude-3-sonnet-20240229-v1:0",
        "name": "US Anthropic Claude 3 Sonnet"
    },
    {
        "id": "us.anthropic.claude-3-opus-20240229-v1:0",
        "name": "US Anthropic Claude 3 Opus"
    },
    {
        "id": "us.anthropic.claude-3-haiku-20240307-v1:0",
        "name": "US Anthropic Claude 3 Haiku"
    },
    {
        "id": "us.meta.llama3-2-11b-instruct-v1:0",
        "name": "US Meta Llama 3.2 11B Instruct"
    },
    {
        "id": "us.meta.llama3-2-3b-instruct-v1:0",
        "name": "US Meta Llama 3.2 3B Instruct"
    },
    {
        "id": "us.meta.llama3-2-90b-instruct-v1:0",
        "name": "US Meta Llama 3.2 90B Instruct"
    },
    {
        "id": "us.meta.llama3-2-1b-instruct-v1:0",
        "name": "US Meta Llama 3.2 1B Instruct"
    },
    {
        "id": "us.anthropic.claude-3-5-haiku-20241022-v1:0",
        "name": "US Anthropic Claude 3.5 Haiku"
    },
    {
        "id": "us.meta.llama3-1-8b-instruct-v1:0",
        "name": "US Meta Llama 3.1 8B Instruct"
    },
    {
        "id": "us.meta.llama3-1-70b-instruct-v1:0",
        "name": "US Meta Llama 3.1 70B Instruct"
    },
    {
        "id": "us.meta.llama3-3-70b-instruct-v1:0",
        "name": "US Meta Llama 3.3 70B Instruct"
    },
    {
        "id": "us.deepseek.r1-v1:0",
        "name": "US DeepSeek-R1"
    },
    {
        "id": "us.mistral.pixtral-large-2502-v1:0",
        "name": "US Mistral Pixtral Large 25.02"
    },
    {
        "id": "us.meta.llama4-scout-17b-instruct-v1:0",
        "name": "US Llama 4 Scout 17B Instruct"
    },
    {
        "id": "us.meta.llama4-maverick-17b-instruct-v1:0",
        "name": "US Llama 4 Maverick 17B Instruct"
    },
    {
        "id": "us.amazon.nova-premier-v1:0",
        "name": "US Nova Premier"
    },
    {
        "id": "us.anthropic.claude-opus-4-20250514-v1:0",
        "name": "US Claude Opus 4"
    },
    {
        "id": "us.twelvelabs.pegasus-1-2-v1:0",
        "name": "US Pegasus v1.2"
    },
    {
        "id": "global.anthropic.claude-sonnet-4-20250514-v1:0",
        "name": "Global Claude Sonnet 4"
    },
    {
        "id": "us.stability.stable-conservative-upscale-v1:0",
        "name": "US Stable Image Conservative Upscale"
    },
    {
        "id": "us.twelvelabs.marengo-embed-2-7-v1:0",
        "name": "US TwelveLabs Marengo Embed v2.7"
    },
    {
        "id": "us.stability.stable-creative-upscale-v1:0",
        "name": "US Stable Image Creative Upscale"
    },
    {
        "id": "us.stability.stable-fast-upscale-v1:0",
        "name": "US Stable Image Fast Upscale"
    },
    {
        "id": "us.stability.stable-image-control-sketch-v1:0",
        "name": "US Stable Image Control Sketch"
    },
    {
        "id": "us.stability.stable-image-control-structure-v1:0",
        "name": "US Stable Image Control Structure"
    },
    {
        "id": "us.stability.stable-image-erase-object-v1:0",
        "name": "US Stable Image Erase Object"
    },
    {
        "id": "us.stability.stable-image-inpaint-v1:0",
        "name": "US Stable Image Inpaint"
    },
    {
        "id": "us.stability.stable-image-remove-background-v1:0",
        "name": "US Stable Image Remove Background"
    },
    {
        "id": "us.stability.stable-image-search-replace-v1:0",
        "name": "US Stable Image Search and Replace"
    },
    {
        "id": "us.stability.stable-image-search-recolor-v1:0",
        "name": "US Stable Image Search and Recolor"
    },
    {
        "id": "us.stability.stable-image-style-guide-v1:0",
        "name": "US Stable Image Style Guide"
    },
    {
        "id": "us.stability.stable-outpaint-v1:0",
        "name": "US Stable Image Outpaint"
    },
    {
        "id": "us.stability.stable-style-transfer-v1:0",
        "name": "US Stable Image Style Transfer"
    },
    {
        "id": "us.anthropic.claude-3-7-sonnet-20250219-v1:0",
        "name": "US Anthropic Claude 3.7 Sonnet"
    },
    {
        "id": "us.cohere.embed-v4:0",
        "name": "US Cohere Embed v4"
    },
    {
        "id": "global.cohere.embed-v4:0",
        "name": "Global Cohere Embed v4"
    },
    {
        "id": "us.twelvelabs.marengo-embed-3-0-v1:0",
        "name": "US TwelveLabs Marengo Embed 3.0"
    },
    {
        "id": "global.twelvelabs.pegasus-1-2-v1:0",
        "name": "GLOBAL TwelveLabs Pegasus v1.2"
    },
    {
        "id": "us.amazon.nova-micro-v1:0",
        "name": "US Nova Micro"
    },
    {
        "id": "us.writer.palmyra-x4-v1:0",
        "name": "US Writer Palmyra X4"
    },
    {
        "id": "us.writer.palmyra-x5-v1:0",
        "name": "US Writer Palmyra X5"
    },
    {
        "id": "global.amazon.nova-2-lite-v1:0",
        "name": "GLOBAL Amazon Nova 2 Lite"
    },
    {
        "id": "us.anthropic.claude-opus-4-1-20250805-v1:0",
        "name": "US Claude Opus 4.1"
    },
    {
        "id": "us.anthropic.claude-sonnet-4-6",
        "name": "US Anthropic Claude Sonnet 4.6"
    },
    {
        "id": "us.anthropic.claude-sonnet-4-5-20250929-v1:0",
        "name": "US Anthropic Claude Sonnet 4.5"
    },
    {
        "id": "us.anthropic.claude-haiku-4-5-20251001-v1:0",
        "name": "US Anthropic Claude Haiku 4.5"
    },
    {
        "id": "us.anthropic.claude-opus-4-6-v1",
        "name": "US Anthropic Claude Opus 4.6"
    },
    {
        "id": "us.anthropic.claude-opus-4-5-20251101-v1:0",
        "name": "US Anthropic Claude Opus 4.5"
    },
    {
        "id": "us.amazon.nova-lite-v1:0",
        "name": "US Nova Lite"
    },
    {
        "id": "us.amazon.nova-pro-v1:0",
        "name": "US Nova Pro"
    },
    {
        "id": "us.amazon.nova-2-lite-v1:0",
        "name": "US Amazon Nova 2 Lite"
    },
    {
        "id": "global.anthropic.claude-haiku-4-5-20251001-v1:0",
        "name": "Global Anthropic Claude Haiku 4.5"
    },
    {
        "id": "us.anthropic.claude-sonnet-4-20250514-v1:0",
        "name": "US Claude Sonnet 4"
    },
    {
        "id": "global.anthropic.claude-sonnet-4-5-20250929-v1:0",
        "name": "Global Claude Sonnet 4.5"
    },
    {
        "id": "global.anthropic.claude-opus-4-5-20251101-v1:0",
        "name": "GLOBAL Anthropic Claude Opus 4.5"
    },
    {
        "id": "global.anthropic.claude-opus-4-6-v1",
        "name": "Global Anthropic Claude Opus 4.6"
    },
    {
        "id": "global.anthropic.claude-sonnet-4-6",
        "name": "Global Anthropic Claude Sonnet 4.6"
    }
]
```

## Image models

### Discovery
```shell
aws bedrock list-foundation-models --region us-east-1 --by-output-modality IMAGE --query 'modelSummaries[*].{id:modelId,name:modelName}' --output json | pbcopy
```
### Models
```json
[
    {
        "id": "stability.stable-creative-upscale-v1:0",
        "name": "Stable Image Creative Upscale"
    },
    {
        "id": "stability.stable-image-remove-background-v1:0",
        "name": "Stable Image Remove Background"
    },
    {
        "id": "stability.stable-image-control-sketch-v1:0",
        "name": "Stable Image Control Sketch"
    },
    {
        "id": "stability.stable-conservative-upscale-v1:0",
        "name": "Stable Image Conservative Upscale"
    },
    {
        "id": "stability.stable-image-search-recolor-v1:0",
        "name": "Stable Image Search and Recolor"
    },
    {
        "id": "stability.stable-fast-upscale-v1:0",
        "name": "Stable Image Fast Upscale"
    },
    {
        "id": "stability.stable-image-erase-object-v1:0",
        "name": "Stable Image Erase Object"
    },
    {
        "id": "stability.stable-image-control-structure-v1:0",
        "name": "Stable Image Control Structure"
    },
    {
        "id": "stability.stable-outpaint-v1:0",
        "name": "Stable Image Outpaint"
    },
    {
        "id": "stability.stable-image-inpaint-v1:0",
        "name": "Stable Image Inpaint"
    },
    {
        "id": "stability.stable-image-style-guide-v1:0",
        "name": "Stable Image Style Guide"
    },
    {
        "id": "stability.stable-style-transfer-v1:0",
        "name": "Stable Image Style Transfer"
    },
    {
        "id": "stability.stable-image-search-replace-v1:0",
        "name": "Stable Image Search and Replace"
    },
    {
        "id": "amazon.titan-image-generator-v2:0",
        "name": "Titan Image Generator G1 v2"
    },
    {
        "id": "amazon.nova-canvas-v1:0",
        "name": "Nova Canvas"
    }
]
```

## Video models

### Discovery
```shell
aws bedrock list-foundation-models --region us-east-1 --by-output-modality VIDEO --query 'modelSummaries[*].{id:modelId,name:modelName}' --output json | pbcopy
```

### Models
```json
[
    {
        "id": "amazon.nova-reel-v1:0",
        "name": "Nova Reel"
    },
    {
        "id": "amazon.nova-reel-v1:1",
        "name": "Nova Reel"
    }
]
```

## Embedding models

### Discovery
```shell
aws bedrock list-foundation-models --region us-east-1 --by-output-modality EMBEDDING --query 'modelSummaries[*].{id:modelId,name:modelName}' --output json | pbcopy
```

### Models

```json
[
    {
        "id": "amazon.nova-2-multimodal-embeddings-v1:0",
        "name": "Amazon Nova Multimodal Embeddings"
    },
    {
        "id": "twelvelabs.marengo-embed-3-0-v1:0",
        "name": "Marengo Embed 3.0"
    },
    {
        "id": "twelvelabs.marengo-embed-2-7-v1:0",
        "name": "Marengo Embed v2.7"
    },
    {
        "id": "cohere.embed-v4:0",
        "name": "Embed v4"
    },
    {
        "id": "amazon.titan-embed-g1-text-02",
        "name": "Titan Text Embeddings v2"
    },
    {
        "id": "amazon.titan-embed-text-v1:2:8k",
        "name": "Titan Embeddings G1 - Text"
    },
    {
        "id": "amazon.titan-embed-text-v1",
        "name": "Titan Embeddings G1 - Text"
    },
    {
        "id": "amazon.titan-embed-text-v2:0:8k",
        "name": "Titan Text Embeddings V2"
    },
    {
        "id": "amazon.titan-embed-text-v2:0",
        "name": "Titan Text Embeddings V2"
    },
    {
        "id": "amazon.titan-embed-image-v1:0",
        "name": "Titan Multimodal Embeddings G1"
    },
    {
        "id": "amazon.titan-embed-image-v1",
        "name": "Titan Multimodal Embeddings G1"
    },
    {
        "id": "cohere.embed-english-v3:0:512",
        "name": "Embed English"
    },
    {
        "id": "cohere.embed-english-v3",
        "name": "Embed English"
    },
    {
        "id": "cohere.embed-multilingual-v3:0:512",
        "name": "Embed Multilingual"
    },
    {
        "id": "cohere.embed-multilingual-v3",
        "name": "Embed Multilingual"
    }
]
```

## Speech models

### Discovery
```shell
aws bedrock list-foundation-models --region us-east-1 --by-output-modality SPEECH --query 'modelSummaries[*].{id:modelId,name:modelName}' --output json | pbcopy
```

### Models

```json
[
    {
        "id": "amazon.nova-2-sonic-v1:0",
        "name": "Nova 2 Sonic"
    },
    {
        "id": "amazon.nova-sonic-v1:0",
        "name": "Nova Sonic"
    }
]
```