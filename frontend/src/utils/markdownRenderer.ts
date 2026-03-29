import MarkdownIt from 'markdown-it'
import katex from 'katex'
import hljs from 'highlight.js'

// 创建 Markdown-it 实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str: string, lang: string) {
    // 代码高亮
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch (__) {}
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// 渲染数学公式
const renderMath = (content: string): string => {
  // 块级公式 $$...$$
  content = content.replace(/\$\$([\s\S]+?)\$\$/g, (match, formula) => {
    try {
      return katex.renderToString(formula.trim(), {
        displayMode: true,
        throwOnError: false,
        errorColor: '#cc0000'
      })
    } catch (e) {
      return `<span class="katex-error">${formula}</span>`
    }
  })

  // 行内公式 $...$
  content = content.replace(/\$([^\$\n]+?)\$/g, (match, formula) => {
    try {
      return katex.renderToString(formula.trim(), {
        displayMode: false,
        throwOnError: false,
        errorColor: '#cc0000'
      })
    } catch (e) {
      return `<span class="katex-error">${formula}</span>`
    }
  })

  return content
}

// 渲染 Markdown
export const renderMarkdown = (content: string): string => {
  if (!content) return ''

  // 先渲染数学公式
  const contentWithMath = renderMath(content)

  // 再渲染 Markdown
  return md.render(contentWithMath)
}
