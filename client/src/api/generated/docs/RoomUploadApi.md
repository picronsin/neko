# RoomUploadApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**uploadDialog**](#uploaddialog) | **POST** /api/room/upload/dialog | Upload File to Dialog|
|[**uploadDialogClose**](#uploaddialogclose) | **DELETE** /api/room/upload/dialog | Close File Chooser Dialog|
|[**uploadDrop**](#uploaddrop) | **POST** /api/room/upload/drop | Upload and Drop File|

# **uploadDialog**
> uploadDialog()

Upload a file to a dialog.

### Example

```typescript
import {
    RoomUploadApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomUploadApi(configuration);

let files: Array<File>; //Files to be uploaded. (optional) (default to undefined)

const { status, data } = await apiInstance.uploadDialog(
    files
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **files** | **Array&lt;File&gt;** | Files to be uploaded. | (optional) defaults to undefined|


### Return type

void (empty response body)

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | File uploaded to dialog successfully. |  -  |
|**400** | Unable to upload file. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | No upload dialog prompt active. |  -  |
|**500** | Unable to process uploaded file. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uploadDialogClose**
> uploadDialogClose()

Close the file chooser dialog.

### Example

```typescript
import {
    RoomUploadApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomUploadApi(configuration);

const { status, data } = await apiInstance.uploadDialogClose();
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
|**204** | File chooser dialog closed successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | No upload dialog prompt active. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uploadDrop**
> uploadDrop()

Upload a file and drop it at a specified location.

### Example

```typescript
import {
    RoomUploadApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoomUploadApi(configuration);

let x: number; //X coordinate of drop. (optional) (default to undefined)
let y: number; //Y coordinate of drop. (optional) (default to undefined)
let files: Array<File>; //Files to be uploaded. (optional) (default to undefined)

const { status, data } = await apiInstance.uploadDrop(
    x,
    y,
    files
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **x** | [**number**] | X coordinate of drop. | (optional) defaults to undefined|
| **y** | [**number**] | Y coordinate of drop. | (optional) defaults to undefined|
| **files** | **Array&lt;File&gt;** | Files to be uploaded. | (optional) defaults to undefined|


### Return type

void (empty response body)

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | File uploaded and dropped successfully. |  -  |
|**400** | Unable to upload file. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**500** | Unable to process uploaded file. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

