import React from "react";
import { ImageResponse } from "@vercel/og";
import { createClient } from "@libsql/client";

let client = null;

const cache = new Map();
const CACHE_TTL = 30 * 60 * 1000;

function stringToColor(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = "";
  for (let i = 0; i < 3; i++) {
    const value = (hash >> (i * 8)) & 0xff;
    color += ("00" + value.toString(16)).slice(-2);
  }
  return color;
}

function generateGradient(str) {
  const color1 = stringToColor(str);
  const color2 = stringToColor(str.split("").reverse().join(""));
  return `linear-gradient(135deg, #${color1}, #${color2})`;
}

function getTursoClient() {
  if (client) return client;

  const url = process.env.TURSO_DATABASE_URL;
  const authToken = process.env.TURSO_AUTH_TOKEN;
  if (!url || !authToken) {
    console.error("Missing Turso environment variables");
    return null;
  }

  client = createClient({ url, authToken });
  return client;
}

async function getStory(id) {
  const cached = cache.get(id);
  if (cached && cached.expires > Date.now()) {
    return { story: { title: cached.title, domain: cached.domain }, cacheHit: true };
  }

  try {
    const tursoClient = getTursoClient();
    if (!tursoClient) {
      return { story: null, error: "Turso client is not configured" };
    }

    const result = await tursoClient.execute({
      sql: "SELECT title, domain FROM top_stories WHERE hn_id = ?",
      args: [id],
    });

    if (result.rows.length === 0) {
      return { story: null, error: null };
    }

    const story = {
      title: result.rows[0].title,
      domain: result.rows[0].domain,
    };

    cache.set(id, { ...story, expires: Date.now() + CACHE_TTL });

    return { story, cacheHit: false, error: null };
  } catch (error) {
    console.error("Turso query failed", { storyId: id, error: error.message });
    return { story: null, error: error.message };
  }
}

function toHeaders(rawHeaders) {
  if (!rawHeaders) return new Headers();
  if (rawHeaders instanceof Headers) return rawHeaders;

  const headers = new Headers();
  for (const [key, value] of Object.entries(rawHeaders)) {
    if (Array.isArray(value)) {
      for (const item of value) {
        if (item != null) headers.append(key, String(item));
      }
    } else if (value != null) {
      headers.set(key, String(value));
    }
  }
  return headers;
}

function normalizeRequest(request) {
  if (request instanceof Request) return request;

  const headers = toHeaders(request?.headers);
  const rawUrl = request?.url || "/";
  const isAbsolute = /^https?:\/\//i.test(rawUrl);
  const protocol = headers.get("x-forwarded-proto") || "https";
  const host = headers.get("x-forwarded-host") || headers.get("host") || "localhost";
  const url = isAbsolute ? rawUrl : `${protocol}://${host}${rawUrl}`;
  const method = request?.method || "GET";

  return { url, method, headers };
}

function isNodeResponse(response) {
  return !!response && typeof response.setHeader === "function" && typeof response.end === "function";
}

async function sendResponse(maybeNodeResponse, fetchResponse) {
  if (!isNodeResponse(maybeNodeResponse)) {
    return fetchResponse;
  }

  maybeNodeResponse.statusCode = fetchResponse.status;
  for (const [key, value] of fetchResponse.headers.entries()) {
    maybeNodeResponse.setHeader(key, value);
  }

  const body = await fetchResponse.arrayBuffer();
  maybeNodeResponse.end(new Uint8Array(body));
  return;
}

export default async function handler(request, response) {

  const normalizedRequest = normalizeRequest(request);

  const url = new URL(normalizedRequest.url);
  const idParam = url.searchParams.get("id");

  if (!idParam) {
    return new Response("Missing id parameter", { status: 400 });
  }

  const id = parseInt(idParam, 10);
  if (isNaN(id) || id <= 0) {
    return new Response("Invalid id parameter", { status: 400 });
  }

  const { story, error } = await getStory(id);

  if (error) {
    return new Response("Failed to load story", { status: 500 });
  }

  if (!story) {
    return new Response("Story not in top 30", { status: 404 });
  }

  const [fontExtraBoldData, fontRegularData] = await Promise.all([
    fetch(`${url.origin}/assets/Inter-ExtraBold.ttf`).then((res) => res.arrayBuffer()),
    fetch(`${url.origin}/assets/Inter-Regular.ttf`).then((res) => res.arrayBuffer()),
  ]);

  const gradient = generateGradient(story.title);

  return sendResponse(response, new ImageResponse(
    React.createElement(
      "div",
      {
        style: {
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "flex-start",
          fontFamily: "Inter",
          backgroundImage: gradient,
          padding: "60px",
        },
      },
      React.createElement(
        "div",
        {
          style: {
            display: "flex",
            flexDirection: "column",
            alignItems: "flex-start",
            maxWidth: "85%",
          },
        },
        React.createElement(
          "div",
          {
            style: {
              fontSize: 80,
              color: "#ffffff",
              textAlign: "left",
              fontWeight: 900,
              lineHeight: 1.1,
              textShadow: "0 2px 10px rgba(0,0,0,0.3)",
              marginBottom: "20px",
            },
          },
          story.title
        ),
        story.domain &&
          React.createElement(
            "div",
            {
              style: {
                fontSize: 28,
                color: "rgba(255,255,255,0.85)",
                textShadow: "0 1px 4px rgba(0,0,0,0.2)",
              },
            },
            story.domain
          )
      )
    ),
    {
      width: 1200,
      height: 630,
      fonts: [
        {
          name: "Inter",
          data: fontRegularData,
          style: "normal",
          weight: 400,
        },
        {
          name: "Inter",
          data: fontExtraBoldData,
          style: "normal",
          weight: 900,
        },
      ],
    }
  ));
}
