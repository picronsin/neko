# SessionsApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**sessionDisconnect**](#sessiondisconnect) | **POST** /api/sessions/{sessionId}/disconnect | Disconnect Session|
|[**sessionGet**](#sessionget) | **GET** /api/sessions/{sessionId} | Get Session|
|[**sessionRemove**](#sessionremove) | **DELETE** /api/sessions/{sessionId} | Remove Session|
|[**sessionsGet**](#sessionsget) | **GET** /api/sessions | List Sessions|

# **sessionDisconnect**
> sessionDisconnect()

Forcefully disconnect a specific session.

### Example

```typescript
import {
    SessionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SessionsApi(configuration);

let sessionId: string; //The identifier of the session. (default to undefined)

const { status, data } = await apiInstance.sessionDisconnect(
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
|**204** | Session disconnected successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **sessionGet**
> SessionData sessionGet()

Retrieve information about a specific session.

### Example

```typescript
import {
    SessionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SessionsApi(configuration);

let sessionId: string; //The identifier of the session. (default to undefined)

const { status, data } = await apiInstance.sessionGet(
    sessionId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sessionId** | [**string**] | The identifier of the session. | defaults to undefined|


### Return type

**SessionData**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Session retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **sessionRemove**
> sessionRemove()

Terminate a specific session.

### Example

```typescript
import {
    SessionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SessionsApi(configuration);

let sessionId: string; //The identifier of the session. (default to undefined)

const { status, data } = await apiInstance.sessionRemove(
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
|**204** | Session removed successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **sessionsGet**
> Array<SessionData> sessionsGet()

Retrieve a list of all active sessions.

### Example

```typescript
import {
    SessionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SessionsApi(configuration);

const { status, data } = await apiInstance.sessionsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<SessionData>**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Sessions retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

