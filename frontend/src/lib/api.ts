async function parseJson<T>(res: Response): Promise<T> {
  const data = (await res.json()) as T & { error?: string };
  if (!res.ok || data.error) {
    throw new Error(data.error || `请求失败：${res.status}`);
  }
  return data;
}

export async function getJson<T>(path: string): Promise<T> {
  const res = await fetch(path, { credentials: "same-origin" });
  return parseJson<T>(res);
}

export async function postJson<T>(path: string, body: unknown = {}): Promise<T> {
  const res = await fetch(path, {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return parseJson<T>(res);
}
