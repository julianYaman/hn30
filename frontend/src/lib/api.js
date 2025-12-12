export async function getTopStories(fetchFn, apiBaseUrl) {
  const res = await fetchFn(`${apiBaseUrl}/api/top`);
  if (res.ok) {
    return await res.json();
  } else {
    throw new Error('Failed to fetch stories');
  }
}

export async function getSummary(storyId) {
  const res = await fetch(`/api/summarize?id=${storyId}`);
  if (res.ok) {
    return await res.json();
  } else {
    const errorText = await res.text();
    throw new Error(`An error occurred: ${errorText}`);
  }
}