import axios from "axios"

const HTTP = axios.create({
  baseURL: `${import.meta.env.VITE_API_HOST}/${import.meta.env.VITE_API_PATH}/${import.meta.env.VITE_API_VERSION}`,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-CLIENT": "web",
  },
  withCredentials: true,
})

// CRSF token
HTTP.interceptors.request.use((config) => {
  // 쿠키에서 토큰을 가져와서 헤더에 넣어준다.  __Host-csrf_
  const token = document.cookie
    .split("; ")
    ?.find((row) => row.startsWith("_Host-csrf_"))
    ?.split("=")[1]
  console.log(token, "token", document.cookie)
  if (token) {
    config.headers["X-CSRF-Token"] = token
  }
  return config
})

const getOnlyData = async (req: Promise<{ data: unknown }>) => {
  try {
    const res = await req
    return res.data
  } catch (error) {
    console.error(error)
  }
}

const repository = <T>(resource: string) => {
  if (!resource) {
    throw new Error("resource is required")
  }
  if (!HTTP) {
    throw new Error("HTTP is required")
  }
  if (/^\//.test(resource)) {
    throw new Error("resource should not start with '/'")
  }
  return {
    fetchAll: () => HTTP.get(resource),
    fetch: (id: string) => getOnlyData(HTTP.get(`${resource}/${id}`)),
    create: (payload: T) => getOnlyData(HTTP.post(resource, payload)),
    update: (id: string, payload: T) =>
      getOnlyData(HTTP.put(`${resource}/${id}`, payload)),
    remove: (id: string) => getOnlyData(HTTP.delete(`${resource}/${id}`)),
  }
}

export default HTTP
export { repository, getOnlyData }
