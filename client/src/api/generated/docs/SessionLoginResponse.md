# SessionLoginResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The unique identifier of the session. | [optional] [default to undefined]
**profile** | [**MemberProfile**](MemberProfile.md) |  | [optional] [default to undefined]
**state** | [**SessionState**](SessionState.md) |  | [optional] [default to undefined]
**token** | **string** | The session token, only if cookie authentication is disabled. | [optional] [default to undefined]

## Example

```typescript
import { SessionLoginResponse } from './api';

const instance: SessionLoginResponse = {
    id,
    profile,
    state,
    token,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
