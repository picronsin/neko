# MemberProfile


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **string** | The name of the member. | [optional] [default to undefined]
**avatar** | **string** | Base64 image data URL for the member avatar. | [optional] [default to undefined]
**is_admin** | **boolean** | Indicates if the member is an admin. | [optional] [default to undefined]
**can_login** | **boolean** | Indicates if the member can log in. | [optional] [default to undefined]
**can_connect** | **boolean** | Indicates if the member can connect. | [optional] [default to undefined]
**can_watch** | **boolean** | Indicates if the member can watch. | [optional] [default to undefined]
**can_host** | **boolean** | Indicates if the member can host. | [optional] [default to undefined]
**can_share_media** | **boolean** | Indicates if the member can share media. | [optional] [default to undefined]
**can_access_clipboard** | **boolean** | Indicates if the member can access the clipboard. | [optional] [default to undefined]
**sends_inactive_cursor** | **boolean** | Indicates if the member sends inactive cursor. | [optional] [default to undefined]
**can_see_inactive_cursors** | **boolean** | Indicates if the member can see inactive cursors. | [optional] [default to undefined]
**plugins** | **{ [key: string]: any; }** | Additional plugin settings. | [optional] [default to undefined]

## Example

```typescript
import { MemberProfile } from './api';

const instance: MemberProfile = {
    name,
    avatar,
    is_admin,
    can_login,
    can_connect,
    can_watch,
    can_host,
    can_share_media,
    can_access_clipboard,
    sends_inactive_cursor,
    can_see_inactive_cursors,
    plugins,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
