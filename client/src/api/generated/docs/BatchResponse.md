# BatchResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**path** | **string** | The API endpoint path. | [optional] [default to undefined]
**method** | **string** | The HTTP method used. | [optional] [default to undefined]
**body** | **any** | The response body from the API call. | [optional] [default to undefined]
**status** | **number** | The HTTP status code of the response. | [optional] [default to undefined]

## Example

```typescript
import { BatchResponse } from './api';

const instance: BatchResponse = {
    path,
    method,
    body,
    status,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
