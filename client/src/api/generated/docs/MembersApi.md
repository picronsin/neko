# MembersApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**membersBulkDelete**](#membersbulkdelete) | **POST** /api/members_bulk/delete | Bulk Delete Members|
|[**membersBulkUpdate**](#membersbulkupdate) | **POST** /api/members_bulk/update | Bulk Update Members|
|[**membersCreate**](#memberscreate) | **POST** /api/members | Create Member|
|[**membersGetProfile**](#membersgetprofile) | **GET** /api/members/{memberId} | Get Member Profile|
|[**membersList**](#memberslist) | **GET** /api/members | List Members|
|[**membersRemove**](#membersremove) | **DELETE** /api/members/{memberId} | Remove Member|
|[**membersUpdatePassword**](#membersupdatepassword) | **POST** /api/members/{memberId}/password | Update Member Password|
|[**membersUpdateProfile**](#membersupdateprofile) | **POST** /api/members/{memberId} | Update Member Profile|

# **membersBulkDelete**
> membersBulkDelete(memberBulkDelete)

Remove multiple members in bulk.

### Example

```typescript
import {
    MembersApi,
    Configuration,
    MemberBulkDelete
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberBulkDelete: MemberBulkDelete; //

const { status, data } = await apiInstance.membersBulkDelete(
    memberBulkDelete
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberBulkDelete** | **MemberBulkDelete**|  | |


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
|**204** | Members removed successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersBulkUpdate**
> membersBulkUpdate(memberBulkUpdate)

Update the profiles of multiple members in bulk.

### Example

```typescript
import {
    MembersApi,
    Configuration,
    MemberBulkUpdate
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberBulkUpdate: MemberBulkUpdate; //

const { status, data } = await apiInstance.membersBulkUpdate(
    memberBulkUpdate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberBulkUpdate** | **MemberBulkUpdate**|  | |


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
|**204** | Members updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersCreate**
> MemberData membersCreate(memberCreate)

Create a new member.

### Example

```typescript
import {
    MembersApi,
    Configuration,
    MemberCreate
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberCreate: MemberCreate; //

const { status, data } = await apiInstance.membersCreate(
    memberCreate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberCreate** | **MemberCreate**|  | |


### Return type

**MemberData**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Member created successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**422** | Member with chosen ID already exists. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersGetProfile**
> MemberProfile membersGetProfile()

Retrieve the profile of a specific member.

### Example

```typescript
import {
    MembersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberId: string; //The identifier of the member. (default to undefined)

const { status, data } = await apiInstance.membersGetProfile(
    memberId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberId** | [**string**] | The identifier of the member. | defaults to undefined|


### Return type

**MemberProfile**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Member profile retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersList**
> Array<MemberData> membersList()

Retrieve a list of all members.

### Example

```typescript
import {
    MembersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let limit: number; // (optional) (default to undefined)
let offset: number; // (optional) (default to undefined)

const { status, data } = await apiInstance.membersList(
    limit,
    offset
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **limit** | [**number**] |  | (optional) defaults to undefined|
| **offset** | [**number**] |  | (optional) defaults to undefined|


### Return type

**Array<MemberData>**

### Authorization

[CookieAuth](../README.md#CookieAuth), [TokenAuth](../README.md#TokenAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Members retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersRemove**
> membersRemove()

Remove a specific member.

### Example

```typescript
import {
    MembersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberId: string; //The identifier of the member. (default to undefined)

const { status, data } = await apiInstance.membersRemove(
    memberId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberId** | [**string**] | The identifier of the member. | defaults to undefined|


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
|**204** | Member removed successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersUpdatePassword**
> membersUpdatePassword(memberPassword)

Update the password of a specific member.

### Example

```typescript
import {
    MembersApi,
    Configuration,
    MemberPassword
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberId: string; //The identifier of the member. (default to undefined)
let memberPassword: MemberPassword; //

const { status, data } = await apiInstance.membersUpdatePassword(
    memberId,
    memberPassword
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberPassword** | **MemberPassword**|  | |
| **memberId** | [**string**] | The identifier of the member. | defaults to undefined|


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
|**204** | Member password updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **membersUpdateProfile**
> membersUpdateProfile(memberProfile)

Update the profile of a specific member.

### Example

```typescript
import {
    MembersApi,
    Configuration,
    MemberProfile
} from './api';

const configuration = new Configuration();
const apiInstance = new MembersApi(configuration);

let memberId: string; //The identifier of the member. (default to undefined)
let memberProfile: MemberProfile; //

const { status, data } = await apiInstance.membersUpdateProfile(
    memberId,
    memberProfile
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberProfile** | **MemberProfile**|  | |
| **memberId** | [**string**] | The identifier of the member. | defaults to undefined|


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
|**204** | Member profile updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |
|**404** | The specified resource was not found. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

