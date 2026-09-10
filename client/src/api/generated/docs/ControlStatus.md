# ControlStatus


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**has_host** | **boolean** | Indicates if there is a host currently. | [default to undefined]
**host_id** | **string** | The ID of the current host, if any. | [optional] [default to undefined]
**epoch** | **number** | Monotonically increasing control lease epoch. | [default to undefined]

## Example

```typescript
import { ControlStatus } from './api';

const instance: ControlStatus = {
    has_host,
    host_id,
    epoch,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
