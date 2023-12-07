# RAG with Amazon Bedrock and Bedrock/Claude as API

## Architecure

## Configure

- Create kendra index
- Set KENDRA_LANGUAGE_CODE in taskfile.yml
- Set KENDRA_REGION in taskfile.yml

```bash
task parameter
```

## Install

```bash
task build
task deploy
```

Create and attach api keys in the console.

### Usage plan
- Create
- Add stage

### API Key
- Create API Key
- Add to usage plan
     Attach key

### API

Redeploy

## Set Kendra configuration
The values are stored in the Systems Manager Parameter store

Content | SSM Parameter
---|---
KENDRA_ID | /rag/KENDRA_INDEX_ID
KENDRA_REGION | /rag/KENDRA_REGION
KENDRA_LANGUAGE_CODE | /rag/KENDRA_LANGUAGE_CODE

The parameters has to be in the same region as the API Gateway.
