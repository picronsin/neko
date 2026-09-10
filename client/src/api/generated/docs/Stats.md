# Stats


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**has_host** | **boolean** | Indicates if there is a host currently. | [optional] [default to undefined]
**host_id** | **string** | The ID of the current host, if any. | [optional] [default to undefined]
**epoch** | **number** | Monotonically increasing control lease epoch used to reject stale input. | [optional] [default to undefined]
**server_started_at** | **string** | The timestamp when the server started. | [optional] [default to undefined]
**total_users** | **number** | The total number of users connected. | [optional] [default to undefined]
**last_user_left_at** | **string** | The timestamp when the last user left, if any. | [optional] [default to undefined]
**total_admins** | **number** | The total number of admins connected. | [optional] [default to undefined]
**last_admin_left_at** | **string** | The timestamp when the last admin left, if any. | [optional] [default to undefined]

## Example

```typescript
import { Stats } from './api';

const instance: Stats = {
    has_host,
    host_id,
    epoch,
    server_started_at,
    total_users,
    last_user_left_at,
    total_admins,
    last_admin_left_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
