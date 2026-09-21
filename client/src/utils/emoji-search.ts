export function filterEmoji(list: string[], keywords: Record<string, string[]>, search: string): string[] {
  const query = search.trim().toLowerCase()
  return list.filter(
    (emoji) =>
      emoji.toLowerCase().includes(query) ||
      (keywords[emoji]?.some((keyword) => keyword.toLowerCase().includes(query)) ?? false),
  )
}
