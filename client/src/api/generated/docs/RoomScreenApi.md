# RoomScreenApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**screenCastImage**](#screencastimage) | **GET** /api/room/screen/cast.jpg | Get Screencast Image|
|[**screenConfiguration**](#screenconfiguration) | **GET** /api/room/screen | Get Screen Configuration|
|[**screenConfigurationChange**](#screenconfigurationchange) | **POST** /api/room/screen | Change Screen Configuration|
|[**screenConfigurationsList**](#screenconfigurationslist) | **GET** /api/room/screen/configurations | Get List of Screen Configurations|
|[**screenShotImage**](#screenshotimage) | **GET** /api/room/screen/shot.jpg | Get Screenshot Image|

# **screenCastImage**
> File screenCastImage()

Retrieve the current screencast image.

### Example

```typescript
import {
    RoomScreenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomScreenApi(configuration);

const { status, data } = await apiInstance.screenCastImage();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**File**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/jpeg, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Screencast image retrieved successfully. |  -  |
|**400** | Screencast is not enabled. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to fetch image. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **screenConfiguration**
> ScreenConfiguration screenConfiguration()

Retrieve the current screen configuration of the room.

### Example

```typescript
import {
    RoomScreenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomScreenApi(configuration);

const { status, data } = await apiInstance.screenConfiguration();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ScreenConfiguration**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Screen configuration retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to get screen configuration. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **screenConfigurationChange**
> ScreenConfiguration screenConfigurationChange(screenConfiguration)

Update the screen configuration of the room.

### Example

```typescript
import {
    RoomScreenApi,
    Configuration,
    ScreenConfiguration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomScreenApi(configuration);

let screenConfiguration: ScreenConfiguration; //

const { status, data } = await apiInstance.screenConfigurationChange(
    screenConfiguration
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **screenConfiguration** | **ScreenConfiguration**|  | |


### Return type

**ScreenConfiguration**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Screen configuration updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | Invalid screen configuration. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **screenConfigurationsList**
> Array<ScreenConfiguration> screenConfigurationsList()

Retrieve a list of all available screen configurations.

### Example

```typescript
import {
    RoomScreenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomScreenApi(configuration);

const { status, data } = await apiInstance.screenConfigurationsList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ScreenConfiguration>**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of screen configurations retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **screenShotImage**
> File screenShotImage()

Retrieve the current screenshot image.

### Example

```typescript
import {
    RoomScreenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomScreenApi(configuration);

let quality: number; //Image quality (0-100). (optional) (default to undefined)

const { status, data } = await apiInstance.screenShotImage(
    quality
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **quality** | [**number**] | Image quality (0-100). | (optional) defaults to undefined|


### Return type

**File**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/jpeg, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Screenshot image retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to create image. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

