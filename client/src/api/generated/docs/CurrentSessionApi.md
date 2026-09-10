# CurrentSessionApi

All URIs are relative to *http://localhost:3000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**login**](#login) | **POST** /api/login | User Login|
|[**logout**](#logout) | **POST** /api/logout | User Logout|
|[**profile**](#profile) | **POST** /api/profile | Update Profile|
|[**profileAvatar**](#profileavatar) | **POST** /api/profile/avatar | Update Current User Avatar|
|[**whoami**](#whoami) | **GET** /api/whoami | Get Current User|

# **login**
> SessionLoginResponse login(sessionLoginRequest)

Authenticate a user and start a new session.

### Example

```typescript
import {
    CurrentSessionApi,
    Configuration,
    SessionLoginRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new CurrentSessionApi(configuration);

let sessionLoginRequest: SessionLoginRequest; //

const { status, data } = await apiInstance.login(
    sessionLoginRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sessionLoginRequest** | **SessionLoginRequest**|  | |


### Return type

**SessionLoginResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User authenticated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **logout**
> logout()

Terminate the current user session.

### Example

```typescript
import {
    CurrentSessionApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CurrentSessionApi(configuration);

const { status, data } = await apiInstance.logout();
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
|**200** | User logged out successfully. |  -  |
|**401** | The request requires user authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **profile**
> profile(memberProfile)

Update the current user\'s display name and avatar; administrators may update the complete profile.

### Example

```typescript
import {
    CurrentSessionApi,
    Configuration,
    MemberProfile
} from './api';

const configuration = new Configuration();
const apiInstance = new CurrentSessionApi(configuration);

let memberProfile: MemberProfile; //

const { status, data } = await apiInstance.profile(
    memberProfile
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberProfile** | **MemberProfile**|  | |


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
|**204** | Profile updated successfully. |  -  |
|**401** | The request requires user authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **profileAvatar**
> profileAvatar(profileAvatarRequest)

Store a PNG, JPEG, or GIF avatar as a validated base64 data URL and broadcast the profile change to the room.

### Example

```typescript
import {
    CurrentSessionApi,
    Configuration,
    ProfileAvatarRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new CurrentSessionApi(configuration);

let profileAvatarRequest: ProfileAvatarRequest; //

const { status, data } = await apiInstance.profileAvatar(
    profileAvatarRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **profileAvatarRequest** | **ProfileAvatarRequest**|  | |


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
|**204** | Avatar updated successfully. |  -  |
|**400** | Invalid or oversized avatar. |  -  |
|**401** | The request requires user authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **whoami**
> SessionData whoami()

Retrieve information about the current user session.

### Example

```typescript
import {
    CurrentSessionApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CurrentSessionApi(configuration);

const { status, data } = await apiInstance.whoami();
```

### Parameters
This endpoint does not have any parameters.


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
|**200** | Current user information retrieved successfully. |  -  |
|**401** | The request requires user authentication. |  -  |
|**403** | The server understood the request, but refuses to authorize it. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

