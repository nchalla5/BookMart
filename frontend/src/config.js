const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || "http://localhost:8080";

export function apiUrl(path) {
  return `${API_BASE_URL}${path}`;
}

export function resolveImageUrl(path) {
  if (!path) {
    return "";
  }

  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  if (path.startsWith("/")) {
    return `${API_BASE_URL}${path}`;
  }

  return `${API_BASE_URL}/${path}`;
}

export { API_BASE_URL };
