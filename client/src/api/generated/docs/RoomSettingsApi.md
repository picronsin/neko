# RoomSettingsApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**settingsGet**](#settingsget) | **GET** /api/room/settings | Get Room Settings|
|[**settingsSet**](#settingsset) | **POST** /api/room/settings | Update Room Settings|

# **settingsGet**
> Settings settingsGet()

Retrieve the current settings of the room.

### Example

```typescript
import {
    RoomSettingsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomSettingsApi(configuration);

const { status, data } = await apiInstance.settingsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Settings**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Room settings retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsSet**
> settingsSet(settings)

Update the settings of the room.

### Example

```typescript
import {
    RoomSettingsApi,
    Configuration,
    Settings
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomSettingsApi(configuration);

let settings: Settings; //

const { status, data } = await apiInstance.settingsSet(
    settings
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **settings** | **Settings**|  | |


### Return type

void (empty response body)

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Room settings updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

