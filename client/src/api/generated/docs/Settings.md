# Settings


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**private_mode** | **boolean** | Indicates if the room is in private mode. | [optional] [default to undefined]
**locked_controls** | **boolean** | Indicates if the room controls are locked. | [optional] [default to undefined]
**implicit_hosting** | **boolean** | Indicates if implicit hosting is enabled. | [optional] [default to undefined]
**inactive_cursors** | **boolean** | Indicates if inactive cursors are shown. | [optional] [default to undefined]
**merciful_reconnect** | **boolean** | Indicates if merciful reconnect is enabled. | [optional] [default to undefined]
**control_lease_ttl** | **number** | Control ownership idle timeout in seconds. | [optional] [default to undefined]
**plugins** | **{ [key: string]: any; }** | Additional plugin settings. | [optional] [default to undefined]

## Example

```typescript
import { Settings } from './api';

const instance: Settings = {
    private_mode,
    locked_controls,
    implicit_hosting,
    inactive_cursors,
    merciful_reconnect,
    control_lease_ttl,
    plugins,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
