import createMarkdown from '../lib/markdown/rules'

export default function renderToHtml(markdown) {
    return createMarkdown({ embeds: {} }).render(markdown).trim()
}
