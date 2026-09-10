# RoomClipboardApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**clipboardGetImage**](#clipboardgetimage) | **GET** /api/room/clipboard/image.png | Get Clipboard Image|
|[**clipboardGetText**](#clipboardgettext) | **GET** /api/room/clipboard | Get Clipboard Content|
|[**clipboardSetText**](#clipboardsettext) | **POST** /api/room/clipboard | Set Clipboard Content|

# **clipboardGetImage**
> File clipboardGetImage()

Retrieve the current image content of the clipboard.

### Example

```typescript
import {
    RoomClipboardApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomClipboardApi(configuration);

const { status, data } = await apiInstance.clipboardGetImage();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**File**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/png, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Clipboard image retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to get clipboard content. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **clipboardGetText**
> ClipboardText clipboardGetText()

Retrieve the current content of the clipboard (rich-text or plain-text).

### Example

```typescript
import {
    RoomClipboardApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomClipboardApi(configuration);

const { status, data } = await apiInstance.clipboardGetText();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ClipboardText**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Clipboard content retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to get clipboard content. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **clipboardSetText**
> clipboardSetText(clipboardText)

Update the content of the clipboard (rich-text or plain-text).

### Example

```typescript
import {
    RoomClipboardApi,
    Configuration,
    ClipboardText
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomClipboardApi(configuration);

let clipboardText: ClipboardText; //

const { status, data } = await apiInstance.clipboardSetText(
    clipboardText
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **clipboardText** | **ClipboardText**|  | |


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
|**204** | Clipboard content updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to set clipboard content. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

