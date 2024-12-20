import DOMPurify from 'dompurify';
import { Converter } from 'showdown';

export const markdownConverter = {
  makeHtml: (markdown: string): string => {
    const unsafeHtml = new Converter({
      tables: true,
      openLinksInNewWindow: true,
      strikethrough: true,
      emoji: true,
      literalMidWordUnderscores: true,
    }).makeHtml(markdown);

    // add hook to transform anchor tags
    DOMPurify.addHook('beforeSanitizeElements', (node) => {
      if (node instanceof HTMLAnchorElement) {
        node.setAttribute('rel', 'noopener noreferrer');
      }
    });

    return DOMPurify.sanitize(unsafeHtml, {
      ALLOWED_TAGS: [
        'b',
        'i',
        'strike',
        's',
        'del',
        'em',
        'strong',
        'a',
        'p',
        'h1',
        'h2',
        'h3',
        'h4',
        'ul',
        'ol',
        'li',
        'code',
        'pre',
        'table',
        'thead',
        'tbody',
        'tr',
        'th',
        'td',
      ],
      ALLOWED_ATTR: ['href', 'target', 'rel', 'class'],
      ALLOWED_URI_REGEXP: /^(?:(?:https?|mailto):|[^a-z]|[a-z+.-]+(?:[^a-z+.\-:]|$))/i,
    });
  },
};
