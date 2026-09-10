# RoomKeyboardApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**keyboardMapGet**](#keyboardmapget) | **GET** /api/room/keyboard/map | Get Keyboard Map|
|[**keyboardMapSet**](#keyboardmapset) | **POST** /api/room/keyboard/map | Set Keyboard Map|
|[**keyboardModifiersGet**](#keyboardmodifiersget) | **GET** /api/room/keyboard/modifiers | Get Keyboard Modifiers|
|[**keyboardModifiersSet**](#keyboardmodifiersset) | **POST** /api/room/keyboard/modifiers | Set Keyboard Modifiers|

# **keyboardMapGet**
> KeyboardMap keyboardMapGet()

Retrieve the current keyboard map configuration.

### Example

```typescript
import {
    RoomKeyboardApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomKeyboardApi(configuration);

const { status, data } = await apiInstance.keyboardMapGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**KeyboardMap**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Keyboard map retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to get keyboard map. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **keyboardMapSet**
> keyboardMapSet(keyboardMap)

Update the keyboard map configuration.

### Example

```typescript
import {
    RoomKeyboardApi,
    Configuration,
    KeyboardMap
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomKeyboardApi(configuration);

let keyboardMap: KeyboardMap; //

const { status, data } = await apiInstance.keyboardMapSet(
    keyboardMap
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **keyboardMap** | **KeyboardMap**|  | |


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
|**204** | Keyboard map updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to change keyboard map. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **keyboardModifiersGet**
> KeyboardModifiers keyboardModifiersGet()

Retrieve the current keyboard modifiers status.

### Example

```typescript
import {
    RoomKeyboardApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomKeyboardApi(configuration);

const { status, data } = await apiInstance.keyboardModifiersGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**KeyboardModifiers**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Keyboard modifiers retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **keyboardModifiersSet**
> keyboardModifiersSet(keyboardModifiers)

Update the keyboard modifiers status.

### Example

```typescript
import {
    RoomKeyboardApi,
    Configuration,
    KeyboardModifiers
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomKeyboardApi(configuration);

let keyboardModifiers: KeyboardModifiers; //

const { status, data } = await apiInstance.keyboardModifiersSet(
    keyboardModifiers
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **keyboardModifiers** | **KeyboardModifiers**|  | |


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
|**204** | Keyboard modifiers updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

