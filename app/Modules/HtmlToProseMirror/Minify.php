<?php

namespace App\Modules\HtmlToProseMirror;

class Minify
{
    protected $_replacementHash;
    protected $_placeholders = [];
    protected $_html;

    public function process($html)
    {
        $this->_html = str_replace("\r\n", "\n", trim($html));

        $this->_replacementHash = 'MINIFYHTML' . md5($_SERVER['REQUEST_TIME']);

        // replace PREs with placeholders
        $this->_html = preg_replace_callback('/\\s*<pre(\\b[^>]*?>[\\s\\S]*?<\\/pre>)\\s*/iu', [$this, '_removePreCB'], $this->_html);

        // trim each line.
        $this->_html = preg_replace('/^\\s+|\\s+$/mu', '', $this->_html);

        // remove ws around block/undisplayed elements
        $this->_html = preg_replace('/\\s+(<\\/?(?:area|article|aside|base(?:font)?|blockquote|body'
            . '|canvas|caption|center|col(?:group)?|dd|dir|div|dl|dt|fieldset|figcaption|figure|footer|form'
            . '|frame(?:set)?|h[1-6]|head|header|hgroup|hr|html|legend|li|link|main|map|menu|meta|nav'
            . '|ol|opt(?:group|ion)|output|p|param|section|t(?:able|body|head|d|h||r|foot|itle)'
            . '|ul|video)\\b[^>]*>)/iu', '$1', $this->_html);

        // fill placeholders
        $this->_html = str_replace(
            array_keys($this->_placeholders),
            array_values($this->_placeholders),
            $this->_html
        );

        return $this->_html;
    }

    protected function _removePreCB($m)
    {
        return $this->_reservePlace("<pre{$m[1]}");
    }

    protected function _reservePlace($content)
    {
        $placeholder = '%' . $this->_replacementHash . count($this->_placeholders) . '%';
        $this->_placeholders[$placeholder] = $content;

        return $placeholder;
    }
}
