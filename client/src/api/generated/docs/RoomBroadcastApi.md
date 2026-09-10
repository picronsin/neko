# RoomBroadcastApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**broadcastStart**](#broadcaststart) | **POST** /api/room/broadcast/start | Start Broadcast|
|[**broadcastStatus**](#broadcaststatus) | **GET** /api/room/broadcast | Get Broadcast Status|
|[**broadcastStop**](#broadcaststop) | **POST** /api/room/broadcast/stop | Stop Broadcast|

# **broadcastStart**
> broadcastStart(broadcastStatus)

Start broadcasting the room\'s content.

### Example

```typescript
import {
    RoomBroadcastApi,
    Configuration,
    BroadcastStatus
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomBroadcastApi(configuration);

let broadcastStatus: BroadcastStatus; //

const { status, data } = await apiInstance.broadcastStart(
    broadcastStatus
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **broadcastStatus** | **BroadcastStatus**|  | |


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
|**204** | Broadcast started successfully. |  -  |
|**400** | Missing broadcast URL. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | Server is already broadcasting. |  -  |
|**500** | Unable to start broadcast. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **broadcastStatus**
> BroadcastStatus broadcastStatus()

Retrieve the current broadcast status of the room.

### Example

```typescript
import {
    RoomBroadcastApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomBroadcastApi(configuration);

const { status, data } = await apiInstance.broadcastStatus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**BroadcastStatus**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Broadcast status retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **broadcastStop**
> broadcastStop()

Stop broadcasting the room\'s content.

### Example

```typescript
import {
    RoomBroadcastApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomBroadcastApi(configuration);

const { status, data } = await apiInstance.broadcastStop();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Broadcast stopped successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | Server is not broadcasting. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

