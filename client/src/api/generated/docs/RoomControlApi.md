# RoomControlApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**controlGive**](#controlgive) | **POST** /api/room/control/give/{sessionId} | Give Control|
|[**controlRelease**](#controlrelease) | **POST** /api/room/control/release | Release Control|
|[**controlRequest**](#controlrequest) | **POST** /api/room/control/request | Request Control|
|[**controlReset**](#controlreset) | **POST** /api/room/control/reset | Reset Control|
|[**controlStatus**](#controlstatus) | **GET** /api/room/control | Get Control Status|
|[**controlTake**](#controltake) | **POST** /api/room/control/take | Take Control|

# **controlGive**
> controlGive()

Give control of the room to a specific session.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

let sessionId: string; //The identifier of the session. (default to undefined)

const { status, data } = await apiInstance.controlGive(
    sessionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sessionId** | [**string**] | The identifier of the session. | defaults to undefined|


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
|**204** | Control given successfully. |  -  |
|**400** | Target session is not allowed to host. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **controlRelease**
> controlRelease()

Release control of the room.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

const { status, data } = await apiInstance.controlRelease();
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
|**204** | Control released successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | There is already a host. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **controlRequest**
> controlRequest()

Request control of the room.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

const { status, data } = await apiInstance.controlRequest();
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
|**204** | Control requested successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | There is already a host. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **controlReset**
> controlReset()

Reset the control status of the room.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

const { status, data } = await apiInstance.controlReset();
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
|**204** | Control reset successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **controlStatus**
> ControlStatus controlStatus()

Retrieve the current control status of the room.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

const { status, data } = await apiInstance.controlStatus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ControlStatus**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Control status retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **controlTake**
> controlTake()

Take control of the room.

### Example

```typescript
import {
    RoomControlApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomControlApi(configuration);

const { status, data } = await apiInstance.controlTake();
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
|**204** | Control taken successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

