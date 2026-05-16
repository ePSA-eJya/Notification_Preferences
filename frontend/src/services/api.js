const BASE_URL = "/api/v1";

function getToken() {
  return localStorage.getItem("token");
}

async function request(method, path, body = null) {
  const headers = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const opts = { method, headers };
  if (body) {
    opts.body = JSON.stringify(body);
  }

  const res = await fetch(`${BASE_URL}${path}`, opts);
  const data = await res.json().catch(() => null);

  if (!res.ok) {
    const message =
      data?.message || data?.error || `Request failed (${res.status})`;
    throw new Error(message);
  }

  return data;
}

// Request function for multipart form data
async function requestMultipart(method, path, formData) {
  const headers = {};
  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const opts = { method, headers, body: formData };

  const res = await fetch(`${BASE_URL}${path}`, opts);
  const data = await res.json().catch(() => null);

  if (!res.ok) {
    const message =
      data?.message || data?.error || `Request failed (${res.status})`;
    throw new Error(message);
  }

  return data;
}

// ===== Auth =====
// export const authAPI = {
//   signup: (email, password, user_handle) =>
//     request('POST', '/auth/signup', { email, password, user_handle }),

//   signin: (email, password) =>
//     request('POST', '/auth/signin', { email, password }),
// };

export const authAPI = {
  signup: (email, password, user_handle, deviceToken) =>
    request("POST", "/auth/signup", {
      email,
      password,
      user_handle,
      deviceToken,
    }),

  // signin: (email, password, deviceToken) =>
  //   request('POST', '/auth/signin', { email, password, deviceToken, }),
  signin: (email, password, deviceToken) => {
    return request("POST", "/auth/signin", {
      email,
      password,
      deviceToken,
    });
  },
};

// ===== Users =====
export const usersAPI = {
  getAll: () => request("GET", "/users/"),
  getById: (id) => request("GET", `/users/${id}`),
  getMe: () => request("GET", "/me"),

  follow: (userId) => request("POST", `/users/${userId}/follow`),
  unfollow: (userId) => request("DELETE", `/users/${userId}/follow`),
  getFollowers: () => request("GET", "/followers"),
  getFollowing: () => request("GET", "/following"),
};

// ===== Feed / Posts =====
export const feedAPI = {
  getFeed: (limit = 20, offset = 0) =>
    request("GET", `/feed?limit=${limit}&offset=${offset}`),

  createPost: (content, mediaFiles = []) => {
    const formData = new FormData();
    formData.append("content", content);

    // Append media files if provided
    if (mediaFiles && mediaFiles.length > 0) {
      mediaFiles.forEach((file) => {
        formData.append("media", file);
      });
    }

    return requestMultipart("POST", "/posts", formData);
  },

  likePost: (postId) => request("POST", `/posts/${postId}/like`),

  unlikePost: (postId) => request("DELETE", `/posts/${postId}/like`),

  isPostLiked: (postId) => request("GET", `/posts/${postId}/liked`),

  commentOnPost: (postId, text) =>
    request("POST", `/posts/${postId}/comment`, { text }),

  getPostComments: (postId, limit = 50, offset = 0) =>
    request("GET", `/posts/${postId}/comments?limit=${limit}&offset=${offset}`),
};

// ===== Preferences =====
export const preferencesAPI = {
  get: () => request("GET", "/preferences"),
  update: (prefs) => request("PUT", "/preferences", prefs),
};

// ===== Notifications =====
export const notificationsAPI = {
  getAll: (limit = 50) => request("GET", `/notifications?limit=${limit}`),
  markAllRead: () => request("POST", "/notifications/read"),
};
